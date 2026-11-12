# Project Context

Part of the project memory: human-readable memory of what this project is.

The reading and writing rules for this file and the rest of `.attunic/` are defined in [attunic-protocol.md](../.agents/attunic-protocol.md). The user fills this file in, manually or by asking an agent to help; change it only on such a request and after the user explicitly confirms that change (see [attunic-protocol.md's Write Rules](../.agents/attunic-protocol.md#write-rules)).

For filled-in examples of this file, see [tf-attunic-examples](https://github.com/tforce-io/tf-attunic-examples); use them as a reference for structure and level of detail, never as a source of this project's content.

## Abstract

A newcomer should be able to tell what the software does without knowing what it is built with. Write the abstract about the domain, not the technology, as one paragraph, in simple terms:
<!-- Project-specific / Abstract -->

## Technical Specifications

The technical aspects of the project: what it is built with and what working on it requires.

Write one bullet per aspect, as `- <Aspect>: <value>.`, such as primary language, frameworks, architecture style, build system, test runner, linter, databases, environment prerequisites.

Write `None` as the value of an aspect the project lacks but an agent would look for, such as a build system, test runner, or linter.

<!-- Project-specific / Technical Specifications -->

## Layout

[project-structure.md](../.agents/project-structure.md) holds the reference layouts of this project's languages and frameworks.

List of important folders, files and their roles in the project:

Write the list as a tree in a fenced block, one entry per line with its role as a `#` comment, mirroring the actual file system.

When the list is not empty and a top-level file/folder is missing from it, prompt the user to update it.

<!-- Project-specific / Layout -->

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

## Project-specific

<!-- Project-specific / Project-specific -->
