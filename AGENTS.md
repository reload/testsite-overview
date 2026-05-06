# AGENTS.md

## Agent: GitHub Copilot

- The agent will keep working until all requested checks pass and the project is
  in a healthy state.
- If no test files exist, the agent will report this after running `go test`.
- The agent follows best practices for Go code style, documentation, and project
  maintenance.

- **Model:** GPT-4.1

- **Purpose:** Automated coding assistant for this project, supporting code quality,
  best practices, and development workflows.

- **Capabilities:**

  - Assists with code cleanup, refactoring, and documentation.
  - Ensures code quality through linting (`golangci-lint`), nil analysis
    (`go tool nilaway`), build/test validation (`go build`, `go test`),
    YAML linting (`yamllint`), and GitHub Actions workflow linting (`actionlint`).
  - Helps maintain consistent code style and project standards.
  - Provides guidance and automation for Go development tasks.
  - Responds to user requests for improvements, bug fixes, and code reviews.

## Typical Workflow

1. Receives user requests for code changes, improvements, or validation.
2. Analyzes relevant code files and project structure.
3. Applies changes or suggestions as needed.
4. Runs linting, nil analysis, build, and test tools to ensure project health.
5. Reports results and next steps to the user.

## Usage Examples

- Request: "Refactor a function for readability."
- Request: "Ensure all code passes lint and nil checks."
- Request: "Add documentation to exported functions."
- Request: "Verify the project builds and tests successfully."

## Linting Tools

- `golangci-lint` for Go code style and static analysis
- `go tool nilaway` for nil analysis
- `yamllint` for YAML file linting
- `actionlint` for GitHub Actions workflow linting
- `markdownlint-cli2` for linting Markdown files (run whenever you change or add
  Markdown files)

## Notes

- The agent will keep working until all requested checks pass and the project is
  in a healthy state.
- If no test files exist, the agent will report this after running `go test`.
- The agent follows best practices for Go code style, documentation, and project
  maintenance.
- The agent ensures all YAML files pass `yamllint` and all GitHub Actions workflows
  pass `actionlint`.
