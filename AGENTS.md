# AGENTS.md for project "ai-workflow"

## 1. Development Environment
- **Language**: Go (>= 1.25.0)
- **Primary Framework**: Gin (`github.com/gin-gonic/gin`)
- **Database ORM**: GORM (`gorm.io/gorm`)
- **To run locally**: Use `go run ./cmd/server/main.go`
- **To add a dependency**: Use `go get <package_path>`

## 2. Testing Instructions
- **Run all tests**: Execute `make test` or `go test -v ./...`
- **Run tests for a specific package**: `go test -v ./internal/handlers`
- **Code Coverage**: Generate coverage report with `make coverage`
- **Linting**: Before committing, run `make lint`. This will execute `golangci-lint run`.

## 3. Git & Pull Request (PR) Workflow
- **Commit Message Format**: Strictly follow Conventional Commits. Example: `feat(api): add user registration endpoint`
- **Branching**: All new features must be developed in branches named `feature/<feature-name>`.
- **PR Title**: The PR title must reference the corresponding Jira issue ID. Example: `[PROJ-123] feat(api): Add User Registration`
- **Code Review**: Before requesting a review, ensure all tests pass and the linter is clean.
