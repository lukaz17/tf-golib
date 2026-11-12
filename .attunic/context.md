# Project Context

Part of the project memory: human-readable memory of what this project is.

The reading and writing rules for this file and the rest of `.attunic/` are defined in [attunic-protocol.md](../.agents/attunic-protocol.md). The user fills this file in, manually or by asking an agent to help; change it only on such a request and after the user explicitly confirms that change (see [attunic-protocol.md's Write Rules](../.agents/attunic-protocol.md#write-rules)).

For filled-in examples of this file, see [tf-attunic-examples](https://github.com/tforce-io/tf-attunic-examples); use them as a reference for structure and level of detail, never as a source of this project's content.

## Abstract

A newcomer should be able to tell what the software does without knowing what it is built with. Write the abstract about the domain, not the technology, as one paragraph, in simple terms:
<!-- Project-specific / Abstract -->

TFgolib is a collection of reusable helpers that save Go developers from rewriting the same small utilities in every application. Each package covers one concern and can be imported on its own, so an application pulls in only what it uses.

## Technical Specifications

The technical aspects of the project: what it is built with and what working on it requires.

Write one bullet per aspect, as `- <Aspect>: <value>.`, such as primary language, frameworks, architecture style, build system, test runner, linter, databases, environment prerequisites.

Write `None` as the value of an aspect the project lacks but an agent would look for, such as a build system, test runner, or linter.

<!-- Project-specific / Technical Specifications -->

- Primary language: Go, minimum version 1.18.
- Project type: Go library, one module with many packages; no binaries.
- Architecture style: flat layout of independent packages, one per concern, with sub-packages for narrower helpers.
- Runtime dependencies: None beyond the standard library.
- Build system: go.
- Test runner: go test, with testify `assert`.
- Linter: gofmt, go vet.
- Databases: None.
- CI: None.

## Layout

[project-structure.md](../.agents/project-structure.md) holds the reference layouts of this project's languages and frameworks.

List of important folders, files and their roles in the project:

Write the list as a tree in a fenced block, one entry per line with its role as a `#` comment, mirroring the actual file system.

When the list is not empty and a top-level file/folder is missing from it, prompt the user to update it.

<!-- Project-specific / Layout -->

```
├── .agents/            # project AI agent instructions (attunic-protocol, coding-conventions, development-guideline, license-header, memory-template, project-structure, pull-request-guideline)
├── .attunic/           # project AI agent memory (context, memory, question)
├── diag/               # diagnostics of a running program
├── multiarch/          # helpers for code that runs on multiple platforms
├── multiplex/          # services running concurrently and exchanging messages
├── opx/                # shorthand expressions and control flow
│   └── slicext/        # slice helpers
├── random/             # random value generation
│   ├── pseudorng/      # fast pseudo-random generation (math/rand), not for security
│   └── securerng/      # cryptographically secure random generation (crypto/rand)
├── stdx/               # extensions of standard library types
│   ├── mathxt/         # math helpers for built-in number types
│   │   └── bigxt/      # math helpers for big numbers (math/big)
│   └── stringxt/       # string helpers
├── strfmt/             # parsing and formatting of well-known string formats
├── .copier-answers.yml # project Copier answers
├── .editorconfig       # EditorConfig configuration
├── .gitignore          # git ignore patterns
├── AGENTS.md           # project AI agent entry point
├── CLAUDE.md           # project Claude Code entry point
├── LICENSE             # MIT license text
├── README.md           # project overview
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Layer Dependencies

List of architecture layers and their relationships:

- Declare the innermost architecture layer first.
- An architecture layer may depend only on itself and on the architecture layers in its `Depends on` column, all declared before it.
- If an architecture layer appears in the table below, the dependency graph must be enforced.

| Architecture layer | Path | Depends on | Must never contain |
|--------------------|------|------------|--------------------|
| <!-- Project-specific / Layer Dependencies --> | | | |

### Deliberate Exceptions

List of documented exceptions:

- Record each exception with its reason and the condition for fixing it.
- Only exceptions to the architecture layers above belong here.
- An undocumented exception is a defect; a documented one is allowed, and revisited once its `Revisit-if` condition holds.

| What | Why | Revisit-if |
|------|-----|------------|
| <!-- Project-specific / Deliberate Exceptions --> | | |

## Supporting Tools

Use the tools declared here over alternatives, so every development environment gets the same results.

List each required tool/library under its own heading, followed by the commands for each use case.

Common use cases include: Build, Test, Lint, Format, Install.

<!-- Project-specific / Supporting Tools -->

Run every Build and Test command with both the local toolchain and Go 1.18.10, the last release of the module's minimum version; a change passes only when both pass, since the local toolchain alone accepts standard library APIs newer than Go 1.18. Switch with `GOTOOLCHAIN=go1.18.10 <command>` if the local toolchain is 1.21 or newer.

### Build

- `go build ./...` - compile every package; the module is a library, so nothing is written to disk.

### Test

- `go test ./...` - run the whole suite; narrow it with a package pattern, e.g. `go test ./stdx/...`.
- `go test -race ./...` - run the suite with the race detector (needs cgo and a C compiler); required when a test touches a goroutine or shared state, per [development-guideline.md](../.agents/development-guideline.md#go-testing).

### Lint

- `gofmt -l .` / `gofmt -w .` - check/apply formatting; the check must list no files.
- `go vet ./...` - static analysis.
- `go mod tidy -diff` / `go mod tidy` - check/fix that `go.mod` and `go.sum` match the imports; the check must report no diff. Local toolchain only (`-diff` needs Go 1.23+).

## Project-specific

<!-- Project-specific / Project-specific -->
