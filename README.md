*This project has been created as part of the 42 curriculum by luluzuri, vali, lepereir, rmarcas-, and dirituay.*

# ft_transcendence — a real-time social network

> ft_transcendence (subject v21.1, the "Surprise." open-ended edition). A multi-user
> social web application: users register, build a profile, follow and add friends,
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
- **Social graph** — follow/unfollow and friend requests (send / accept / remove),
  with followers / following / friends lists.
- **Content** — create/read/update/delete posts with image attachments, likes, and reply.
- **Real-time chat** — WebSocket-based direct messages with message history.
- **Notifications** — in-app notifications for social actions (friend requests, followers, likes,
  messages, replies) and create and delete actions, also for gamification actions.
- **File management** — secure uploads (avatars, wallpapers, post media) with per-file
  visibility and access control.
- **Privacy** — GDPR data export and account deletion.
- **PWA** — installable, offline-capable progressive web app.

The application is a three-service stack (React frontend, Go API, PostgreSQL) plus Redis (Real-time caching and message broker) and
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
| **vali** | Developer | Real-time backend: WebSocket chat and message persistence. |

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
- **Task tracking** — GitHub Issues and GitHub Projects for sprint planning and task management.
- **Communication** — Whatsapp for day-to-day communication and coordination between frontend and backend tracks.

---

## Technical Stack

| Layer | Technology | Why |
|---|---|---|
| **Frontend** | React 18 + Vite 5 + Tailwind CSS 4.2.2 | React is treated as a framework by the subject (ecosystem + architecture); Vite gives a fast dev server and an easy PWA pipeline (`vite-plugin-pwa`); Tailwind is the chosen styling solution. Routing via `react-router-dom`, HTTP via `axios`. |
| **Backend** | Go 1.25 + Gin | Gin is a mature, high-performance HTTP framework. Go's static typing, goroutines, and standard tooling suit a concurrent real-time API (WebSocket hub + REST). |
| **ORM** | GORM (Postgres driver) | Clean model definitions, auto-migration, soft deletes, and relations — satisfies the ORM minor module. |
| **Database** | PostgreSQL 15 | Relational data with clear foreign-key relationships (users, posts, friendships, messages…), UUID primary keys, and strong consistency for concurrent multi-user writes. |
| **Cache / pub-sub** | Redis 7 | Token/session helpers, per-IP rate-limit buckets, and pub/sub for broadcasting real-time notifications and chat across connections. |
| **Auth** | `golang-jwt/jwt/v5`, `pquerna/otp` (TOTP), `golang.org/x/crypto` (bcrypt), GitHub OAuth 2.0 | Standard, audited libraries for JWT, RFC-6238 TOTP, and password hashing. |
| **Real-time** | `gorilla/websocket` | Battle-tested WebSocket implementation for the chat hub. |
| **Reverse proxy** | nginx (edge) | Single public entrypoint: terminates **TLS/HTTPS**, redirects HTTP→HTTPS, and routes `/api` (incl. WebSocket upgrade) + `/uploads` → backend and `/` → the static frontend. |
| **Containerisation** | Docker + Docker Compose | One-command, reproducible multi-service deployment. |
| **Testing / CI** | testcontainers-go, GitHub Actions, golangci-lint | Real-Postgres integration tests, a coverage gate, and static analysis. |

**Significant libraries:** `gin-contrib/cors`, `redis/go-redis/v9`, `gorm.io/driver/sqlite`
(test fallback), `vite-plugin-pwa`, `@heroicons/react` + `lucide-react` (icons).

---

## Database Schema

PostgreSQL with UUID primary keys. The entire schema — all 11 tables — is created and
kept in sync by GORM `AutoMigrate` from the structs in `backend/internal/models/`
(registered in `backend/internal/config/database.go`). Unless noted, strings map to
`text`, integers to `bigint`, and timestamps to `timestamptz`. Soft deletes
(`deleted_at`) are used on user-facing content.

### Tables

- **`users`** — accounts (local credentials or GitHub OAuth).
  `id uuid PK` (default `gen_random_uuid()`), `github_id varchar(255) unique nullable`,
  `provider varchar(50)` (default `'local'`), `display_name text`,
  `username text unique not null`, `email text unique not null`,
  `password text nullable` (bcrypt hash; null for OAuth-only accounts),
  `date_of_birth timestamptz nullable`, `two_fa_secret varchar(255) nullable`,
  `two_fa_enabled boolean` (default `false`), `avatar text`, `wallpaper text`, `bio text`,
  `created_at` / `updated_at`, `deleted_at` (soft delete).
- **`friends`** — directed friendship/follow edges between users.
  `id uuid PK`, `user_id uuid FK → users.id`, `friend_id uuid FK → users.id`,
  `status text` (`pending` | `accepted` | `following`);
  unique on `(user_id, friend_id, status)`.
- **`posts`** — user publications.
  `id uuid PK`, `author_id uuid FK → users.id not null`, `content text not null`,
  `media_url text nullable`, `media_mime varchar(100) nullable`,
  `tags text[]` (GIN-indexed), denormalized counters `likes_count` /
  `dislikes_count` / `comments_count bigint` (default `0`),
  `created_at` / `updated_at`, `deleted_at` (soft delete).
- **`likes`** — post reactions (likes **and** dislikes).
  `id uuid PK`, `user_id uuid FK → users.id`, `post_id uuid FK → posts.id`,
  `value bigint` (`1` like, `-1` dislike), `created_at`;
  unique on `(user_id, post_id)` — one reaction per user per post.
- **`replies`** — comments on posts.
  `id uuid PK`, `post_id uuid FK → posts.id` (indexed), `author_id uuid FK → users.id`,
  `content text not null`, `file_id uuid FK → files.id nullable` (attachment),
  `likes_count` / `dislikes_count bigint` (default `0`),
  `created_at` / `updated_at`, `deleted_at` (soft delete).
- **`reply_reactions`** — reactions on comments.
  `id uuid PK`, `user_id uuid FK → users.id`, `reply_id uuid FK → replies.id`,
  `value bigint` (`1` | `-1`), `created_at`; unique on `(user_id, reply_id)`.
- **`reposts`** — posts shared again by another user.
  `id uuid PK`, `post_id uuid FK → posts.id` (indexed), `author_id uuid FK → users.id`,
  `created_at`, `deleted_at` (soft delete).
- **`messages`** — chat messages.
  `id varchar(36) PK`, `sender_id varchar(36) not null` (indexed),
  `recipient_id varchar(36) nullable` (indexed; set for direct messages),
  `username text` (denormalized sender name), `type text`, `content text not null`,
  `file_id uuid nullable` (attachment), `room_id text nullable` (group/room chat),
  `parent_id varchar nullable` (self-reference → threaded replies), `created_at`.
- **`notifications`** — alerts delivered to users.
  `id varchar(36) PK`, `user_id text not null` (recipient),
  `actor_id text not null` (who triggered it), `user_username` / `actor_username text`
  (denormalized), `type text` (`friend_request` | `follow` | `like` | `comment` |
  `message`), `content text`, `post_id text` (related post, when applicable),
  `read boolean` (default `false`), `created_at`.
- **`files`** — uploaded files (avatars, post media, attachments).
  `id uuid PK`, `owner_id uuid not null` (indexed), `path text not null`
  (server-side storage path, never exposed), `filename varchar(255) not null`,
  `mime_type varchar(100) not null`, `size bigint not null` (bytes),
  `visibility varchar(20)` (`public` | `friends` | `private`, default `'public'`),
  `created_at`.
- **`file_access`** — per-file ACL join table.
  Composite PK `(file_id uuid, user_id uuid)`; grants a user access to a
  non-public file.

### Relationships

- A **user** has many posts, reactions (`likes`, `reply_reactions`), comments
  (`replies`), reposts, sent/received **messages**, **notifications** (as recipient
  and as actor), and **files**.
- **friends** stores both friendship and follow edges as directed `(user_id → friend_id)`
  rows distinguished by `status`.
- A **post** has many `likes`, `replies`, and `reposts`; a **reply** has many
  `reply_reactions`. Reaction tables store likes and dislikes in `value` (`1`/`-1`),
  one per user per target.
- **messages** support both direct messages (`recipient_id`) and rooms (`room_id`), plus
  threaded replies (`parent_id` self-reference) and optional attachments (`file_id`).
- **replies** and **messages** may reference a **file** as attachment (`file_id`).
- **file_access** is a join table implementing per-file access control for
  non-public files.

---

## Modules

The mandatory minimum is **14 points** (Major = 2, Minor = 1). The bonus part is capped at
**5 points** and counts only once all 14 mandatory points are fully functional — so the
maximum gradeable total is **14 + 5 = 19 points (125%)**, which is our target. We claim **20
points** (a 1-point buffer, since a module that fails to validate at the defense scores 0).
Per the subject, we list our full intended module set with an **honest implementation
status**, since only fully functional modules count at evaluation:

- ✅ **Done** — works end-to-end (UI + API).

### Major modules (2 pts each)

| Module | Category | Pts | Status | How it was implemented | Owner |
|---|---|---|---|---|---|
| Use a framework for frontend **and** backend | Web | 2 | ✅ | Chosen to leverage mature ecosystems and separation of concerns. React (Vite) for a reactive, component-based UI; Gin (Go) for a high-performance, concurrent REST API.   | whole team |
| Real-time features via WebSockets | Web | 2 | ✅ | Chosen to enable instant notifications and live chat without polling. gorilla/websocket hub with broadcast, connect/disconnect handling, Redis pub/sub; real-time chat UI fully wired with message history and online status. | vali, lepereir |
| User interaction (chat + profile + friends) | Web | 2 | ✅ | Core social feature required by the subject. Profiles, friends, follows, and real-time chat all implemented end-to-end. Chat UI with message history, online status, and user ordering by last message. | vali, dirituay, rmarcas- |
| Public API | Web | 2 | ✅ | Chosen to expose the platform data securely to external consumers. 30+ REST endpoints (GET/POST/PUT/DELETE), per-IP rate limiting, Swagger docs. | luluzuri, lepereir, vali |
| Standard user management | User Mgmt | 2 | ✅ | Required by the subject for a complete social platform. Profile edit, avatar with default, friends system, online status indicator (green dot), and level badge on profile. | whole team |

**Major subtotal: 10 pts**

### Minor modules (1 pt each)

| Module | Category | Pts | Status | How it was implemented | Owner |
|---|---|---|---|---|---|
| ORM | Web | 1 | ✅ | Chosen to simplify database interactions. GORM with auto-migration, relations, and soft deletes. | luluzuri |
| Notification system | Web | 1 | ✅ | Chosen to improve user engagement with real-time feedback. Notifications for all social actions via WebSocket; mark-all-read and per-notification read status. | dirituay, vali |
| File upload & management | Web | 1 | ✅ | Chosen to support rich media content. Multi-type uploads, client+server validation, visibility/ACL, delete, serving. | luluzuri, rmarcas- |
| Progressive Web App (PWA) | Web | 1 | ✅ | Chosen to make the app installable and available offline. vite-plugin-pwa, manifest, Workbox service worker. | rmarcas- |
| OAuth 2.0 (GitHub) | User Mgmt | 1 | ✅ | Chosen to simplify onboarding with an existing account. OAuth login + callback, account linking. |  luluzuri |
| Two-Factor Authentication (TOTP) | User Mgmt | 1 | ✅ | RFC-6238 TOTP setup/enable/disable/verify with authenticator apps | luluzuri |
| GDPR compliance | Data & Analytics | 1 | ✅ | Required by European data protection law. Data export in JSON and account deletion with confirmation overlay. | lepereir, dirituay |
| Custom-made design system | UX UI  | 1 | ✅ | Chosen to ensure visual consistency across the entire app. Color palette, typography, and 10+ reusable components: Button, Input, Modal, PostCard, CreatePost, NotificationBell, FollowButton, CookiesBanner, LoginModal, ProfileBadges. | dirituay |
| Support for additional browsers | Accessibility  | 1 | ✅ | Chosen to maximize accessibility. Full compatibility with Firefox, Brave, and Chrome verified and tested. | whole team |
| Search | Web | 1 | ✅ |  Implement advanced search functionality with filters, sorting, and pagination. |  lepereir, dirituay, vali |
| Gamification system  | Gaming | 1 | ✅ |  A gamification system to reward users for their actions. |  lepereir, dirituay, rmarcas- |

**Minor subtotal: 11 pts**

### Point calculation

```
Major modules :  5 × 2 = 10 pts
Minor modules :  11 × 1 =  11 pts
                       -------
Total claimed          = 21 pts   (mandatory: 14 · bonus cap: 19 = 125%)
```

We target the **125% ceiling (19 points = 14 mandatory + 5 bonus)** and claim exactly **19**,
which leaves no buffer if a module fails to validate during evaluation. No "Module of choice" /
custom module is claimed, so no extra justification is required there.

---

## Individual Contributions

- **dirituay — Product Owner.** Defined the social-network product direction and feature
  priorities; validated delivered features. On the codebase, implemented the notification system UI, the design system, profile improvements, the setting page, the responsive mobile layout improvements, and the right sidebar with user suggestions. _Challenge:_ integrating real-time badge notifications without interfering with the social notification system.
- **rmarcas- — Project Manager / Scrum Master & Frontend lead.** Coordinated planning and
  delivery, organised sprints and managed blockers. On the codebase, built the profile pages, the post feed and post UI, auth forms, the gamification system and the PWA setup. _Challenge:_ keeping a fast-moving UI aligned with an
  evolving API — addressed by centralising HTTP in a shared axios instance.
- **lepereir — Tech Lead / Architect.** Restructured the backend to the
  golang-standards/project-layout, renamed the Go module, drove the integration test suite to
  ~81% coverage with a CI coverage gate, achieved a clean golangci-lint pass,
  implemented GDPR export/delete, and authored the curl-based HTTP seeder.
  _Challenge:_ flaky tests from DB connection exhaustion under testcontainers — solved by
  sharing a single DB/Redis connection across the suite, needed to pay to send emails.
- **luluzuri — Developer.** Core backend: friends/follows, posts/likes/comments, file
  uploads, and large parts of auth; plus infrastructure _Challenge:_ correct
  like/follow uniqueness and counters, manage the upload of differents files — solved with composite unique indexes and cached counts.
- **vali — Developer.** Real-time subsystem: the WebSocket and message persistence
  (DMs, search of message, search of users), added user services and enpoints, notification system . _Challenge:_ delivering
  message history — solved with redis.

---

## Instructions

### Prerequisites

- **Podman** (preferred — the school runs Fedora) **or Docker**, with the `compose`
  subcommand. The Makefile auto-detects the engine (Podman first, Docker fallback); override
  with `make ENGINE=docker <target>`.
- **GNU Make** (wraps the common workflows).
- Optional for host-side dev: **Go 1.25+**, **Node 18+**.

### Configuration (`.env`)

Environment variables live in `infra/.env`, which is git-ignored. A template is provided at
`infra/env.example` — copy it and adjust as needed:

```bash
cp infra/env.example infra/.env
```

The variables are: `FT_TRANSCENDENCE_URL` (the app's public home URL), `JWT_SECRET`,
`DB_USER` / `DB_PASSWORD` / `DB_NAME`, and the `GITHUB_CLIENT_ID` / `GITHUB_CLIENT_SECRET` /
`GITHUB_REDIRECT_URL` trio (fill these to enable GitHub login). Ports are **fixed** in
`infra/docker-compose.yml`, not configured here. **Do not commit real secrets.**

### Run (single command)

```bash
make up
```

This builds the images and starts the **nginx reverse proxy + frontend + backend + PostgreSQL**
+ **Redis** in the background; run `make logs` to follow output. Other useful targets:

| Command | Description |
|---|---|
| `make logs` | Follow all logs (`logs-backend` / `logs-frontend` / `logs-db` for one) |
| `make re` | Clean (removes volumes) + rebuild + up |
| `make down` / `make clean` | Stop / stop and remove volumes |
| `make seed` | Seed the DB (Go seeder container) |
| `make shell-postgres` / `shell-redis` | `psql` / `redis-cli` into the database / cache |
| `make help` | List all targets (and the detected engine) |

### Default ports

| Service | URL |
|---|---|
| **Everything** — frontend, API (`/api`), uploads — via the nginx proxy (HTTPS) | **https://localhost:3000** |
| backend / PostgreSQL / Redis | **not published** — internal only (use `make shell-{backend,postgres,redis}`) |

The proxy is the **only published port** — all access goes through nginx. It's fixed at **3000**
(set in `infra/docker-compose.yml`) because **rootless Podman can't bind privileged ports < 1024**,
and serves **HTTPS** there with a **self-signed certificate**. The cert is generated once on the
host into `infra/certs/` by `make certs` (run automatically by `make up`) and bind-mounted into
the proxy, so it stays stable across rebuilds — the one-time browser "not secure" warning
(Advanced → proceed) you accept therefore keeps working. A plain-HTTP request to the port is
auto-redirected to HTTPS. The
backend listens on a fixed internal port **8080** (the proxy's upstream; not exposed to the host).

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

- **Browser console.** The app uses a self-signed certificate generated by nginx for local development. Browsers may show a security warning on first access — accept it once via Advanced → Proceed and it will not reappear.
