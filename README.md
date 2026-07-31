# orbit

[![CI](https://github.com/Pittsburgh-Steelers-Fans/orbit/actions/workflows/ci.yml/badge.svg)](https://github.com/Pittsburgh-Steelers-Fans/orbit/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/go-1.24-00ADD8)
![License](https://img.shields.io/badge/license-MIT-green)

A REST API for team task and project management. Orbit provides authentication,
users, projects, tasks, comments, and notifications behind a small, well-factored
Go service, plus an `orbitctl` command line client.

## Quick start

```bash
cp .env.example .env
make docker-up          # starts Postgres + the orbit server
# or run the server directly:
make build && ./bin/orbit
```

## Environment variables

| Variable       | Default                          | Description                        |
| -------------- | -------------------------------- | ---------------------------------- |
| `PORT`         | `8080`                           | HTTP listen port                   |
| `DATABASE_URL` | `postgres://.../orbit`           | Postgres DSN                       |
| `JWT_SECRET`   | (required)                       | Secret used to sign access tokens  |
| `JWT_TTL`      | `15m`                            | Access token lifetime              |
| `LOG_LEVEL`    | `info`                           | zerolog level                      |
| `RATE_LIMIT_RPS`| `20`                            | Per-client request rate limit      |

## API reference

| Method | Path                         | Description                     |
| ------ | ---------------------------- | ------------------------------- |
| POST   | `/auth/register`             | Create an account               |
| POST   | `/auth/login`                | Exchange credentials for a JWT  |
| GET    | `/users`                     | List users (paginated)          |
| GET    | `/projects`                  | List projects for the caller    |
| POST   | `/projects`                  | Create a project                |
| POST   | `/tasks`                     | Create a task                   |
| PATCH  | `/tasks/{id}/status`         | Transition task status          |
| GET    | `/health`                    | Liveness probe                  |
| GET    | `/ready`                     | Readiness probe                 |

## Contributing

1. Branch from `main` using `feature/*`, `fix/*`, `chore/*`, or `refactor/*`.
2. Keep commits small and focused; every PR references an issue.
3. `make test lint` must pass before review.

## License

MIT
