# Dependency policy

Orbit keeps a small set of direct dependencies and updates them deliberately. Patch releases are safe to take during routine maintenance; minor and major releases should be reviewed for behavior changes before they are merged.

## Direct dependencies

- `github.com/go-chi/chi/v5`: HTTP router and middleware composition for the REST API.
- `github.com/golang-jwt/jwt/v5`: JWT parsing and signing for access tokens.
- `github.com/joho/godotenv`: Local development support for loading `.env` files.
- `github.com/rs/zerolog`: Structured logging with low allocation overhead.
- `github.com/stretchr/testify`: Assertions and helpers used by the test suite.

## Update policy

1. Prefer patch updates whenever CI is green and changelogs do not mention breaking changes.
2. Review minor updates in a dedicated pull request with focused testing around affected packages.
3. Avoid major updates unless the migration plan and compatibility impact are documented.
4. Do not vendor dependencies; rely on Go modules and repeatable CI builds.
