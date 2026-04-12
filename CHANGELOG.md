# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Clean Architecture with domain/usecase/infrastructure layers
- GitHub Actions CI/CD pipelines (ci, release, docker workflows)
- GoReleaser configuration for cross-platform releases
- Pull request template
- Unit tests with mocks for usecase layer
- Edge case tests for strutil, file, flags, and usecase
- Dockerfile with tini for proper signal handling
- Graceful shutdown via `SIGINT`/`SIGTERM`
- Conventional Commits validation hook
- Documentation guides (`docs/CI-CD.md`, `docs/DEVELOPMENT.md`)

### Changed

- Migrated from `github.com/xanzy/go-gitlab` to `gitlab.com/gitlab-org/api/client-go`
- Upgraded to Go 1.26
- Renamed `tests/` to `test/` (Go standard)
- Renamed `internal/text/` to `internal/strutil/`
- CLI flag methods: `GetValue()` → `Value()`, `GetValueInt()` → `Int()`
- Renamed `GitlabFilePathsStruct` to `FileList`
- Simplified README with links to detailed docs

### Fixed

- Pluralization for Russian numbers 11–14, 111–114
- `sort.StringSlice` → `sort.Strings` (was not sorting)
- Ignored error from `GetProjects()` in app.go
- Potential nil dereference in `scanDir` response handling
- Closure capture bug in concurrent directory scanning
- CRLF line endings in `go.sum` breaking CI
- golangci-lint issues: errcheck, staticcheck

### Removed

- Legacy `internal/file/`, `internal/flags/`, `internal/gitlab/api/`, `internal/text/`
- Dead code in `internal/gitlab/api/api.go`
- Unused `SortFilePaths` function

[Unreleased]: https://github.com/xT10r/gitlab-file-scanner/compare/v0.1.0...HEAD
