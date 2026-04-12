# Contributing to gitlab-file-scanner

## Commit Message Format

This project uses **Conventional Commits** format. All commit messages must follow this structure:

```
type(scope): description

optional body (for complex changes)
```

### Types

| Type       | Description                                    |
|------------|------------------------------------------------|
| `feat`     | New feature or functionality                   |
| `fix`      | Bug fix                                        |
| `refactor` | Code refactoring (no behavior change)          |
| `docs`     | Documentation changes                          |
| `test`     | Adding or updating tests                       |
| `chore`    | Maintenance tasks, dependencies, config        |
| `ci`       | CI/CD, Docker, build pipeline changes          |
| `build`    | Build system, Makefile, dependencies           |

### Scopes

Use the package or module name as the scope:

| Scope      | Description                      |
|------------|----------------------------------|
| `app`      | `internal/app` - main logic      |
| `flags`    | `internal/flags` - CLI flags     |
| `api`      | `internal/gitlab/api` - GitLab   |
| `file`     | `internal/file` - file ops       |
| `text`     | `internal/text` - utilities      |
| `cmd`      | `cmd/` - entry point             |
| `docker`   | Dockerfile                       |
| `ci`       | Makefile, CI config              |
| `deps`     | Dependencies (go.mod)            |

### Description Rules

1. **English only** - all descriptions in English
2. **Imperative mood** - "add" not "added" or "adds"
3. **No period** at the end
4. **Lowercase** first letter (type handles casing)
5. **Be specific** - describe WHAT changed, not HOW

### Examples

```
✅ feat(api): add context support for GitLab requests
✅ fix(text): handle pluralization for numbers 11-14
✅ refactor(file): replace sort.StringSlice with sort.Strings
✅ test(file): add FilterFilesByMask test cases
✅ chore(deps): migrate to official gitlab.com/api/client-go
✅ ci(docker): add tini for proper signal handling
✅ build(deps): update go.mod to Go 1.26

❌ fixed the bug
❌ Добавил поддержку контекста
❌ Refactoring of file module.
❌ fix: update dependencies
```

### Breaking Changes

For breaking changes, add `!` after the type and document in the body:

```
feat(api)!: change GetProjects to return int64 IDs

Project ID type changed from int to int64 to match
the official GitLab API client.
```

## Branching Strategy

| Branch        | Purpose                           |
|---------------|-----------------------------------|
| `main`        | Stable releases                   |
| `develop`     | Active development                |
| `feature/*`   | New features (from develop)       |
| `fix/*`       | Bug fixes (from develop)          |
| `release/*`   | Release preparation (from develop)|
| `hotfix/*`    | Urgent fixes (from main)          |

## Pull Request Guidelines

1. One PR = one logical change
2. All tests must pass: `go test ./test/... -v`
3. Code must build: `go build ./cmd`
4. Update tests for new/changed functionality
5. Use draft PRs for work-in-progress

## Local Development

```bash
# Run tests
go test ./test/... -v

# Build binary
go build -o gitlab-file-scanner ./cmd

# Install commit hook
make install-hooks
```
