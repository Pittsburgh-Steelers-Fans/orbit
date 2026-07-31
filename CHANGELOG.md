# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project adheres
to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Full-text search, task priorities, due dates, and project member management (in progress).

## [0.3.0] - 2026-06-18
### Added
- Cursor-based pagination across list endpoints.
- `/health` and `/ready` probes.
- Structured request logging with zerolog.

## [0.2.0] - 2026-05-20
### Added
- Task status state machine and transition endpoint.
- Comments on tasks and fan-out notifications.
### Fixed
- JWT expiry now validated on refresh.

## [0.1.0] - 2026-04-12
### Added
- Initial service scaffold: auth, users, projects, and tasks CRUD.
