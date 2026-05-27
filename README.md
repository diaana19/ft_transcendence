*This project has been created as part of the 42 curriculum by luluzuri, vali, lepereir, rmarcas-, and dirituay.*

# ft_transcendence — a real-time social network

> ft_transcendence (subject v21.1, the "Surprise." open-ended edition). A multi-user
> social web application: users register, build a profile, follow and befriend each other,
> publish posts with images, like and comment, exchange real-time private messages, and
> receive notifications — secured with JWT, optional 2FA, and GitHub OAuth, and shipped as
> an installable PWA.

---

## Table of contents

- [Description](#description)
- [Team Information](#team-information)
- [Project Management](#project-management)
- [Technical Stack](#technical-stack)
- [Database Schema](#database-schema)
- [Features List](#features-list)
- [Modules](#modules)
- [Individual Contributions](#individual-contributions)
- [Instructions](#instructions)
- [Resources](#resources)
- [Known Limitations & Before Evaluation](#known-limitations--before-evaluation)

---

## Description

**ft_transcendence** is a social-network web application. Our goal was to build a real,
multi-user product that demonstrates a clean layered backend, real-time messaging, and a
modern reactive frontend — the "Social Network" project archetype from the subject (§V.3).

**Key features**

- **Accounts & security** — email/password sign-up with hashed + salted passwords (bcrypt),
  JWT access/refresh tokens, optional **TOTP 2FA**, and **GitHub OAuth 2.0** login.
- **Social graph** — follow/unfollow and friend requests (send / accept / reject / remove),
  with followers / following / friends lists.
- **Content** — create/read/update/delete posts with image attachments, likes, and threaded
  comments.
- **Real-time chat** — WebSocket-based direct messages and rooms with message history,
  threaded replies, and file attachments (backend complete; see limitations).
- **Notifications** — in-app notifications for social actions (friend requests, likes,
  messages, replies).
- **Search** — query users, posts, and messages with type filters, sorting, and pagination.
- **File management** — secure uploads (avatars, wallpapers, post media) with per-file
  visibility and access control.
- **Privacy** — GDPR data export and account deletion.
- **PWA** — installable, offline-capable progressive web app.

The application is a three-service stack (React frontend, Go API, PostgreSQL) plus Redis and
an nginx reverse proxy, all orchestrated with Docker Compose and started with a single
command.

---

## Team Information

The team is composed of 5 members. Roles follow the subject's recommended structure (§II.1).

| Member (42 login) | Role | Responsibilities |
|---|---|---|
| **dirituay** | Product Owner (PO) | Owns the product vision and backlog, prioritises features, validates completed work, and represents user needs to the team. |
| **rmarcas-** | Project Manager / Scrum Master | Organises planning and meetings, tracks progress and deadlines, manages blockers, and ensures team communication. Also led the **frontend** (UI, notifications, PWA). |
| **lepereir** | Tech Lead / Architect | Owns technical architecture and code quality: the golang-standards project layout, Go module structure, CI pipeline, linting, the integration test suite, and critical-code review. |
| **luluzuri** | Developer | Backend feature development (users, friends, posts, uploads, auth) and infrastructure/seed tooling. |
| **vali** | Developer | Real-time backend: WebSocket chat, message persistence, and search. |

> All members contributed to both the mandatory part and the modules. See
> [Individual Contributions](#individual-contributions) for details.

---

## Project Management

- **Roles & organisation** — work was split along the technical stack: a backend track
  (luluzuri, vali) and a frontend track (rmarcas-, dirituay), with the Tech Lead (lepereir)
  responsible for cross-cutting concerns (project layout, CI, tests, linting). Features were
  developed on dedicated branches (e.g. `S5-*` story branches) and merged via pull requests
  with review.
- **Version control** — Git with feature/story branches and reviewed pull requests; commits
  are SSH-signed. Work distribution is visible in `git shortlog -sne`.
- **Task tracking & communication** — _team to confirm the exact tools used_ (e.g. GitHub
  Issues / a project board for task tracking, and Discord for day-to-day communication).

---

## Technical Stack

| Layer | Technology | Why |
|---|---|---|
| **Frontend** | React 18 + Vite 5 + Tailwind CSS 4 | React is treated as a framework by the subject (ecosystem + architecture); Vite gives a fast dev server and an easy PWA pipeline (`vite-plugin-pwa`); Tailwind is the chosen styling solution. Routing via `react-router-dom`, HTTP via `axios`. |
| **Backend** | Go 1.25 + Gin | Gin is a mature, high-performance HTTP framework. Go's static typing, goroutines, and standard tooling suit a concurrent real-time API (WebSocket hub + REST). |
| **ORM** | GORM (Postgres driver) | Clean model definitions, auto-migration, soft deletes, and relations — satisfies the ORM minor module. |
| **Database** | PostgreSQL 15 | Relational data with clear foreign-key relationships (users, posts, friendships, messages…), UUID primary keys, and strong consistency for concurrent multi-user writes. |
| **Cache / pub-sub** | Redis 7 | Token/session helpers, per-IP rate-limit buckets, and pub/sub for broadcasting real-time notifications and chat across connections. |
| **Auth** | `golang-jwt/jwt/v5`, `pquerna/otp` (TOTP), `golang.org/x/crypto` (bcrypt), GitHub OAuth 2.0 | Standard, audited libraries for JWT, RFC-6238 TOTP, and password hashing. |
| **Real-time** | `gorilla/websocket` | Battle-tested WebSocket implementation for the chat hub. |
| **Reverse proxy** | nginx | Single public entrypoint: serves the frontend and proxies `/api` → backend. |
| **Containerisation** | Docker + Docker Compose | One-command, reproducible multi-service deployment. |
| **Testing / CI** | testcontainers-go, GitHub Actions, golangci-lint | Real-Postgres integration tests, a coverage gate, and static analysis. |

**Significant libraries:** `gin-contrib/cors`, `redis/go-redis/v9`, `gorm.io/driver/sqlite`
(test fallback), `vite-plugin-pwa`, `@heroicons/react` + `lucide-react` (icons).

---

## Database Schema

PostgreSQL with UUID primary keys. The canonical `users` table is created by
`database/init.sql`; all tables are kept in sync by GORM `AutoMigrate` from
`backend/internal/models/`. Soft deletes (`deleted_at`) are used on user-facing content.

```
users ───────────────┐
  id (uuid, PK)       │ 1
  github_id, provider │
  name                │
  username  (unique)  │        ┌──────────────── friends ────────────────┐
  email     (unique)  │        │ id (uuid, PK)                            │
  password  (hashed)  ├───────<│ user_id   (FK → users.id)                │  friendship +
  date_of_birth       │        │ friend_id (FK → users.id)                │  follow edges
  two_fa_secret       │        │ status  (pending | accepted | following) │
  two_fa_enabled      │        └──────────────────────────────────────────┘
  avatar, wallpaper   │
  bio                 │ 1      ┌──────────────── posts ───────────────────┐
  created/updated_at  ├───────<│ id (uuid, PK)                            │
  deleted_at          │        │ author_id (FK → users.id)                │
                      │        │ content, media_url                       │
                      │        │ likes_count, comments_count              │
                      │        │ created/updated_at, deleted_at           │
                      │        └───┬───────────────┬───────────────┬──────┘
                      │           <│              <│              <│
                      │     ┌── likes ──┐   ┌─ replies ──┐   ┌─ reposts ─┐
                      │     │ user_id   │   │ post_id    │   │ post_id   │
                      │     │ post_id   │   │ author_id  │   │ author_id │
                      │     │ UNIQUE    │   │ content    │   │           │
                      │     │(user,post)│   │ (comments) │   └───────────┘
                      │     └───────────┘   └────────────┘
                      │ 1
                      ├───< messages       (id, sender_id, recipient_id, content,
                      │                      file_id, room_id, parent_id → threaded
                      │                      replies, created_at)
                      │ 1
                      ├───< notifications  (id, user_id → recipient, actor_id, type
                      │                      [friend_request|like|message|reply],
                      │                      read, created_at)
                      │ 1
                      └───< files          (id, owner_id, path, filename, mime_type,
                                            size, visibility [public|private])
                                                  │ 1
                                                  └───< file_access (file_id, user_id)  -- per-file ACL
```

**Relationships (summary)**

- A **user** has many posts, likes, comments (`replies`), reposts, sent/received **messages**,
  **notifications**, and **files**.
- **friends** stores both friendship and follow edges as directed `(user_id → friend_id)`
  rows distinguished by `status`.
- **likes** has a unique `(user_id, post_id)` constraint (one like per user per post).
- **messages** support both direct messages (`recipient_id`) and rooms (`room_id`), plus
  threaded replies (`parent_id`) and optional attachments (`file_id`).
- **file_access** is a join table implementing per-file access control for private files.

---

## Features List

Attribution is derived from Git history (commit topics and per-file authorship); see
[Individual Contributions](#individual-contributions).

| Feature | Description | Main contributor(s) |
|---|---|---|
| Registration & login | Email/password with bcrypt hashing, input validation (FE + BE) | luluzuri, dirituay, rmarcas- |
| JWT auth + refresh + logout | Access/refresh tokens, Redis-assisted session/invalidation | luluzuri, vali |
| Two-Factor Auth (TOTP) | Setup / enable / disable / verify, with authenticator QR | luluzuri, vali (BE) · rmarcas-, dirituay (UI) |
| GitHub OAuth 2.0 | OAuth login + callback, account linking | vali, luluzuri |
| Profiles | View/edit profile, avatar + wallpaper, follower/following counts | rmarcas-, luluzuri |
| Friends & follows | Request/accept/reject/remove friends; follow/unfollow; lists | luluzuri, vali (BE) · rmarcas- (UI) |
| Posts | Create/read/update/delete posts with image attachments | luluzuri, rmarcas-, dirituay |
| Likes & comments | Toggle likes; threaded comments CRUD | rmarcas-, luluzuri, vali |
| Real-time chat | WebSocket DMs & rooms, history, threaded replies, attachments | **vali** (backend) |
| Notifications | In-app notifications for social actions; mark-all-read | rmarcas-, vali |
| Search | Users / posts / messages with filters, sorting, pagination | vali (BE) · lepereir (wiring + e2e tests) |
| File upload & management | Multi-type uploads, validation, visibility/ACL, serving | luluzuri, lepereir |
| GDPR | Data export + account deletion | lepereir |
| PWA | Installable + offline service worker | rmarcas- |
| Tests, CI & code quality | ~81% backend coverage, CI coverage gate, golangci-lint, project layout | **lepereir** |
| Infra & seed | Docker Compose, nginx, DB seeding (Go + curl scripts) | luluzuri, lepereir |

---

## Modules

The mandatory minimum is **14 points** (Major = 2, Minor = 1). The bonus part is capped at
**5 points** and counts only once all 14 mandatory points are fully functional — so the
maximum gradeable total is **14 + 5 = 19 points (125%)**, which is our target. We claim **20
points** (a 1-point buffer, since a module that fails to validate at the defense scores 0).
Per the subject, we list our full intended module set with an **honest implementation
status**, since only fully functional modules count at evaluation:

- ✅ **Done** — works end-to-end (UI + API).
- 🟧 **Backend complete / frontend pending** — the API and tests exist, but the feature is not
  yet wired into the UI (or is partially wired). See
  [Known Limitations](#known-limitations--before-evaluation).
- ⬜ **Planned** — targeted for the 125% bonus; not yet started.

### Major modules (2 pts each)

| Module | Category | Pts | Status | How it was implemented | Owner |
|---|---|---|---|---|---|
| Use a framework for frontend **and** backend | Web | 2 | ✅ | React (Vite) frontend + Gin (Go) backend | whole team |
| Real-time features via WebSockets | Web | 2 | 🟧 | `gorilla/websocket` hub with broadcast, connect/disconnect handling, Redis pub/sub; chat UI not yet wired | vali |
| User interaction (chat + profile + friends) | Web | 2 | 🟧 | Profiles & friends done end-to-end; chat backend complete, chat UI pending | luluzuri, vali, rmarcas- |
| Public API (rate-limit, docs, ≥5 CRUD endpoints) | Web | 2 | 🟧 | 30+ REST endpoints (GET/POST/PUT/DELETE), per-IP token-bucket rate limiting, Swagger + Postman docs. _Auth is JWT bearer rather than a separate API key._ | luluzuri, lepereir |
| Standard user management | User Mgmt | 2 | 🟧 | Profile edit, avatar (with default), friends, profile page. _Online-status indicator not yet implemented._ | luluzuri, rmarcas- |

**Major subtotal: 10 pts**

### Minor modules (1 pt each)

| Module | Category | Pts | Status | How it was implemented | Owner |
|---|---|---|---|---|---|
| ORM | Web | 1 | ✅ | GORM with auto-migration, relations, soft deletes | luluzuri |
| Notification system | Web | 1 | 🟧 | Notifications for friend requests/likes/messages/replies; basic UI + mark-all-read | rmarcas-, vali |
| Advanced search (filters/sort/pagination) | Web | 1 | 🟧 | `/api/search` over users/posts/messages, e2e-tested; search UI pending | vali, lepereir |
| File upload & management | Web | 1 | ✅ | Multi-type uploads, client+server validation, visibility/ACL, delete, serving | luluzuri, lepereir |
| Progressive Web App (PWA) | Web | 1 | ✅ | `vite-plugin-pwa`, manifest, Workbox service worker, installable + offline | rmarcas- |
| OAuth 2.0 (GitHub) | User Mgmt | 1 | ✅ | OAuth login + callback, account linking | vali, luluzuri |
| Two-Factor Authentication (TOTP) | User Mgmt | 1 | ✅ | RFC-6238 TOTP setup/enable/disable/verify with authenticator apps | luluzuri, vali |
| GDPR compliance | Data & Analytics | 1 | 🟧 | Data export + account deletion with confirmation; confirmation emails pending | lepereir |
| Multiple languages (i18n) | Accessibility & i18n | 1 | ⬜ | Planned: `react-i18next`, 3 languages + UI language switcher, all user-facing text translatable | rmarcas-, dirituay |
| Data export / import | Data & Analytics | 1 | ⬜ | Planned: extend the GDPR export with **import** + multiple formats (JSON/CSV/XML) and validation | lepereir |

**Minor subtotal: 10 pts**

### Point calculation

```
Major modules :  5 × 2 = 10 pts
Minor modules : 10 × 1 = 10 pts
                       -------
Total claimed          = 20 pts   (mandatory: 14 · bonus cap: 19 = 125%)
```

We target the **125% ceiling (19 points = 14 mandatory + 5 bonus)** and claim **20** so the
project still reaches 19 even if one module is not validated during evaluation (as recommended
in §IV). No "Module of choice" / custom module is claimed, so no extra justification is
required there.

---

## Individual Contributions

- **dirituay — Product Owner.** Defined the social-network product direction and feature
  priorities; validated delivered features. On the codebase, contributed to the auth and
  posts frontend.
- **rmarcas- — Project Manager / Scrum Master & Frontend lead.** Coordinated planning and
  delivery. Built most of the frontend: profile pages, the post feed and post UI, auth forms,
  the notification UI, and the PWA setup. _Challenge:_ keeping a fast-moving UI aligned with an
  evolving API — addressed by centralising HTTP in a shared axios instance.
- **lepereir — Tech Lead / Architect.** Restructured the backend to the
  golang-standards/project-layout, renamed the Go module, drove the integration test suite to
  ~81% coverage with a CI coverage gate, achieved a clean golangci-lint pass, wired and
  e2e-tested search, implemented GDPR export/delete, and authored the curl-based HTTP seeder.
  _Challenge:_ flaky tests from DB connection exhaustion under testcontainers — solved by
  sharing a single DB/Redis connection across the suite.
- **luluzuri — Developer.** Core backend: users, friends/follows, posts/likes/comments, file
  uploads, and large parts of auth; plus infrastructure and DB seeding. _Challenge:_ correct
  like/follow uniqueness and counters — solved with composite unique indexes and cached counts.
- **vali — Developer.** Real-time subsystem: the WebSocket chat hub, message persistence
  (DMs, rooms, threaded replies, attachments), and the search backend. _Challenge:_ delivering
  message history without loading entire conversations — solved with `since`-based pagination
  on room messages.

---

## Instructions

### Prerequisites

- **Podman** (preferred — the school runs Fedora) **or Docker**, with the `compose`
  subcommand. The Makefile auto-detects the engine (Podman first, Docker fallback); override
  with `make ENGINE=docker <target>`.
- **GNU Make** (wraps the common workflows).
- Optional for host-side dev: **Go 1.25+**, **Node 18+**, and **[mise](https://mise.jdx.dev/)**
  (auto-loads `infra/.env` so host-side `go run` sees the right env vars).

### Configuration (`.env`)

Environment variables live in `infra/.env`, which is git-ignored. A template is provided at
`infra/env.example` — copy it and adjust as needed:

```bash
cp infra/env.example infra/.env
```

Key variables (see the file for the full list): `API_PORT`, `FRONTEND_PORT`, `JWT_SECRET`,
`RATE_LIMIT_MAX`, `DB_USER` / `DB_PASSWORD` / `DB_NAME`, and the `GITHUB_CLIENT_ID` /
`GITHUB_CLIENT_SECRET` / `GITHUB_REDIRECT_URL` trio (fill these to enable GitHub login).
**Do not commit real secrets.**

### Run (single command)

```bash
make up
```

This builds the images and starts **backend + frontend (with nginx) + PostgreSQL + Redis** in
the background; run `make logs` to follow output. Other useful targets:

| Command | Description |
|---|---|
| `make logs` | Follow all logs (`logs-backend` / `logs-frontend` / `logs-db` for one) |
| `make re` | Clean (removes volumes) + rebuild + up |
| `make down` / `make clean` | Stop / stop and remove volumes |
| `make seed` | Seed the DB (Go seeder container) |
| `make shell-db` | `psql` into the running Postgres |
| `make help` | List all targets (and the detected engine) |

### Default ports

| Service | URL |
|---|---|
| App — frontend + nginx (proxies `/api` → backend) | http://localhost:3000 |
| Backend API (direct) | http://localhost:8000 |
| PostgreSQL | localhost:5432 |
| Redis | localhost:6380 |

### Tests

```bash
make test          # backend suite inside a golang:1.25 container
```

The suite uses **testcontainers** to spin up a real Postgres 15, so the Go container talks to
the host engine through its socket. With **Podman (rootless)**, enable the API socket once:

```bash
systemctl --user enable --now podman.socket
```

CI (GitHub Actions) runs the suite on every PR (open/sync) and on demand, and **fails the
build if backend coverage drops below 80%**; a golangci-lint job enforces formatting,
`go vet`, and static analysis.

### Seeding via the public API

`seed/seed.sh` populates a running instance through the HTTP API (curl + jq) — 10 users with
friendships, follows, posts, likes and comments:

```bash
./seed/seed.sh                                   # against http://localhost:8000
BASE_URL=http://localhost/api ./seed/seed.sh     # through nginx
```

---

## Resources

**Documentation & references**

- [Gin web framework](https://gin-gonic.com/docs/) · [GORM](https://gorm.io/docs/)
- [React](https://react.dev/) · [Vite](https://vitejs.dev/) · [Tailwind CSS](https://tailwindcss.com/docs)
- [vite-plugin-pwa](https://vite-pwa-org.netlify.app/) · [Web App Manifest (MDN)](https://developer.mozilla.org/docs/Web/Manifest)
- [JWT (RFC 7519)](https://datatracker.ietf.org/doc/html/rfc7519) · [JSON Web Tokens intro](https://jwt.io/introduction)
- [TOTP (RFC 6238)](https://datatracker.ietf.org/doc/html/rfc6238) · [`pquerna/otp`](https://github.com/pquerna/otp)
- [GitHub OAuth 2.0 apps](https://docs.github.com/apps/oauth-apps/building-oauth-apps)
- [gorilla/websocket](https://pkg.go.dev/github.com/gorilla/websocket) · [The WebSocket Protocol (RFC 6455)](https://datatracker.ietf.org/doc/html/rfc6455)
- [PostgreSQL 15](https://www.postgresql.org/docs/15/) · [Redis](https://redis.io/docs/)
- [testcontainers-go](https://golang.testcontainers.org/) · [golangci-lint](https://golangci-lint.run/)
- [golang-standards/project-layout](https://github.com/golang-standards/project-layout) · [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [GDPR — rights of the data subject](https://gdpr-info.eu/chapter-3/)

**How AI was used**

We used AI assistance for a few specific, well-scoped tasks; **all AI-generated output was
reviewed, adapted, and tested by the team member responsible**, and nothing was merged that
the author could not explain (per the subject's AI guidance, §I). Concretely, AI was used to:

- **Increase test coverage** — expanding the backend integration suite and filling gaps in
  test cases.
- **Review documents** — proofreading and sanity-checking written material such as the
  subject requirements and this README for completeness and consistency.
- **Research unfamiliar concepts and tools** — explaining technologies and patterns new to
  the team before we implemented them ourselves.

---

## Known Limitations & Before Evaluation

We are transparent about what is **not** finished, so the team can close these gaps before the
defense (several are mandatory per the subject):

- **⚠️ Privacy Policy & Terms of Service pages (mandatory).** The register form references them
  but the actual pages do not yet exist. The subject states that missing/inadequate legal pages
  **result in project rejection** — these must be added (accessible from the app, with real
  content) before evaluation.
- **⚠️ HTTPS (mandatory).** nginx currently listens on `:80` only. Browser↔backend traffic must
  use HTTPS; add a TLS listener (e.g. self-signed certs for the dev container) before the
  defense.
- **Chat UI.** The WebSocket chat backend is complete and message persistence works, but the
  frontend chat components (`features/chat/*`) and the `useWebSocket` hook are not yet wired —
  required to demonstrate the *real-time* and *user-interaction* majors end-to-end.
- **Search UI.** The search API is implemented and e2e-tested; the search bar/results UI is not
  yet wired.
- **Online status.** Presence indicators are not implemented yet (needed to fully satisfy the
  *standard user management* major).
- **Notifications coverage & confirmation emails.** Notifications cover the main social actions;
  GDPR confirmation emails are not sent yet.
- **Browser console.** A pass to ensure **no warnings/errors in the console** (subject
  requirement) is pending.

**Planned for the 125% target (⬜ modules):**

- **Multiple languages (i18n) — +1 pt.** No internationalization yet; add `react-i18next`
  with 3 complete languages and a UI language switcher.
- **Data import — +1 pt.** GDPR export exists; add an import path with validation and
  JSON/CSV/XML support to complete the data export/import module.

> Module statuses in this README reflect the current state of the repository. Items marked 🟧
> have a working backend and tests; finishing their UI (and the mandatory items above) is the
> remaining work before the project is evaluation-ready.
