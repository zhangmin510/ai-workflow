# The Founder's Playbook: Building an AI-Native Startup

> *A practitioner's guide to restructuring software development so AI is a first-class participant — not a bolt-on tool.*

---

Your team just shipped the wrong feature. Not because the developers were bad. Not because the product manager was unclear. But because your AI coding assistant — despite being brilliant at generating code — had no idea that three months ago you decided to never use an ORM, that your authentication module is off-limits to automated refactoring, and that your current sprint is locked to three specific user stories.

The AI built exactly what you asked for in that session. It just didn't know what you'd already decided everywhere else.

This is the central problem of AI-assisted development in 2025: teams adopt AI tools without redesigning the system those tools operate within. The result is high-velocity chaos — fast code generation, slow coherence.

AI-native development solves this. It means restructuring how your team builds software so that AI agents can participate consistently, safely, and with institutional memory — not just session by session, but across weeks, team members, and product pivots.

This playbook covers exactly how to do it. It's written for technical founders, CTOs, and early-stage engineering leads who are serious about turning AI from a productivity supplement into a structural competitive advantage.

---

## The Three Frictions Killing Your AI Velocity

Before proposing solutions, it's worth naming the problem precisely. Most teams struggle with AI collaboration because of three distinct friction types. Each has a different root cause and a different cure.

### 1. Context Friction

Every new AI session starts from zero. Your assistant doesn't remember that you rejected microservices last quarter, that `internal/auth` is owned by a contractor who hates code changes during sprints, or that your Go version is 1.25 and you're using `log/slog` for structured logging — not `logrus`.

The symptom: you spend the first 10 minutes of every session re-briefing the AI. The root cause: you've treated AI like a contractor you hire by the hour, not a team member who accumulates context over time.

**The cure is Pillar Two: Persistent, Layered Context.**

### 2. Workflow Friction

Ad-hoc AI usage produces ad-hoc results. When one developer asks the AI to "write a user registration endpoint" and another asks it to "build the auth flow," you get two implementations with different error handling, different logging, different test styles, and different dependency choices — even if both used the same AI model on the same day.

The symptom: AI-generated code looks wildly different across your codebase. The root cause: you have no repeatable process for what the AI receives as input before it writes code.

**The cure is Pillar One: Executable Specifications.**

### 3. Cognitive Friction

When developers must supervise AI output manually — reading every generated line, correcting style violations, enforcing lint rules by hand — the cognitive overhead erodes the productivity gain. You end up with AI that writes fast and humans who review slowly.

The symptom: AI feels like a net positive on greenfield files but a liability on existing code. The root cause: you have no automated guardrails, so human attention compensates for missing system constraints.

**The cure is Pillar Three: Orchestrable Actions.**

---

## Pillar One: Executable Specifications

The most powerful change you can make is shifting from conversational AI usage to specification-driven AI usage.

In conversational mode, a developer types: *"Build me a booking endpoint that lets users reserve a time slot."* The AI does its best. The result depends entirely on what the AI infers about your system, your conventions, and your requirements — which is to say, it mostly guesses.

In specification-driven mode, the same request goes through a pipeline:

```
spec.md (WHAT/WHY) → plan.md (HOW) → tasks.md (breakdown) → AI agent (code)
```

Each stage transforms vague intent into machine-readable instructions. By the time an AI agent touches the codebase, it has a detailed, structured brief — not a freeform request.

### The spec.md: Define Intent Before Design

A `spec.md` file defines the WHAT and WHY of a feature in technology-agnostic terms. It's written for humans and AI alike. Here's the key structural element — user stories written as independently testable journeys:

```markdown
### User Story 1 - Core Booking Flow (Priority: P1)

**Why this priority**: Delivers the core user value; all other stories depend on it.

**Independent Test**: Can be fully tested by creating a booking end-to-end without any other story implemented.

**Acceptance Scenarios**:
1. **Given** an authenticated user, **When** they select an available date,
   **Then** a booking is created and a confirmation is returned.
2. **Given** a date is already fully booked, **When** the user selects it,
   **Then** they see a clear error message with alternative dates.
```

The "independently testable" constraint is critical. It forces you to decompose features into slices that can each be shipped, demonstrated, and validated on their own. AI agents work best on bounded, verifiable units of work — not sprawling feature bundles.

Below the user stories, the spec captures functional requirements (`FR-001`, `FR-002`...) and measurable success criteria (`SC-001`, `SC-002`...). This structure gives the AI a checklist it can verify against when generating code and tests.

### The plan.md: Design Before Building

Once the spec is locked, `plan.md` answers HOW the system will implement the spec. But here's the part most teams skip: before any design work begins, the plan must pass a **constitution check gate**.

```markdown
## Constitution Check (合宪性审查)

*GATE: Must pass before technical design begins.*

- [ ] **Simplicity Gate (Article I)**: Does this plan prioritize stdlib? Avoid unnecessary abstractions?
- [ ] **Test-First Gate (Article II)**: Does the plan include a "write tests first" step?
- [ ] **Explicitness Gate (Article III)**: Is dependency injection explicit? No globals?

*If any gate fails, document the justification in the Complexity Tracking section.*
```

This gate is not bureaucracy. It's the moment you verify that your architectural decisions align with your team's non-negotiable principles before an AI agent spends hours building something you'll have to tear down. We'll cover the constitution itself in depth in a later section.

### The tasks.md: Atomic, Parallelizable Work

The final stage of the pipeline is `tasks.md` — a breakdown of implementation into atomic tasks with explicit dependencies, phase ordering, and traceability back to user stories:

```markdown
## Phase 1: Core Data Layer [US1]

- [P] [US1] Define Booking struct with fields: ID, UserID, SlotID, CreatedAt
- [P] [US1] Write table-driven tests for BookingRepository.Create()
- [US1] Implement BookingRepository.Create() — tests must pass first
```

The `[P]` marker means "parallelizable" — an AI agent can run these simultaneously. The `[US1]` tag ties each task back to its user story, so the agent knows when a story is fully implemented (all its tasks are checked off).

When you run `/speckit.implement`, the agent reads `tasks.md`, executes tasks phase by phase, checks off each one as it completes, and moves to the next phase only when the current one is green. The spec-to-code pipeline becomes a single, auditable workflow.

---

## Pillar Two: Persistent, Layered Context

Context files are not documentation. They are infrastructure.

A README explains your project to a human who might contribute someday. A context file tells an AI agent, right now, exactly how to operate within your project — what tools to use, what patterns to follow, what to never touch.

The key insight from this repository's design is that context must be **layered**. You don't want every agent reading everything. You want each layer to provide the right context at the right level of specificity.

### The Four Memory Layers

**Layer 1: Enterprise (Company-Wide)**

Security red lines and compliance requirements that apply to every project in your organization. Examples: "Never log PII," "All secrets must use the secrets manager," "No direct database access from API handlers." These live at the organization level and every project inherits them.

**Layer 2: Project (Repository Root)**

This is `AGENTS.md` — the agent-portable specification of your project's development environment. It's not Claude-specific or Cursor-specific. Any AI agent that reads it gets the same information:

```markdown
# AGENTS.md for project "ai-workflow"

## 1. Development Environment
- **Language**: Go (>= 1.25.0)
- **Primary Framework**: Gin (github.com/gin-gonic/gin)
- **To run locally**: go run ./cmd/server/main.go
- **To add a dependency**: go get <package_path>

## 2. Testing Instructions
- **Run all tests**: go test -v ./...
- **Linting**: Before committing, run golangci-lint run

## 3. Git & PR Workflow
- **Commit Message Format**: Strictly follow Conventional Commits.
- **Branching**: All new features use branches named feature/<feature-name>.
```

`AGENTS.md` is the single source of truth that survives tool changes. When your team switches from Claude Code to another agent, the project knowledge transfers automatically.

**Layer 3: Tool-Specific Overrides (`.claude/CLAUDE.md`)**

On top of `AGENTS.md`, each AI tool gets a thin override file that adds tool-specific capabilities. In Claude Code, this looks like:

```markdown
# Import shared agent standards
@../AGENTS.md

# Claude-specific overrides

## Sub-agents
- For security review, invoke the `security-reviewer` sub-agent.

## Hooks
- Auto-run gofmt after every Go file edit.

# Personal preferences (git-ignored, locally loaded)
@~/.claude/personal-preferences.md
```

The `@` import syntax is powerful. The shared AGENTS.md is imported at the top, so the tool-specific file only contains what's unique to that tool. And the personal preferences import at the bottom loads each developer's individual style choices — from their home directory, so it's never committed to the repo and never overrides team conventions.

**Layer 4: Personal (`~/.claude/personal-preferences.md`)**

Individual preferences that should never be team policy: "I prefer verbose variable names," "Always suggest type aliases for complex function signatures," "Format code before explaining it." These live in the developer's home directory and overlay on top of everything else without polluting shared context.

### The Go + Claude Template as a Starting Point

This repository ships a `go-claude-template.md` that captures the AI collaboration directives for a Go project:

```markdown
## 4. AI Collaboration Directives

- [Principle] Stdlib First: Use stdlib when a reasonable solution exists.
- [Process] Review Before Build: Read related code first, propose a plan,
  wait for confirmation, then implement.
- [Practice] Table-Driven Tests: All tests MUST use table-driven style.
- [Practice] Concurrency Safety: Flag race conditions explicitly and document
  the safety mechanism used (mutex, channel, etc.).
```

These directives are not suggestions. They're instructions that shape every AI session in the project. An agent that reads this file will propose a review step before building, will write table-driven tests, and will flag any goroutine it introduces with an explicit safety note.

Start from a template like this. Fill in your stack. Add your non-negotiables. That's your project-level context layer.

---

## Pillar Three: Orchestrable Actions

Once your context is layered and your specifications are executable, the third pillar turns your AI setup into an automated system.

There are three categories of orchestrable actions: slash commands, hooks, and MCP servers. Together they form an extensible toolbox that grows with your needs.

### Slash Commands: Repeatable AI Workflows

A slash command is a Markdown file in `.claude/commands/` that defines a reusable AI workflow. It's not a script — it's a structured prompt with metadata:

```yaml
---
description: Create or update the feature specification from a natural language description.
handoffs:
  - label: Build Technical Plan
    agent: speckit.plan
    prompt: Create a plan for the spec I just wrote...
---

You are a product analyst. Read the user's feature description and generate
a spec.md following the template in .specify/templates/spec-template.md.
...
```

The `handoffs` field is what makes this powerful: when `/speckit.specify` finishes, it offers to hand off directly to `/speckit.plan`. The pipeline stages chain together, each feeding the next, each with its own context and constraints.

The full command suite in this repository covers the entire development lifecycle:
- `/speckit.specify` — transforms a feature idea into a structured spec
- `/speckit.clarify` — asks up to 3 targeted clarification questions before locking the spec
- `/speckit.plan` — generates a technical plan with constitution gate
- `/speckit.tasks` — breaks the plan into phase-ordered, parallelizable tasks
- `/speckit.implement` — executes tasks and marks completions
- `/commit` — generates a Conventional Commits message from the staged diff
- `/review-go-code` — runs a structured code review against the constitution

Each command is version-controlled. New team members get the full workflow the moment they clone the repo.

### Hooks: Automated Guardrails

Hooks run automatically in response to AI tool events. They enforce rules that should never depend on AI self-discipline.

This repository's hook configuration in `.claude/settings.json`:

```json
"hooks": {
  "PreToolUse": [{
    "matcher": "Edit|Write|MultiEdit",
    "hooks": [{
      "type": "command",
      "command": "\"$CLAUDE_PROJECT_DIR\"/.claude/hooks/check_main_branch.py"
    }]
  }],
  "PostToolUse": [{
    "matcher": "Edit|Write|MultiEdit",
    "hooks": [{
      "type": "command",
      "command": "FILE=$(jq -r '.tool_input.file_path'); case $FILE in *.go) gofmt -w $FILE ;; esac"
    }]
  }]
}
```

The `PreToolUse` hook runs `check_main_branch.py` before every file edit. If the agent is on the `main` branch, the hook exits with code 2, blocking the operation entirely and displaying an error. No AI session can directly edit `main` — the guardrail is structural, not advisory.

The `PostToolUse` hook auto-runs `gofmt` after every Go file edit. Code style is enforced the moment a file is written, not during a later review step.

Hooks belong in version control. They're team policy, not personal settings.

### MCP Servers: Extending the AI's Reach

Model Context Protocol (MCP) servers give your AI agent access to external systems. This project's `.mcp.json` wires in the GitHub MCP server:

```json
{
  "mcpServers": {
    "github": {
      "type": "http",
      "url": "https://api.githubcopilot.com/mcp/",
      "headers": {
        "Authorization": "Bearer ${GITHUB_TOKEN}"
      }
    }
  }
}
```

With this in place, an AI agent can read issues, comment on PRs, check CI status, and create branches — all without leaving the coding session. The `/speckit.taskstoissues` command, for example, converts `tasks.md` into actual GitHub issues with proper dependencies and labels, using the MCP connection.

The combinatorial power here is significant: slash commands define what to do, hooks enforce how it's done, and MCP servers expand where it can operate. The toolbox is nearly unlimited.

---

## Your Development Constitution

The constitution is the single most important artifact a founding team can produce. Every other pillar depends on it.

A constitution is a short document that defines your team's non-negotiable development principles. It has the highest priority in your context hierarchy — superseding any `CLAUDE.md` or session-level instruction. When an AI agent is about to do something that violates the constitution, it should stop and flag it, not proceed and apologize later.

Here's the full constitution from this repository, translated and adapted:

### Article I: Simplicity First

**Core**: Follow the "less is more" philosophy. Never introduce unnecessary abstraction. Never add a dependency that isn't strictly required.

- **YAGNI**: Implement only what the `spec.md` explicitly requires.
- **Stdlib First**: Unless there is a compelling reason, prefer the standard library. Use `net/http`, not Gin. Use `log/slog`, not `logrus`.
- **Anti-Over-Engineering**: Avoid complex design patterns. Simple functions and data structures outperform intricate interfaces and inheritance hierarchies.

### Article II: Test-First Imperative (Non-Negotiable)

**Core**: Every new feature or bug fix must begin with a failing test.

- **TDD Cycle**: Strictly follow Red-Green-Refactor. Write the failing test first. Make it pass. Then refactor.
- **Table-Driven Tests**: Unit tests must prefer table-driven style to cover multiple inputs and boundary conditions in a single, scannable structure.
- **No Mocks by Default**: Prefer integration tests with real dependencies or in-memory fakes over mock objects. Mocks test implementation details; real dependencies test behavior.

### Article III: Clarity and Explicitness

**Core**: Code's primary purpose is to be understood by humans. Machine execution is secondary.

- **Error Handling**: All errors must be explicitly handled. Never discard an error with `_`. Always wrap errors with context: `fmt.Errorf("booking: create: %w", err)`.
- **No Global Variables**: Never use globals to pass state. All dependencies must be explicitly injected through function parameters or struct fields.
- **WHY Comments**: Comments explain why the code does something unusual, not what it does. Well-named identifiers already do that.

### Article IV: Single Responsibility

**Core**: Each package, file, and function does exactly one thing well.

- **Package Cohesion**: The `github` package only talks to the GitHub API. It never contains Markdown conversion logic. Boundaries are strict.
- **Interface Segregation**: Define small, purpose-built interfaces. Avoid "god interfaces" that combine unrelated behaviors.

### Governance

The constitution has the highest authority in your project. Any plan generated by an AI agent must pass a constitutional review before design begins. The plan template enforces this with explicit gates that must be checked off before proceeding.

Version your constitution. When a principle changes, that's a major decision — treat it like one. Ratify it as a team, bump the version, and commit the change with a note explaining why the principle evolved.

---

## Bootstrapping Your AI-Native Stack

Here's a practical setup sequence for a new project:

**Step 1: Write your constitution**

Before writing a line of code, write `.specify/memory/constitution.md`. Gather your founding team for 90 minutes. Write down the things that would make you reject a PR — not stylistic preferences, but architectural non-negotiables. Ratify it, version it, commit it.

**Step 2: Create `AGENTS.md` at the repo root**

Document your development environment: language version, frameworks, test commands, build commands, git workflow. Keep it technology-specific but tool-agnostic. This file should make sense to any AI agent, not just Claude.

**Step 3: Create `.claude/CLAUDE.md` with the import chain**

Import `AGENTS.md`, add Claude-specific overrides (sub-agent definitions, hook configurations), and leave the personal preferences import at the bottom. Keep this file thin — it should extend `AGENTS.md`, not duplicate it.

**Step 4: Set up permissions and hooks**

In `.claude/settings.json`, explicitly allow safe commands (`go run`, `git add`, `go test`), deny destructive ones (`rm -rf`, `git push --force`), and configure your hooks. The main-branch protection hook and auto-formatter hook are the minimum viable set.

**Step 5: Install the spec pipeline**

Copy the `.specify/` directory into your repo. This gives you the spec, plan, tasks, and checklist templates. Register the `/speckit.*` slash commands from `.claude/commands/`. Your team now has a full spec-to-code workflow available as standard commands.

**Step 6: Connect MCP servers**

Create `.mcp.json` and wire in your critical external systems. Start with GitHub — it's the highest-leverage connection for development workflows. Add domain-specific APIs as your needs grow.

---

## Common Mistakes Founders Make

**Skipping the constitution.** Teams treat `CLAUDE.md` as a to-do list and skip writing the constitution. The result: AI makes different architectural decisions in every session. You get code that's internally consistent within a file but contradictory across the codebase.

**Treating spec.md as optional.** Jumping straight to "build it" means the AI builds the wrong thing correctly. User stories with acceptance scenarios force clarity before any planning begins. The 30 minutes you spend writing the spec saves hours of rework.

**One giant context file.** A 500-line `CLAUDE.md` that mixes enterprise security rules, project conventions, and personal style preferences becomes unmanageable. Use the four-layer model. Each layer contains only what belongs to it.

**No hooks, no guardrails.** Relying on the AI's self-discipline to not push to main or not skip formatting is not a system. Hooks enforce rules deterministically, regardless of what any session does or says.

**Tool lock-in through context files.** Writing instructions that only work in Claude Code means you're locked to Claude Code. Keep `AGENTS.md` as the canonical, tool-agnostic source. Tool-specific files (`CLAUDE.md`, `CURSOR.md`, `GEMINI.md`) extend it — they don't replace it.

**Over-specifying before validating.** Writing 50-point requirement lists before any user feedback is a waterfall antipattern with an AI wrapper. Keep specs lean. Use `/speckit.clarify` to ask at most three targeted questions, make informed guesses for the rest, and validate incrementally. The spec pipeline is designed for iteration, not perfection on the first pass.

---

## The Compounding Advantage

Here's the strategic case for doing this work upfront.

Every hour your team spends writing the constitution, layering context files, and building slash commands is an investment that compounds. Session 1 is slower than ad-hoc AI usage — there's scaffolding to set up. Session 50 is nearly frictionless. The AI already knows your architecture, your test style, your branching convention, and your non-negotiables. It doesn't ask. It doesn't guess. It executes.

Contrast this with traditional teams: when a senior developer leaves, they take their mental model with them. In an AI-native system, that mental model is encoded in `AGENTS.md`, `constitution.md`, and the specification pipeline. It stays with the project. It's readable by every new team member and every AI agent that joins the workflow.

The competitive moat is not which AI model you use. It's how well your team has systematized AI collaboration. Organizations that build this infrastructure compound faster than those that treat AI as a search box.

**Three things you can do this week:**

1. Write your first constitution — even just two principles. Ratify it with your co-founder or first engineer.
2. Create `AGENTS.md` at your repo root. Fill in the development environment section completely.
3. Add the main-branch protection hook to your `.claude/settings.json`. It takes five minutes and eliminates an entire class of accidents.

You don't need all three pillars in place to start. Pick one. Build the habit. The system compounds from there.

---

*This playbook is based on the open-source [ai-workflow](https://github.com/zhangmin510/ai-workflow) repository, which implements the three-pillar architecture described here. The templates, slash commands, hooks, and specification pipeline referenced throughout are available in that repository.*
