# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Fixed

## [v1.0.0] - 2026-04-17

Initial release of GitLab File Scanner, a CLI tool for scanning GitLab repositories and exporting file lists.

### Added

- Clean Architecture with domain/usecase/infrastructure layers
- GitHub Actions CI/CD pipelines (ci, release, docker workflows)
- GoReleaser configuration for cross-platform releases
- Pull request template
- Unit tests with thread-safe mocks for usecase layer
- Edge case tests for strutil, filefilter, filesystem, consolelog, and usecase
- Dockerfile with tini for proper signal handling
- Graceful shutdown via `SIGINT`/`SIGTERM`
- Conventional Commits validation
- Documentation guides (`docs/CI-CD.md`, `docs/DEVELOPMENT.md`)
- Logger level constants (`LevelDebug`, `LevelInfo`, `LevelWarn`, `LevelError`)
- `Debug()` method to logger interface
- License header check in CI (`HEADER` file, `make license-check`)

### Changed

- Migrated from `github.com/xanzy/go-gitlab` to `gitlab.com/gitlab-org/api/client-go`
- Upgraded to Go 1.26 (1.26.2 required for CVE fixes)
- Renamed `tests/` to `test/` (Go standard)
- Renamed `internal/text/` to `internal/strutil/`
- CLI flag methods: `GetValue()` → `Value()`, `GetValueInt()` → `Int()`
- Renamed `GitlabFilePathsStruct` to `FileList`
- Simplified README with links to detailed docs

### Fixed

- Goroutine leak in `scanRepository` (semaphore acquired before defer)
- Duplicate files in filter results when masks overlap
- Pluralization for Russian numbers 11–14, 111–114
- Format string panic when log message contains `%` without args
- Regex injection via unescaped `[`, `]`, `\`, `+`, `?`, `{`, `}`, `^`, `$`
- Race conditions in mock implementations (added `sync.RWMutex`)
- `sort.StringSlice` → `sort.Strings` (was not sorting)
- Ignored error from `GetProjects()` in app.go
- Potential nil dereference in `scanDir` response handling
- Closure capture bug in concurrent directory scanning
- CRLF line endings in `go.sum` breaking CI
- Duration output: `[1d 2h]` → `1d 2h` (`fmt.Sprintf` → `strings.Join`)

[v1.0.0]: https://github.com/xT10r/gitlab-file-scanner/releases/tag/v1.0.0
