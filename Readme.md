# CorrelateApp Backend

## 🚀 Getting Started

To set up and run the CorrelateApp Backend locally, follow these steps:

1. Clone the repository

```bash
git clone https://github.com/dim2k2006/correlateapp-be.git
cd correlateapp-be
```

2. Install the dependencies

```bash
go mod download
```

3. Install golangci-lint

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

4. Install goimports

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

5. Run the Application:

```bash
go run ./cmd/api/main.go
```

## Managing dependencies

### To add a new dependency, run:

```bash
go get <package-name>
```

Note: Direct dependencies are those you explicitly import in your code. Indirect dependencies are required by your direct dependencies.

### Removing Dependencies

To remove a dependency, you can use the `go mod tidy` command. This command will remove any dependencies that are no longer required by your code.

```bash
go mod tidy
```

### Tidying Up

Regularly running go mod tidy ensures that your go.mod and go.sum files are clean and free of unnecessary dependencies.

```bash
go mod tidy
```

Benefits:

- Removes unused dependencies.
- Adds missing dependencies required by your imports.
- Ensures go.mod and go.sum are in sync with your codebase.

👍👍👍