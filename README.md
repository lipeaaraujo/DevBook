# DevBook

A small social network, created as a study project: a **Go** REST API backed by PostgreSQL, and a **React** client that consumes it. Users can post, follow each other and read a feed of the people they follow.

## Features

**API**

- User CRUD (create, list, get by ID, update, delete)
- JWT-based authentication
- Password hashing with bcrypt and email validation
- Posts with author details, and a feed built from who you follow
- Controller → Service → Repository architecture with dependency injection
- Custom API errors and consistent JSON responses
- Database migrations with `golang-migrate`
- Live reload during development via [Air](https://github.com/air-verse/air)
- API collection documented with Bruno

**Web**

- Log in and register
- Feed of posts from the people you follow
- Create, edit and delete your own posts
- Profiles with follower counts, and follow / unfollow
- Search for other users
- Edit your profile and change your password

## Tech Stack

**Web:** React 19, TypeScript, [Vite](https://vite.dev/), [Mantine](https://mantine.dev/) for UI, [TanStack Query](https://tanstack.com/query) for server state, React Router.

**API**

- **Language:** Go 1.25
- **Router:** [gorilla/mux](https://github.com/gorilla/mux)
- **Database:** PostgreSQL (with pgAdmin via Docker Compose)
- **Auth:** JWT ([dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go))
- **Other:** `lib/pq`, `golang.org/x/crypto`, `badoux/checkmail`, `joho/godotenv`

## Project Structure

```
.
├── api/                  # Go backend application
│   ├── main.go
│   ├── docker-compose.yml
│   ├── Makefile          # migration helpers
│   ├── migrations/       # SQL migrations
│   └── src/
│       ├── apierrors/    # custom API error types
│       ├── config/       # env loading and configuration
│       ├── db/           # database connection
│       ├── login/        # authentication (controller/service/routes)
│       ├── middlewares/  # logging and auth middlewares
│       ├── responses/    # JSON response helpers
│       ├── router/       # route registration
│       ├── users/        # user domain (controller/service/repository/model)
│       └── utils/        # hashing and token utilities
├── web/                  # React frontend
│   └── src/
│       ├── api.ts        # the only place that talks to the network
│       ├── hooks.ts      # TanStack Query hooks, one per endpoint
│       ├── routes.tsx    # route table and the auth guard
│       ├── pages/        # login, register, feed, profile, settings
│       └── components/   # shell, post card, post form, user search
└── bruno/                # Bruno API collection for documentation/testing
```

## Getting Started

### Prerequisites

- [Go](https://go.dev/) 1.25+
- [Docker](https://www.docker.com/) and Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI (for migrations)
- [Node.js](https://nodejs.org/) 20+ (for the web app)
- [Air](https://github.com/air-verse/air) (optional, for live reload)

### 1. Configure environment

```bash
cd api
cp .env.example .env
```

Edit `.env` and set the values, including a `JWT_TOKEN_SECRET`. You can also set `PORT` (defaults to `4000`).

### 2. Start the database

```bash
docker compose up -d
```

This starts PostgreSQL and pgAdmin (available at `http://localhost:${PGADMIN_PORT}`).

### 3. Run migrations

```bash
make migrate_up
```

To roll back:

```bash
make migrate_down
```

### 4. Run the API

```bash
go run main.go
```

Or with live reload:

```bash
air
```

The API will be available at `http://localhost:4000` (or your configured `PORT`).

### 5. Run the web app

```bash
cd web
npm install
npm run dev
```

Available at `http://localhost:5173`, with the API already running.

The API has no CORS middleware, so the browser never calls it directly — Vite proxies `/api` to `http://localhost:4000` instead, which makes every request same-origin. If you run the API on a different port, change the target in `web/vite.config.ts`. Deploying the two on separate origins needs a CORS middleware first (issue #2).

## Migrations

Migrations live in `api/migrations/` and are managed with the [golang-migrate](https://github.com/golang-migrate/migrate) CLI through the `Makefile`.

Create a new migration:

```bash
make create_migration NAME=add_some_table
```

This generates timestamped `*.up.sql` and `*.down.sql` files in `migrations/`. Edit them with your schema changes, then apply with `make migrate_up` (or roll back with `make migrate_down`).

## API Documentation (Bruno)

The `bruno/` directory contains a [Bruno](https://www.usebruno.com/) collection documenting all endpoints, with a `Local` environment preconfigured. Open the `bruno/DevBook` workspace in Bruno to explore and run the requests.

## Testing

```bash
cd api && go test ./...
cd web && npm test
```

## License

MIT License.
