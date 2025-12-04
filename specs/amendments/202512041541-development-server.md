# 202512041541 - Technical amendment specification

## Objective

During development the site should be served, and ready to be inspected. Take inspiration from: `https://gohugo.io/commands/hugo_server/#hugo-server`

Use Go `https://pkg.go.dev/net/http` to learn about setting up web servers.

The server should be very simple. Just be able to serve the static website incl. styles and images. 

## Technology

- Fetch and read: `https://pkg.go.dev/net/http`

### Skills

Use these Skills (in order):

1. `/skill:tdd-workflow` - Follow TDD process for all implementation
2. `/skill:go-scaffolder` - For creating new `internal/server/` package
3. `/skill:go-writer` - For implementing Go code (standard library only)
4. `/skill:go-tester` - For writing tests (table-driven, parallel)
5. `/skill:go-documenter` - For doc.go package documentation
6. `/skill:go-reviewer` - For code review after implementation
7. `/skill:go-analyzer` - For static analysis validation

## Examples

```bash
ssg serve                     # Serve output_dir on :8080
ssg serve -p 3000             # Serve on port 3000
ssg serve --build             # Build first, then serve
ssg serve --dir ./other-dir   # Serve specific directory
```

## Implementation Plan

### CLI Design

**Command**: `ssg serve`

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `ssg.yaml` | Path to config file |
| `--port` | `-p` | `8080` | Port to serve on |
| `--dir` | `-d` | (from config's output_dir) | Directory to serve |
| `--build` | `-b` | `false` | Build before serving |

### Files to Create

| File | Purpose |
|------|---------|
| `internal/server/doc.go` | Package documentation |
| `internal/server/server.go` | Server implementation |
| `internal/server/server_test.go` | Server tests |

### Files to Modify

| File | Changes |
|------|---------|
| `cmd/ssg/main.go` | Add `newServeCmd()`, `runServe()`, register in root command |
| `cmd/ssg/main_test.go` | Add tests for serve command |

### Test List (TDD)

#### Server Package (`internal/server/server_test.go`)

1. Server starts and serves files from directory
2. Server serves static assets (CSS, images) with correct content-type
3. Server serves nested directories (clean URLs: /about/ -> /about/index.html)
4. Server returns 404 for non-existent files
5. Server can be gracefully shutdown
6. Server reports address it's listening on

#### CLI Tests (`cmd/ssg/main_test.go`)

7. Serve command shows in help
8. Serve command has expected flags (--port, --config, --dir, --build)
9. Serve command validates port range (1-65535)
10. Serve command validates directory exists (when --build=false)

### Implementation Tasks

#### Phase 1: Server Package

**Task 1**: Create `internal/server/doc.go`
- `/skill:go-documenter`
- Package documentation following existing pattern

**Task 2**: Create Server struct and Config
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test Config struct exists with Port and Dir fields
- GREEN: Create struct

**Task 3**: Implement `New(Config) *Server`
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test Server can be created from Config
- GREEN: Implement constructor

**Task 4**: Implement `Server.Start() error`
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test server starts and accepts connections
- GREEN: Use `http.FileServer` and `http.ListenAndServe`

**Task 5**: Implement `Server.Addr() string`
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test returns listening address
- GREEN: Store and return address

**Task 6**: Implement `Server.Shutdown(ctx) error`
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test graceful shutdown stops server
- GREEN: Use `http.Server.Shutdown`

#### Phase 2: CLI Integration

**Task 7**: Add `newServeCmd()` to main.go
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test serve command appears in help
- GREEN: Create command with flags

**Task 8**: Implement `runServe()` function
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test serve loads config and starts server
- GREEN: Wire up config loading, optional build, server start

**Task 9**: Add graceful shutdown on Ctrl+C
- `/skill:tdd-workflow` + `/skill:go-writer`
- Setup signal handler calling Shutdown with 5-second timeout

#### Phase 3: Validation

**Task 10**: Add directory validation
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test error when directory doesn't exist (with --build=false)
- GREEN: Check directory before serving

**Task 11**: Add port validation
- `/skill:tdd-workflow` + `/skill:go-writer`
- RED: Test error for invalid port numbers
- GREEN: Validate port range 1-65535

#### Phase 4: Quality Assurance

**Task 12**: Code review
- `/skill:go-reviewer`
- Review all new code against Go style guide

**Task 13**: Static analysis
- `/skill:go-analyzer`
- Run `go vet`, `staticcheck`, `gopls codeaction`

### Key Implementation Details

#### Server Core (standard library only)

```go
handler := http.FileServer(http.Dir(s.dir))
s.server = &http.Server{Addr: addr, Handler: handler}
```

#### Graceful Shutdown

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
go func() {
    <-sigChan
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    server.Shutdown(ctx)
}()
```

#### Testing with Random Port

```go
// Use port 0 for tests to avoid port conflicts
srv := server.New(server.Config{Port: 0, Dir: dir})
```
