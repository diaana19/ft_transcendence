#!/usr/bin/env bash
#
# seed.sh — populate a running ft_transcendence stack with realistic demo data
# using nothing but bash + curl (plus jq to build/parse JSON).
#
# What it creates:
#   - N users (default 500) with real names + photos pulled from randomuser.me
#     (an open demo-data API), profile bios and avatars set via the REST API
#   - POSTS_TARGET posts (default 2000) with unique, randomly-composed content
#   - likes and comments (the app's "replies") across them
#   - a social graph: follows + friend requests that get accepted
#
# Everything goes through curl against the public REST API. The content-creation
# phases fan out up to $PAR requests in parallel to stay fast at this scale.
#
# Usage:
#   scripts/seed.sh                       # 500 users, 2000 posts
#   USERS=200 POSTS_TARGET=800 scripts/seed.sh
#   PAR=20 scripts/seed.sh                # more parallelism
#   API=https://host:3000/api scripts/seed.sh
#
# Requirements: bash, curl, jq on the host; a running stack (`make up`).
set -euo pipefail

# ---------------------------------------------------------------------------
# Config (all overridable via environment)
# ---------------------------------------------------------------------------
API="${API:-https://localhost:3000/api}"      # public API base (through nginx)
USERS="${USERS:-500}"                         # how many users to create
POSTS_TARGET="${POSTS_TARGET:-2000}"          # total posts to create (unique content)
PAR="${PAR:-12}"                              # max concurrent curl requests
PASSWORD="${SEED_PASSWORD:-Tr@nscend42}"      # shared password (meets policy)
RU_SEED="${RU_SEED:-ftt}"                     # randomuser seed → reproducible people
EMAIL_DOMAIN="${EMAIL_DOMAIN:-seed.test}"     # avoids clashing with real domains

# Resolve the repo root so compose (-f infra/...) and relative paths work no
# matter where the script is invoked from.
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

# Container engine / compose command — same autodetection as the Makefile so the
# Redis rate-limit reset can find the stack.
ENGINE="${ENGINE:-$(command -v podman >/dev/null 2>&1 && echo podman || echo docker)}"
COMPOSE="${COMPOSE:-$ENGINE compose -f infra/docker-compose.yml}"

CURL=(curl -sk --max-time 30)   # -k: the proxy serves a self-signed cert locally

# ---------------------------------------------------------------------------
# Small helpers
# ---------------------------------------------------------------------------
c_reset=$'\033[0m'; c_dim=$'\033[2m'; c_grn=$'\033[32m'; c_yel=$'\033[33m'; c_red=$'\033[31m'; c_bld=$'\033[1m'
log()  { printf '%s\n' "${c_dim}· $*${c_reset}"; }
step() { printf '\n%s\n' "${c_bld}▸ $*${c_reset}"; }
ok()   { printf '%s\n' "${c_grn}✓ $*${c_reset}"; }
warn() { printf '%s\n' "${c_yel}! $*${c_reset}" >&2; }
die()  { printf '%s\n' "${c_red}✗ $*${c_reset}" >&2; exit 1; }

# api METHOD PATH [JSON_BODY] [TOKEN] → echoes the response body
api() {
  local method=$1 path=$2 body=${3:-} token=${4:-}
  local args=("${CURL[@]}" -X "$method" "$API$path" -H 'Content-Type: application/json')
  [ -n "$token" ] && args+=(-H "Authorization: Bearer $token")
  [ -n "$body" ]  && args+=(-d "$body")
  "${args[@]}"
}

# The /auth group is rate-limited to ~20 req/min per IP, and every request here
# shares the proxy's IP. Flushing the rate_limit:* keys between auth batches
# keeps seeding fast without tripping 429s. No-op if the redis service isn't
# reachable (e.g. seeding a remote API) — we fall back to a short sleep.
reset_rate_limit() {
  # </dev/null is critical: this runs inside the `while read` user loop, and
  # `compose exec` would otherwise swallow the loop's stdin and stop it early.
  if ! $COMPOSE exec -T redis sh -c "redis-cli --scan --pattern 'rate_limit:*' | xargs -r redis-cli del" </dev/null >/dev/null 2>&1; then
    sleep 3
  fi
}

# pick a random 0-based index < $1
rand() { echo $(( RANDOM % $1 )); }

# throttle blocks until fewer than $PAR background jobs are running, so the
# fan-out phases keep at most $PAR curl requests in flight at once.
throttle() { while [ "$(jobs -r | wc -l)" -ge "$PAR" ]; do wait -n; done; }

require() { command -v "$1" >/dev/null 2>&1 || die "missing required tool: $1"; }

# ---------------------------------------------------------------------------
# Content corpora for generated comments and bios
# ---------------------------------------------------------------------------
COMMENTS=(
  "This is great, thanks for sharing."
  "Totally agree with this."
  "Can you say more about how you did it?"
  "Saving this for later."
  "Haha so relatable."
  "Strong opinion, respect it."
  "I learned this the hard way too."
  "Underrated take honestly."
  "Count me in for Friday."
  "This made my day."
  "Bookmarked. Pure gold."
  "Wait this is brilliant."
)
BIOS=(
  "Builder of things. Breaker of builds."
  "42 student. Caffeine to code converter."
  "I make computers do tricks."
  "Backend by day, frontend by mistake."
  "Pong world champion in my own head."
  "Shipping bugs and the occasional feature."
  "Probably refactoring something right now."
  "Distributed systems, local coffee."
)

# Pull a random element from an array named by $1
pick() { local -n arr=$1; echo "${arr[$(rand ${#arr[@]})]}"; }

# Word banks for synthesizing post text. The combinatorial space across the
# templates below is in the hundreds of thousands, so 2000+ unique posts is easy.
PO_OPEN=(Just Finally Today Somehow Yesterday "This week I" "Once again I" "Pretty sure I" "Not gonna lie I" "Against all odds I")
PO_VERB=(shipped refactored broke debugged rewrote deployed optimized untangled benchmarked containerized documented deleted automated profiled migrated patched)
PO_VERBING=(shipping refactoring breaking debugging rewriting deploying optimizing untangling benchmarking documenting deleting automating profiling migrating patching fighting)
PO_NOUN=("the auth flow" "the WebSocket layer" "a flaky test" "the build pipeline" "the legacy module" "the rate limiter" "our Postgres queries" "the CI config" "a memory leak" "the cache layer" "the deploy script" "three race conditions" "the API gateway" "a gnarly regex" "the migration" "the login page" "the retry logic" "the seed script" "the docker setup" "the nginx config")
PO_TOPIC=("tabs vs spaces" "microservices" "rewriting it in Rust" "dark mode" "monorepos" "kubernetes" "serverless" "the new framework" "strict type safety" "test coverage" "code review" "pair programming" "daily standups" "tech debt" "the legacy codebase" "premature optimization")
PO_ADJ=(underrated overrated cursed elegant chaotic flawless fragile beautiful unhinged questionable inevitable suspicious)
PO_TAIL=("Feels good." "No idea how." "Ship it." "Tests are green." "Worth it." "Send coffee." "Small wins." "Never again." "Big if true." "Living dangerously." "10 out of 10." "Do not ask." "It compiles." "We move." "Chef kiss." "Pain.")

# gen_post prints a unique, randomly-composed post. It loops until it lands on a
# sentence not yet seen this run, so no two posts share the same text.
declare -A SEEN_POST
gen_post() {
  local s
  while :; do
    case $(( RANDOM % 4 )) in
      0) s="$(pick PO_OPEN) $(pick PO_VERB) $(pick PO_NOUN). $(pick PO_TAIL)" ;;
      1) s="Hot take: $(pick PO_TOPIC) is $(pick PO_ADJ). $(pick PO_TAIL)" ;;
      2) s="Honestly, $(pick PO_TOPIC) might be the most $(pick PO_ADJ) thing in tech. $(pick PO_TAIL)" ;;
      3) s="Day $(( RANDOM % 90 + 1 )) of $(pick PO_VERBING) $(pick PO_NOUN). $(pick PO_TAIL)" ;;
    esac
    if [ -z "${SEEN_POST[$s]:-}" ]; then SEEN_POST[$s]=1; printf '%s' "$s"; return; fi
  done
}

# ---------------------------------------------------------------------------
# Preflight
# ---------------------------------------------------------------------------
require curl; require jq
step "ft_transcendence seeder  ${c_dim}(api=$API · users=$USERS)${c_reset}"

# Wait for the API to come up — `make seed-clean` runs this right after `make up`,
# so the backend healthcheck may still be going green.
WAIT="${WAIT:-90}"
for attempt in $(seq 1 "$WAIT"); do
  if "${CURL[@]}" --max-time 5 "$API/health" | grep -q '"status"'; then
    [ "$attempt" -gt 1 ] && printf '%s\n' "$c_reset"
    ok "API is reachable"; break
  fi
  [ "$attempt" -eq "$WAIT" ] && die "API not reachable at $API after ${WAIT}s — is the stack up? (make up)"
  [ "$attempt" -eq 1 ] && printf '%s' "${c_dim}  waiting for API"
  printf '.'; sleep 1
done

# ---------------------------------------------------------------------------
# 1. Fetch real demo people (names + photos) from randomuser.me
# ---------------------------------------------------------------------------
step "Fetching $USERS demo profiles from randomuser.me"
people=$("${CURL[@]}" --max-time 30 \
  "https://randomuser.me/api/?results=${USERS}&inc=name,login,email,dob,picture&noinfo&seed=${RU_SEED}" \
  | jq -r '.results[] | [.login.username, .name.first, .name.last, (.dob.date[:10]), .picture.large] | @tsv' || true)

if [ -z "$people" ]; then
  warn "randomuser.me unavailable — falling back to synthetic profiles + pravatar photos"
  firsts=(Alex Sam Jordan Taylor Casey Riley Morgan Jamie Avery Quinn Drew Skyler Cameron Reese Devon)
  lasts=(Stone Reyes Khan Silva Park Nguyen Costa Walsh Mertz Frami Lubowitz Hane Cole Vega Ortiz)
  people=""
  for i in $(seq 1 "$USERS"); do
    f=${firsts[$((RANDOM % ${#firsts[@]}))]}; l=${lasts[$((RANDOM % ${#lasts[@]}))]}
    u="$(echo "$f$l$i" | tr 'A-Z' 'a-z')"
    img=$(( (i % 70) + 1 ))
    people+="$u\t$f\t$l\t1996-04-0$(( (i%9)+1 ))\thttps://i.pravatar.cc/300?img=$img\n"
  done
  people=$(printf '%b' "$people")
fi

# ---------------------------------------------------------------------------
# 2. Register + login + set profile (avatar/bio). Rate-limit aware.
# ---------------------------------------------------------------------------
step "Registering users and setting profiles"
declare -a UID_ARR=() TOK_ARR=() NAME_ARR=() UNAME_ARR=()
i=0; created=0
while IFS=$'\t' read -r -u 3 username first last dob avatar; do
  [ -z "$username" ] && continue
  # keep usernames API-valid + deterministic email domain to avoid clashes
  username=$(echo "$username" | tr -cd 'a-zA-Z0-9_')
  email="${username}@${EMAIL_DOMAIN}"
  name="$first $last"

  (( i % 8 == 0 )) && reset_rate_limit   # stay under 20 auth req/min

  reg=$(api POST /auth/register \
    "$(jq -n --arg u "$username" --arg e "$email" --arg p "$PASSWORD" --arg d "$dob" \
        '{username:$u,email:$e,password:$p,dateOfBirth:$d}')")
  uid=$(echo "$reg" | jq -r '.id // empty')

  # If we hit the limiter, flush and retry once.
  if [ -z "$uid" ] && echo "$reg" | grep -qi "too many"; then
    reset_rate_limit
    reg=$(api POST /auth/register \
      "$(jq -n --arg u "$username" --arg e "$email" --arg p "$PASSWORD" --arg d "$dob" \
          '{username:$u,email:$e,password:$p,dateOfBirth:$d}')")
    uid=$(echo "$reg" | jq -r '.id // empty')
  fi
  if [ -z "$uid" ]; then
    log "skip $username: $(echo "$reg" | jq -r '.error // "unknown error"')"
    i=$((i+1)); continue
  fi

  login=$(api POST /auth/login "$(jq -n --arg e "$email" --arg p "$PASSWORD" '{email:$e,password:$p}')")
  tok=$(echo "$login" | jq -r '.token // empty')
  [ -z "$tok" ] && { log "skip $username: login failed"; i=$((i+1)); continue; }

  # avatar + bio + display name (partial update — username/email untouched)
  api PUT "/users/$uid" \
    "$(jq -n --arg n "$name" --arg b "$(pick BIOS)" --arg a "$avatar" '{name:$n,bio:$b,avatar:$a}')" \
    "$tok" >/dev/null

  UID_ARR+=("$uid"); TOK_ARR+=("$tok"); NAME_ARR+=("$name"); UNAME_ARR+=("$username")
  created=$((created+1)); i=$((i+1))
  printf '\r  %s%d/%d%s users' "$c_dim" "$created" "$USERS" "$c_reset"
done 3<<< "$people"
printf '\n'
[ "$created" -gt 1 ] || die "could not create users (created=$created)"
ok "$created users created with photos + profiles"

N=${#UID_ARR[@]}

# ---------------------------------------------------------------------------
# 3. Posts — POSTS_TARGET posts with unique content, spread across all users.
#    Content is composed in this (single) shell so uniqueness holds; the HTTP
#    POSTs fan out up to $PAR at a time and report their ids back via a file.
# ---------------------------------------------------------------------------
step "Creating posts  ${c_dim}(target $POSTS_TARGET, unique content, $PAR-way parallel)${c_reset}"
declare -a POST_ID=() POST_OWNER=()
postsfile=$(mktemp)
for ((k = 0; k < POSTS_TARGET; k++)); do
  idx=$(rand "$N")
  body=$(jq -n --arg c "$(gen_post)" '{content:$c}')
  {
    pid=$(api POST /posts "$body" "${TOK_ARR[$idx]}" | jq -r '.id // empty')
    [ -n "$pid" ] && printf '%s\t%s\n' "$idx" "$pid" >> "$postsfile"
  } &
  throttle
  (( k % 200 == 0 )) && printf '\r  %s%d/%d%s queued' "$c_dim" "$k" "$POSTS_TARGET" "$c_reset"
done
wait
printf '\r%*s\r' 40 ''   # clear the progress line
while IFS=$'\t' read -r oidx pid; do POST_ID+=("$pid"); POST_OWNER+=("$oidx"); done < "$postsfile"
rm -f "$postsfile"
ok "${#POST_ID[@]} posts created"
[ "${#POST_ID[@]}" -gt 0 ] || die "no posts created — aborting"
P=${#POST_ID[@]}

# Count the lines a fan-out phase appended to a temp file (each = one success).
count_file() { wc -l < "$1" | tr -d ' '; }

# ---------------------------------------------------------------------------
# 4. Likes — each user likes a random handful of posts (parallel)
# ---------------------------------------------------------------------------
step "Adding likes  ${c_dim}($PAR-way parallel)${c_reset}"
lf=$(mktemp)
for idx in "${!UID_ARR[@]}"; do
  for _ in $(seq 1 $(( (RANDOM % 8) + 3 ))); do
    p=$(rand "$P")
    { api POST "/posts/${POST_ID[$p]}/like" "" "${TOK_ARR[$idx]}" >/dev/null && echo 1 >> "$lf"; } &
    throttle
  done
done
wait
ok "$(count_file "$lf") likes added"; rm -f "$lf"

# ---------------------------------------------------------------------------
# 5. Comments (the app's post replies, parallel)
# ---------------------------------------------------------------------------
step "Adding comments  ${c_dim}($PAR-way parallel)${c_reset}"
cf=$(mktemp)
for idx in "${!UID_ARR[@]}"; do
  for _ in $(seq 1 $(( (RANDOM % 4) + 2 ))); do
    p=$(rand "$P")
    body=$(jq -n --arg c "$(pick COMMENTS)" '{content:$c}')
    { api POST "/posts/${POST_ID[$p]}/comments" "$body" "${TOK_ARR[$idx]}" >/dev/null && echo 1 >> "$cf"; } &
    throttle
  done
done
wait
ok "$(count_file "$cf") comments added"; rm -f "$cf"

# ---------------------------------------------------------------------------
# 6. Social graph — follows + friend requests that get accepted (parallel)
# ---------------------------------------------------------------------------
step "Building the social graph  ${c_dim}($PAR-way parallel)${c_reset}"
ff=$(mktemp); frf=$(mktemp)
for idx in "${!UID_ARR[@]}"; do
  # follow 3-6 distinct random others
  for _ in $(seq 1 $(( (RANDOM % 4) + 3 ))); do
    t=$(rand "$N"); [ "$t" -eq "$idx" ] && continue
    { api POST "/friends/follow/${UID_ARR[$t]}" "" "${TOK_ARR[$idx]}" >/dev/null && echo 1 >> "$ff"; } &
    throttle
  done
  # send 1-2 friend requests, then have the target accept
  for _ in $(seq 1 $(( (RANDOM % 2) + 1 ))); do
    t=$(rand "$N"); [ "$t" -eq "$idx" ] && continue
    {
      req=$(api POST "/friends/request/${UID_ARR[$t]}" "" "${TOK_ARR[$idx]}")
      if echo "$req" | grep -qi '"message"'; then
        api POST "/friends/accept/${UID_ARR[$idx]}" "" "${TOK_ARR[$t]}" >/dev/null && echo 1 >> "$frf"
      fi
    } &
    throttle
  done
done
wait
ok "$(count_file "$ff") follows · $(count_file "$frf") friendships"; rm -f "$ff" "$frf"

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------
printf '\n%s\n' "${c_grn}${c_bld}Seed complete.${c_reset}"
printf '%s\n' "${c_dim}Log in with any seeded account — e.g. ${UNAME_ARR[0]}@${EMAIL_DOMAIN} / ${PASSWORD}${c_reset}"
printf '%s\n' "${c_dim}Open the app: https://localhost:3000${c_reset}"
