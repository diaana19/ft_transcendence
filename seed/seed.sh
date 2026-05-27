#!/usr/bin/env bash
#
# seed.sh — populate a running Transcendence instance through its public HTTP
# API (curl), mirroring the data created by backend/cmd/seed/main.go but WITHOUT
# touching the database directly.
#
# It creates 10 users and, between them: accepted + pending friend requests,
# follows, posts, likes and comments. Friend-request notifications are produced
# automatically by the API as a side effect of the friend-request calls.
#
# Requirements: bash 4+, curl, jq. The backend must be running.
#
# Usage:
#   ./seed.sh                                  # against the proxy at https://localhost:3000
#   BASE_URL=https://host:port ./seed.sh       # custom proxy URL (uses -k for self-signed TLS)
#   SEED_PASSWORD='MyPass123!' ./seed.sh
#
# Notes / intentional differences from the Go seeder:
#   - The register endpoint validates passwords, so every seeded user shares one
#     valid password (SEED_PASSWORD) rather than the placeholder in users.json.
#   - Names and bios are set with a follow-up profile update (register only
#     accepts username/email/password/dateOfBirth).
#   - Avatars/wallpapers are skipped: they require uploading real image files.
#   - The Go seeder's like rule compares UUID lengths, which are always equal,
#     so it never actually likes anything; this script instead likes on an
#     index parity so the dataset is useful.
#   - /auth/register and /auth/login are rate-limited per IP (RATE_LIMIT_MAX,
#     default 20/min). Seeding 10 users issues 20 such calls — right at the
#     default — so run it on a fresh window. Users whose login is throttled or
#     who already exist with a different password are logged and skipped.

set -uo pipefail

BASE_URL="${BASE_URL:-https://localhost:3000}"
API="${BASE_URL%/}"
case "$API" in
*/api) ;;            # caller already included the /api prefix
*) API="$API/api" ;;
esac
PASSWORD="${SEED_PASSWORD:-SeedPass123!}"

for dep in curl jq; do
	command -v "$dep" >/dev/null 2>&1 || { echo "error: '$dep' is required" >&2; exit 1; }
done

# ---------------------------------------------------------------------------
# Seed data (mirrors backend/cmd/seed)
# ---------------------------------------------------------------------------
usernames=(alice01 carlos02 sophie03 john04 emma05 lucas06 mia07 ahmed08 olivia09 noah10)
emails=(
	alice01@example.com carlos02@example.com sophie03@example.com john04@example.com
	emma05@example.com lucas06@example.com mia07@example.com ahmed08@example.com
	olivia09@example.com noah10@example.com
)
names=(
	"Alice Dupont" "Carlos Méndez" "Sophie Martin" "John Smith" "Emma Garcia"
	"Lucas Silva" "Mia Chen" "Ahmed Ali" "Olivia Brown" "Noah Wilson"
)
bios=(
	"Loves photography and travel" "Backend developer" "UI/UX Designer" "Tech enthusiast"
	"Digital artist" "Music lover" "Food blogger" "Entrepreneur" "Marketing specialist"
	"Fitness coach"
)
dobs=(
	1995-03-15 1990-07-21 1998-11-05 1985-01-10 1992-09-30
	1997-04-18 1993-06-25 1988-12-12 1996-02-14 1991-08-08
)

post_contents=(
	"Just finished an amazing project! Feeling proud 🎉"
	"Anyone else love coding at 3am? Coffee is life ☕"
	"Check out this cool new tech stack I'm learning about!"
	"Weekend is here! Time to build something cool 🚀"
	"Finally deployed to production! No more bugs (hopefully) 😅"
	"Open source contributions are the best way to learn"
	"Just discovered this incredible library, game changer!"
	"Working on something secret, can't wait to share soon..."
	"The debugging journey never ends... but that's what makes it fun!"
	"New blog post is live! Check it out, feedback welcome 📝"
	"Sometimes the best code is the code you delete 🗑️"
	"Excited to announce we're hiring! Great team, great mission 💼"
)
comment_contents=(
	"Great post!" "Totally agree." "Could you share more details?" "Thanks for sharing!"
	"This is awesome!" "Same here haha" "Love this perspective." "Interesting take!"
)

N=${#usernames[@]}
declare -A user_id user_token
declare -a post_id post_author

# ---------------------------------------------------------------------------
# HTTP helpers
# ---------------------------------------------------------------------------

# post_json PATH JSON [TOKEN] — POST a JSON body, echo the response body.
post_json() {
	local path="$1" json="$2" token="${3:-}"
	local args=(-sSk -X POST "$API$path" -H 'Content-Type: application/json' -d "$json")
	[[ -n "$token" ]] && args+=(-H "Authorization: Bearer $token")
	curl "${args[@]}"
}

# post_auth PATH TOKEN — POST with no body, just auth; echo response body.
post_auth() {
	curl -sSk -X POST "$API$1" -H "Authorization: Bearer $2"
}

wait_for_api() {
	echo "Waiting for API at $API ..."
	for _ in $(seq 1 30); do
		if curl -sSk -o /dev/null "$BASE_URL/health" 2>/dev/null || curl -sSk -o /dev/null "$API/posts" 2>/dev/null; then
			return 0
		fi
		sleep 1
	done
	echo "error: API not reachable at $API" >&2
	exit 1
}

# ---------------------------------------------------------------------------
# 1. Users: register, log in, set name/bio
# ---------------------------------------------------------------------------
seed_users() {
	echo "Seeding users..."
	for i in $(seq 0 $((N - 1))); do
		local u="${usernames[$i]}" e="${emails[$i]}"
		local body resp token id

		body=$(jq -nc --arg u "$u" --arg e "$e" --arg p "$PASSWORD" --arg d "${dobs[$i]}" \
			'{username:$u, email:$e, password:$p, dateOfBirth:$d}')
		post_json /auth/register "$body" >/dev/null

		resp=$(post_json /auth/login "$(jq -nc --arg e "$e" --arg p "$PASSWORD" '{email:$e, password:$p}')")
		token=$(jq -r '.token // empty' <<<"$resp")
		id=$(jq -r '.user.id // empty' <<<"$resp")
		if [[ -z "$token" || -z "$id" ]]; then
			echo "  ! could not log in $u (already seeded with a different password?) — skipping" >&2
			continue
		fi
		user_id["$u"]="$id"
		user_token["$u"]="$token"

		# set display name + bio (register only accepts username/email/password/dob)
		curl -sSk -X PUT "$API/users/$id" -H 'Content-Type: application/json' \
			-H "Authorization: Bearer $token" \
			-d "$(jq -nc --arg n "${names[$i]}" --arg b "${bios[$i]}" '{name:$n, bio:$b}')" >/dev/null
		echo "  + $u ($id)"
	done
}

# token/id by index, only for users that registered successfully
tok() { echo "${user_token[${usernames[$1]}]:-}"; }
uid() { echo "${user_id[${usernames[$1]}]:-}"; }

# ---------------------------------------------------------------------------
# 2. Friendships: accepted (i -> i+1) and pending (i -> i+2)
# ---------------------------------------------------------------------------
seed_friendships() {
	echo "Seeding friendships..."
	local accepted=0 pending=0
	for i in $(seq 0 $((N - 1))); do
		local a=$(((i + 1) % N)) b=$(((i + 2) % N))
		local ti tia ida idi idb
		ti=$(tok "$i"); tia=$(tok "$a")
		ida=$(uid "$a"); idi=$(uid "$i"); idb=$(uid "$b")
		[[ -n "$ti" && -n "$ida" && -n "$tia" && -n "$idi" ]] || continue

		# accepted: i requests (i+1), then (i+1) accepts i
		post_auth "/friends/request/$ida" "$ti" >/dev/null
		post_auth "/friends/accept/$idi" "$tia" >/dev/null
		accepted=$((accepted + 1))

		# pending: i requests (i+2), left unaccepted
		if [[ -n "$idb" ]]; then
			post_auth "/friends/request/$idb" "$ti" >/dev/null
			pending=$((pending + 1))
		fi
	done
	echo "  + $accepted accepted, $pending pending"
}

# ---------------------------------------------------------------------------
# 3. Follows: each user follows the next two
# ---------------------------------------------------------------------------
seed_follows() {
	echo "Seeding follows..."
	local created=0
	for i in $(seq 0 $((N - 1))); do
		local ti; ti=$(tok "$i")
		[[ -n "$ti" ]] || continue
		for offset in 1 2; do
			local t=$(((i + offset) % N)) idt
			idt=$(uid "$t")
			[[ -n "$idt" ]] || continue
			post_auth "/friends/follow/$idt" "$ti" >/dev/null
			created=$((created + 1))
		done
	done
	echo "  + $created follow relationships"
}

# ---------------------------------------------------------------------------
# 4. Posts: post_contents[i] authored by user i % N
# ---------------------------------------------------------------------------
seed_posts() {
	echo "Seeding posts..."
	local created=0
	for ci in $(seq 0 $((${#post_contents[@]} - 1))); do
		local ai=$((ci % N)) ti resp id
		ti=$(tok "$ai")
		[[ -n "$ti" ]] || continue
		resp=$(post_json /posts "$(jq -nc --arg c "${post_contents[$ci]}" '{content:$c}')" "$ti")
		id=$(jq -r '.id // empty' <<<"$resp")
		post_id[ci]="$id"
		post_author[ci]=$ai
		[[ -n "$id" ]] && created=$((created + 1))
	done
	echo "  + $created posts"
}

# ---------------------------------------------------------------------------
# 5. Likes: every user likes posts (not their own) on an index parity
# ---------------------------------------------------------------------------
seed_likes() {
	echo "Seeding likes..."
	local created=0
	for ci in $(seq 0 $((${#post_contents[@]} - 1))); do
		local pid=${post_id[$ci]:-}
		[[ -n "$pid" ]] || continue
		for ui in $(seq 0 $((N - 1))); do
			[[ "$ui" -eq "${post_author[$ci]}" ]] && continue
			[[ $(((ci + ui) % 2)) -eq 1 ]] || continue
			local ti; ti=$(tok "$ui")
			[[ -n "$ti" ]] || continue
			post_auth "/posts/$pid/like" "$ti" >/dev/null
			created=$((created + 1))
		done
	done
	echo "  + $created likes"
}

# ---------------------------------------------------------------------------
# 6. Comments: 1 or 2 per post, from rotating non-author users
# ---------------------------------------------------------------------------
seed_comments() {
	echo "Seeding comments..."
	local created=0 nc=${#comment_contents[@]}
	for ci in $(seq 0 $((${#post_contents[@]} - 1))); do
		local pid=${post_id[$ci]:-}
		[[ -n "$pid" ]] || continue
		local count=$((1 + ci % 2))
		for j in $(seq 0 $((count - 1))); do
			local ai=$(((ci + j + 1) % N))
			[[ "$ai" -eq "${post_author[$ci]}" ]] && ai=$(((ai + 1) % N))
			local ti; ti=$(tok "$ai")
			[[ -n "$ti" ]] || continue
			local content="${comment_contents[$(((ci + j) % nc))]}"
			post_json "/posts/$pid/comments" "$(jq -nc --arg c "$content" '{content:$c}')" "$ti" >/dev/null
			created=$((created + 1))
		done
	done
	echo "  + $created comments"
}

main() {
	wait_for_api
	seed_users
	seed_friendships
	seed_follows
	seed_posts
	seed_likes
	seed_comments
	echo ""
	echo "✅ Seeding finished via the HTTP API ($API)"
}

main "$@"
