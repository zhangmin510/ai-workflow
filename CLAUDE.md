# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is an AI-native workflow project built with Go (>= 1.25.0). The project implements a three-pillar architecture:
1. **Executable Specifications (Specs)**: Define intent (WHAT/WHY) in `spec.md`, technical planning (HOW) in `plan.md`, and task breakdown in `tasks.md`
2. **Persistent Context**: Maintain consistent development context across sessions
3. **Orchestrable Actions**: Enable automated implementation through AI agents

## Development Commands
- **Run the application**: `go run hello.go`
- **Run all tests**: `go test -v ./...`
- **Run tests for a specific package**: `go test -v ./utils`
- **Code coverage**: `go test -cover ./...`
- **Linting**: `golangci-lint run` (required before commits)
- **Format code**: `gofmt -w .` (automatically applied after edits)
- **Add dependencies**: `go get <package_path>`
- **Build**: `go build -o bin/app .`

## Project Structure
- **Root**: Contains main example files (`hello.go`) and specification documents (`spec.md`, `plan.md`, `tasks.md`)
- **utils/**: Utility functions and their tests
- **.claude/**: Claude Code configuration and custom slash commands
- **.specify/**: Templates and scripts for the AI-native workflow system

## Key Dependencies
- `github.com/google/uuid`: For UUID generation (currently the only external dependency)
- Framework references in AGENTS.md suggest future use of:
  - Gin (`github.com/gin-gonic/gin`) for web framework
  - GORM (`gorm.io/gorm`) for database ORM

## Development Workflow
- **Branching**: Use `feature/<feature-name>` branches for new work
- **Commit messages**: Follow Conventional Commits format (e.g., `feat(api): add user registration`)
- **Testing**: All tests must pass before PR submission
- **Linting**: Run `golangci-lint run` before committing
- **PR titles**: Include Jira issue ID (e.g., `[PROJ-123] feat(api): Add User Registration`)

## AI Agent Integration
- Custom slash commands are available in `.claude/commands/` for workflow automation
- The project uses specification-driven development with `spec.md`, `plan.md`, and `tasks.md`
- Security reviews should use the `security-reviewer` sub-agent when needed