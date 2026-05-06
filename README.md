# Test Site Overview

This project provides an overview dashboard for test environments (test sites)
managed via the Upsun platform. It fetches environment data using the Upsun API
and displays a list of active PR environments in a web interface.

## Features

- Lists active PR environments for a given Upsun project
- OAuth-based authentication to the Upsun API
- Simple, fast Go web server
- Uses [templ](https://templ.guide/) for HTML rendering

## Getting Started

### Prerequisites

- Go 1.21 or newer
- An Upsun API token and project ID

### Setup

1. Clone the repository:

   ```sh
   git clone https://github.com/reload/testsite-overview.git
   cd testsite-overview
   ```

2. Set the required environment variables:

   - `UPSUN_API_TOKEN`: Your Upsun API token
   - `UPSUN_PROJECT_ID`: The ID of your Upsun project
   - (Optional) `PORT`: Port for the web server (default: 80)
   - (Optional) `TITLE`: Custom title for the dashboard
   - (Optional) `LINK_REGEXP`: Custom regexp for filtering environment URLs

### Build and Run

```sh
go build
./testsite-overview
```

### Linting and Analysis

To ensure code quality, run:

```sh
golangci-lint run ./...
go tool nilaway ./...
```

### Testing

Currently, there are no automated tests. You can add Go test files to increase
coverage.

## Project Structure

- `main.go` — Main application logic and HTTP server
- `page_templ.go` / `page.templ` — Templ components for HTML rendering
- `Dockerfile` — Containerization support
- `AGENTS.md` — Agent automation and workflow documentation

## Development Notes

- Follows Go best practices for code style and documentation
- Uses `golangci-lint` and `nilaway` for static analysis
- See AGENTS.md for automation and agent workflow details

## License

See [LICENSE.md](LICENSE.md)
