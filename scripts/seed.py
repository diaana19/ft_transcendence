#!/usr/bin/env python3
"""
seed.py — populate a running ft_transcendence stack with realistic demo data
using nothing but the Python standard library (urllib). No third-party packages
required, and it talks only to the HTTP API — no container engine access.

What it creates:
  - N users (default 50) with real names + photos pulled from randomuser.me
    (an open demo-data API), profile bios and avatars set via the REST API
  - POSTS_TARGET posts (default 500) with unique, randomly-composed content
  - likes and comments (the app's "replies") across them
  - a social graph: follows + friend requests that get accepted

Everything goes through HTTP against the public REST API. The content-creation
phases fan out up to $PAR requests in parallel (thread pool) to stay fast at
this scale.

The /auth endpoints are rate-limited to RATE_LIMIT_MAX requests/min per IP
(default 20, matching the backend). The seeder paces its own register/login
calls to stay under that ceiling, so it never trips a 429 — set RATE_LIMIT_MAX
to the same value the backend runs with to seed at full speed.

Usage:
  scripts/seed.py                       # 50 users, 500 posts
  USERS=200 POSTS_TARGET=800 scripts/seed.py
  PAR=20 scripts/seed.py                # more parallelism
  RATE_LIMIT_MAX=1000 scripts/seed.py   # match a backend with a raised limit
  API=https://host:3000/api scripts/seed.py

Requirements: python3 on the host; a running stack (`make up`).
"""

import json
import os
import random
import ssl
import sys
import time
import urllib.error
import urllib.request
from collections import deque
from concurrent.futures import ThreadPoolExecutor
from pathlib import Path
from threading import Lock

# ---------------------------------------------------------------------------
# Config (all overridable via environment)
# ---------------------------------------------------------------------------
API = os.environ.get("API", "https://localhost:3000/api")  # public API base (through nginx)
USERS = int(os.environ["USERS"])                           # how many users to create
POSTS_TARGET = int(os.environ["POSTS_TARGET"])             # total posts to create (unique content)
PAR = int(os.environ["PAR"])                               # max concurrent requests
PASSWORD = os.environ.get("SEED_PASSWORD", "Tr@nscend42")  # shared password (meets policy)
RU_SEED = os.environ.get("RU_SEED", "ftt")                 # randomuser seed → reproducible people
EMAIL_DOMAIN = os.environ.get("EMAIL_DOMAIN", "seed.test") # avoids clashing with real domains
RATE_LIMIT_MAX = int(os.environ["RATE_LIMIT_MAX"])         # backend /auth ceiling, req/min/IP

# Resolve the repo root so relative paths work no matter where we're invoked from.
ROOT = Path(__file__).resolve().parent.parent
os.chdir(ROOT)

# The proxy serves a self-signed cert locally, so skip TLS verification (curl -k).
_SSL_CTX = ssl.create_default_context()
_SSL_CTX.check_hostname = False
_SSL_CTX.verify_mode = ssl.CERT_NONE

# ---------------------------------------------------------------------------
# Small helpers
# ---------------------------------------------------------------------------
c_reset = "\033[0m"; c_dim = "\033[2m"; c_grn = "\033[32m"
c_yel = "\033[33m"; c_red = "\033[31m"; c_bld = "\033[1m"


def log(msg):  print(f"{c_dim}· {msg}{c_reset}")
def step(msg): print(f"\n{c_bld}▸ {msg}{c_reset}")
def ok(msg):   print(f"{c_grn}✓ {msg}{c_reset}")
def warn(msg): print(f"{c_yel}! {msg}{c_reset}", file=sys.stderr)


def die(msg):
    print(f"{c_red}✗ {msg}{c_reset}", file=sys.stderr)
    sys.exit(1)


def http(method, url, body=None, token=None, timeout=30):
    """Perform an HTTP request, returning (status, text). Never raises on HTTP
    errors — mirrors curl echoing the response body regardless of status."""
    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    data = body.encode() if isinstance(body, str) else body
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=timeout, context=_SSL_CTX) as resp:
            return resp.status, resp.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8", "replace")
    except (urllib.error.URLError, TimeoutError, OSError):
        return 0, ""


def api(method, path, body=None, token=None):
    """api METHOD PATH [JSON_BODY] [TOKEN] → returns the parsed JSON (or {})."""
    _, text = http(method, f"{API}{path}", body=body, token=token)
    try:
        return json.loads(text) if text else {}
    except json.JSONDecodeError:
        return {"_raw": text}


class AuthThrottle:
    """Paces auth-group requests (register/login) so we never exceed the
    backend's per-IP ceiling of RATE_LIMIT_MAX requests/minute. Every seed
    request shares the proxy's IP, so without pacing a 500-user run would trip
    429s almost immediately. A trailing 60s sliding window of at most
    `per_minute` calls is stricter than the server's fixed window, so it's safe.

    Set RATE_LIMIT_MAX to the value the backend actually runs with: too low just
    seeds slower; too high risks 429s (caught and retried below as a backstop)."""

    def __init__(self, per_minute):
        self.per_minute = max(1, per_minute)
        self.window = deque()
        self.lock = Lock()

    def wait(self):
        with self.lock:
            now = time.monotonic()
            while self.window and now - self.window[0] >= 60:
                self.window.popleft()
            if len(self.window) >= self.per_minute:
                sleep_for = 60 - (now - self.window[0])
                if sleep_for > 0:
                    time.sleep(sleep_for)
                now = time.monotonic()
                while self.window and now - self.window[0] >= 60:
                    self.window.popleft()
            self.window.append(now)


_auth_throttle = AuthThrottle(RATE_LIMIT_MAX)


def auth_api(method, path, body=None, token=None):
    """Like api(), but paced to respect the backend's /auth rate limit."""
    _auth_throttle.wait()
    return api(method, path, body, token)


# ---------------------------------------------------------------------------
# Content corpora for generated comments and bios
# ---------------------------------------------------------------------------
COMMENTS = [
    "This is great, thanks for sharing.",
    "Totally agree with this.",
    "Can you say more about how you did it?",
    "Saving this for later.",
    "Haha so relatable.",
    "Strong opinion, respect it.",
    "I learned this the hard way too.",
    "Underrated take honestly.",
    "Count me in for Friday.",
    "This made my day.",
    "Bookmarked. Pure gold.",
    "Wait this is brilliant.",
]
BIOS = [
    "Builder of things. Breaker of builds.",
    "42 student. Caffeine to code converter.",
    "I make computers do tricks.",
    "Backend by day, frontend by mistake.",
    "Pong world champion in my own head.",
    "Shipping bugs and the occasional feature.",
    "Probably refactoring something right now.",
    "Distributed systems, local coffee.",
]

# Word banks for synthesizing post text. The combinatorial space across the
# templates below is in the hundreds of thousands, so 2000+ unique posts is easy.
PO_OPEN = ["Just", "Finally", "Today", "Somehow", "Yesterday", "This week I",
           "Once again I", "Pretty sure I", "Not gonna lie I", "Against all odds I"]
PO_VERB = ["shipped", "refactored", "broke", "debugged", "rewrote", "deployed",
           "optimized", "untangled", "benchmarked", "containerized", "documented",
           "deleted", "automated", "profiled", "migrated", "patched"]
PO_VERBING = ["shipping", "refactoring", "breaking", "debugging", "rewriting",
              "deploying", "optimizing", "untangling", "benchmarking", "documenting",
              "deleting", "automating", "profiling", "migrating", "patching", "fighting"]
PO_NOUN = ["the auth flow", "the WebSocket layer", "a flaky test", "the build pipeline",
           "the legacy module", "the rate limiter", "our Postgres queries", "the CI config",
           "a memory leak", "the cache layer", "the deploy script", "three race conditions",
           "the API gateway", "a gnarly regex", "the migration", "the login page",
           "the retry logic", "the seed script", "the docker setup", "the nginx config"]
PO_TOPIC = ["tabs vs spaces", "microservices", "rewriting it in Rust", "dark mode",
            "monorepos", "kubernetes", "serverless", "the new framework",
            "strict type safety", "test coverage", "code review", "pair programming",
            "daily standups", "tech debt", "the legacy codebase", "premature optimization"]
PO_ADJ = ["underrated", "overrated", "cursed", "elegant", "chaotic", "flawless",
          "fragile", "beautiful", "unhinged", "questionable", "inevitable", "suspicious"]
PO_TAIL = ["Feels good.", "No idea how.", "Ship it.", "Tests are green.", "Worth it.",
           "Send coffee.", "Small wins.", "Never again.", "Big if true.",
           "Living dangerously.", "10 out of 10.", "Do not ask.", "It compiles.",
           "We move.", "Chef kiss.", "Pain."]
# Hashtags appended to most posts so the trends feature has data. The pool is
# deliberately small-ish so some tags recur often and a clear ranking emerges.
PO_TAGS = ["#golang", "#rust", "#docker", "#kubernetes", "#postgres", "#redis",
           "#webdev", "#devops", "#testing", "#opensource", "#typescript",
           "#react", "#ci", "#42born2code", "#cleancode", "#performance"]


def pick(arr):
    """Pull a random element from a list."""
    return random.choice(arr)


def gen_tags():
    """Return a leading-space string of 0-2 distinct hashtags. Most posts get at
    least one so the trends aggregation has something to rank."""
    n = random.choices([0, 1, 2], weights=[2, 5, 3])[0]
    if n == 0:
        return ""
    return " " + " ".join(random.sample(PO_TAGS, n))


# gen_post returns a unique, randomly-composed post. It loops until it lands on a
# sentence not yet seen this run, so no two posts share the same text.
_SEEN_POST = set()


def gen_post():
    while True:
        kind = random.randint(0, 3)
        if kind == 0:
            s = f"{pick(PO_OPEN)} {pick(PO_VERB)} {pick(PO_NOUN)}. {pick(PO_TAIL)}"
        elif kind == 1:
            s = f"Hot take: {pick(PO_TOPIC)} is {pick(PO_ADJ)}. {pick(PO_TAIL)}"
        elif kind == 2:
            s = (f"Honestly, {pick(PO_TOPIC)} might be the most {pick(PO_ADJ)} "
                 f"thing in tech. {pick(PO_TAIL)}")
        else:
            s = (f"Day {random.randint(1, 90)} of {pick(PO_VERBING)} "
                 f"{pick(PO_NOUN)}. {pick(PO_TAIL)}")
        s += gen_tags()
        if s not in _SEEN_POST:
            _SEEN_POST.add(s)
            return s


# ---------------------------------------------------------------------------
# Preflight
# ---------------------------------------------------------------------------
def wait_for_api():
    step(f"ft_transcendence seeder  {c_dim}(api={API} · users={USERS}){c_reset}")
    # Wait for the API to come up — `make seed-clean` runs this right after
    # `make up`, so the backend healthcheck may still be going green.
    wait = int(os.environ.get("WAIT", "90"))
    for attempt in range(1, wait + 1):
        _, text = http("GET", f"{API}/health", timeout=5)
        if '"status"' in text:
            if attempt > 1:
                print(c_reset)
            ok("API is reachable")
            return
        if attempt == 1:
            sys.stdout.write(f"{c_dim}  waiting for API")
        sys.stdout.write(".")
        sys.stdout.flush()
        time.sleep(1)
    die(f"API not reachable at {API} after {wait}s — is the stack up? (make up)")


# ---------------------------------------------------------------------------
# 1. Fetch real demo people (names + photos) from randomuser.me
# ---------------------------------------------------------------------------
def fetch_people():
    step(f"Fetching {USERS} demo profiles from randomuser.me")
    url = (f"https://randomuser.me/api/?results={USERS}"
           f"&inc=name,login,email,dob,picture&noinfo&seed={RU_SEED}")
    status, text = http("GET", url, timeout=30)
    people = []
    if status == 200 and text:
        try:
            for r in json.loads(text).get("results", []):
                people.append((
                    r["login"]["username"],
                    r["name"]["first"],
                    r["name"]["last"],
                    r["dob"]["date"][:10],
                    r["picture"]["large"],
                ))
        except (json.JSONDecodeError, KeyError, TypeError):
            people = []

    if not people:
        warn("randomuser.me unavailable — falling back to synthetic profiles + pravatar photos")
        firsts = ["Alex", "Sam", "Jordan", "Taylor", "Casey", "Riley", "Morgan",
                  "Jamie", "Avery", "Quinn", "Drew", "Skyler", "Cameron", "Reese", "Devon"]
        lasts = ["Stone", "Reyes", "Khan", "Silva", "Park", "Nguyen", "Costa", "Walsh",
                 "Mertz", "Frami", "Lubowitz", "Hane", "Cole", "Vega", "Ortiz"]
        for i in range(1, USERS + 1):
            f = random.choice(firsts)
            l = random.choice(lasts)
            u = f"{f}{l}{i}".lower()
            img = (i % 70) + 1
            people.append((u, f, l, f"1996-04-0{(i % 9) + 1}",
                           f"https://i.pravatar.cc/300?img={img}"))
    return people


# ---------------------------------------------------------------------------
# 2. Register + login + set profile (avatar/bio). Rate-limit aware.
# ---------------------------------------------------------------------------
def seed_users(people):
    step("Registering users and setting profiles")
    uids, toks, names, unames = [], [], [], []
    created = 0
    for i, (username, first, last, dob, avatar) in enumerate(people):
        if not username:
            continue
        # keep usernames API-valid (GitHub rules: alphanumeric only here, which
        # always satisfies them) + deterministic email domain to avoid clashes
        username = "".join(ch for ch in username if ch.isalnum())
        email = f"{username}@{EMAIL_DOMAIN}"
        name = f"{first} {last}"

        reg_body = json.dumps({"username": username, "email": email,
                               "password": PASSWORD, "dateOfBirth": dob})
        reg = auth_api("POST", "/auth/register", reg_body)
        uid = reg.get("id")

        # Backstop: if the throttle and the backend disagree on the ceiling, wait
        # out a full window and retry once before giving up on this user.
        if not uid and "too many" in json.dumps(reg).lower():
            time.sleep(60)
            reg = auth_api("POST", "/auth/register", reg_body)
            uid = reg.get("id")
        if not uid:
            log(f"skip {username}: {reg.get('error', 'unknown error')}")
            continue

        login = auth_api("POST", "/auth/login",
                         json.dumps({"email": email, "password": PASSWORD}))
        tok = login.get("token")
        if not tok:
            log(f"skip {username}: login failed")
            continue

        # avatar + bio + display name (partial update — username/email untouched)
        api("PUT", f"/users/{uid}",
            json.dumps({"displayname": name, "bio": pick(BIOS), "avatar": avatar}), tok)

        uids.append(uid); toks.append(tok); names.append(name); unames.append(username)
        created += 1
        sys.stdout.write(f"\r  {c_dim}{created}/{USERS}{c_reset} users")
        sys.stdout.flush()
    print()
    if created <= 1:
        die(f"could not create users (created={created})")
    ok(f"{created} users created with photos + profiles")
    return uids, toks, names, unames


# ---------------------------------------------------------------------------
# Parallel fan-out helper — run `fn` over `tasks` with up to PAR workers,
# returning the list of truthy results.
# ---------------------------------------------------------------------------
def fan_out(tasks, fn):
    results = []
    with ThreadPoolExecutor(max_workers=PAR) as pool:
        for r in pool.map(fn, tasks):
            if r:
                results.append(r)
    return results


# ---------------------------------------------------------------------------
# 3. Posts — POSTS_TARGET posts with unique content, spread across all users.
#    Content is composed up-front (single-threaded) so uniqueness holds; the
#    HTTP POSTs fan out up to PAR at a time.
# ---------------------------------------------------------------------------
def seed_posts(uids, toks):
    n = len(uids)
    step(f"Creating posts  {c_dim}(target {POSTS_TARGET}, unique content, "
         f"{PAR}-way parallel){c_reset}")
    # Pre-compose all (owner_idx, content) pairs single-threaded for uniqueness.
    jobs = [(random.randrange(n), gen_post()) for _ in range(POSTS_TARGET)]

    def make_post(job):
        idx, content = job
        res = api("POST", "/posts", json.dumps({"content": content}), toks[idx])
        pid = res.get("id")
        return (idx, pid) if pid else None

    created = fan_out(jobs, make_post)
    post_owner = [idx for idx, _ in created]
    post_id = [pid for _, pid in created]
    ok(f"{len(post_id)} posts created")
    if not post_id:
        die("no posts created — aborting")
    return post_id, post_owner


# ---------------------------------------------------------------------------
# 4. Likes — each user likes a random handful of posts (parallel)
# ---------------------------------------------------------------------------
def seed_likes(uids, toks, post_id):
    p = len(post_id)
    step(f"Adding likes  {c_dim}({PAR}-way parallel){c_reset}")
    tasks = []
    for idx in range(len(uids)):
        for _ in range(random.randint(3, 10)):
            tasks.append((idx, random.randrange(p)))

    def like(task):
        idx, pi = task
        status, _ = http(
            "POST",
            f"{API}/posts/{post_id[pi]}/react",
            body='{"value": 1}',
            token=toks[idx],
        )
        return 1 if status == 200 else None

    n = len(fan_out(tasks, like))
    ok(f"{n} likes added")


# ---------------------------------------------------------------------------
# 5. Comments (the app's post replies, parallel)
# ---------------------------------------------------------------------------
def seed_comments(uids, toks, post_id):
    p = len(post_id)
    step(f"Adding comments  {c_dim}({PAR}-way parallel){c_reset}")
    tasks = []
    for idx in range(len(uids)):
        for _ in range(random.randint(2, 5)):
            tasks.append((idx, random.randrange(p), pick(COMMENTS)))

    def comment(task):
        idx, pi, content = task
        res = api("POST", f"/posts/{post_id[pi]}/comments",
                  json.dumps({"content": content}), toks[idx])
        return 1 if res else None

    n = len(fan_out(tasks, comment))
    ok(f"{n} comments added")


# ---------------------------------------------------------------------------
# 6. Social graph — follows + friend requests that get accepted (parallel)
# ---------------------------------------------------------------------------
def seed_social(uids, toks):
    n = len(uids)
    step(f"Building the social graph  {c_dim}({PAR}-way parallel){c_reset}")

    follow_tasks, friend_tasks = [], []
    for idx in range(n):
        # follow 3-6 distinct random others
        for _ in range(random.randint(3, 6)):
            t = random.randrange(n)
            if t != idx:
                follow_tasks.append((idx, t))
        # send 1-2 friend requests, then have the target accept
        for _ in range(random.randint(1, 2)):
            t = random.randrange(n)
            if t != idx:
                friend_tasks.append((idx, t))

    def follow(task):
        idx, t = task
        status, _ = http("POST", f"{API}/friends/follow/{uids[t]}", token=toks[idx])
        return 1 if status else None

    def friend(task):
        idx, t = task
        req = api("POST", f"/friends/request/{uids[t]}", token=toks[idx])
        if "message" in req:
            status, _ = http("POST", f"{API}/friends/accept/{uids[idx]}", token=toks[t])
            return 1 if status else None
        return None

    follows = len(fan_out(follow_tasks, follow))
    friendships = len(fan_out(friend_tasks, friend))
    ok(f"{follows} follows · {friendships} friendships")


# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------
def main():
    wait_for_api()
    people = fetch_people()
    uids, toks, names, unames = seed_users(people)
    post_id, post_owner = seed_posts(uids, toks)
    seed_likes(uids, toks, post_id)
    seed_comments(uids, toks, post_id)
    seed_social(uids, toks)

    print(f"\n{c_grn}{c_bld}Seed complete.{c_reset}")
    print(f"{c_dim}Log in with any seeded account — e.g. "
          f"{unames[0]}@{EMAIL_DOMAIN} / {PASSWORD}{c_reset}")
    print(f"{c_dim}Open the app: https://localhost:3000{c_reset}")


if __name__ == "__main__":
    main()
