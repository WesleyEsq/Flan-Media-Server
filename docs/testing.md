# Testing Strategy & TDD Guidelines

This document details the architectural boundaries, dependency decoupling, and Test-Driven Development (TDD) practices for Flan Media Server.

---

## 1. Philosophy: Fast, Hermetic, and Decoupled

A major flaw in legacy web systems is tight coupling to external dependencies: tests require spinning up real databases, creating temporary files on disk, or binding to real network sockets. Over time, these suites become painfully slow, flaky, and hard to maintain.

Flan Media Server enforces three core testing rules:

1. **Zero Disk I/O in Unit Tests:** Logic that interacts with media files runs against in-memory virtual filesystems, never touching physical disk storage.
2. **Zero Database Coupling in Handlers:** HTTP handlers and business logic depend on narrow Go interfaces, not concrete SQLite connections (*sql.DB). Handlers are tested against lightweight in-memory mocks.
3. **Sub-Second Execution:** The entire unit test suite must execute in under 100 milliseconds via `go test ./...`.

---

## 2. Decoupling the Filesystem: In-Memory Virtual FS

Instead of writing tests that create real folders or copy dummy MP4 files on disk, the library scanner and media intake components depend on Go's standard `io/fs.FS` interface.

### The Production Interface

```go
// internal/scraper/scanner.go
func ScanDirectory(fileSystem fs.FS, mediaType string) ([]MediaCandidate, error) {
    // Walks fileSystem using fs.WalkDir
}
```

### The In-Memory Unit Test

In tests, use Go's standard `testing/fstest.MapFS` to construct a virtual directory in RAM:

```go
func TestScanDirectory_DetectsMediaAndLocalArtwork(t *testing.T) {
    mockFS := fstest.MapFS{
        "The Matrix (1999).mp4": &fstest.MapFile{Data: []byte("fake video data")},
        "poster.jpg":            &fstest.MapFile{Data: []byte("fake cover art")},
    }

    candidates, err := ScanDirectory(mockFS, "movies")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    if len(candidates) != 1 {
        t.Fatalf("expected 1 candidate, got %d", len(candidates))
    }

    item := candidates[0]
    if item.Title != "The Matrix" || item.ReleaseYear != 1999 {
        t.Errorf("unexpected metadata: %+v", item)
    }
}
```

#### Why This Matters

+ Runs in microseconds.
+ Cannot fail due to host file permissions or missing directories.
+ Requires no cleanup code (no os.RemoveAll, no temp directories).

---

## 3. Decoupling the Database: Consumer-Driven Interfaces

HTTP controllers (such as login, catalog viewing, and progress saving) do not accept `*sql.DB`. Instead, each package defines small, consumer-driven interfaces containing only the methods it actually uses.

### The Consumer Interface

```go
// internal/controller/auth.go (or internal/model/user.go)
type UserStore interface {
    GetByUsername(username string) (*User, error)
    VerifyPIN(userID int, pin string) (bool, error)
    RecordFailedAttempt(userID int) (lockedUntil *time.Time, err error)
    ResetFailedAttempts(userID int) error
}
```

### The Unit Test Mock

In unit tests, satisfy the interface with an inline mock struct:

```go
type mockUserStore struct {
    verifyPINFunc func(userID int, pin string) (bool, error)
}

func (m *mockUserStore) VerifyPIN(userID int, pin string) (bool, error) {
    return m.verifyPINFunc(userID, pin)
}

func (m *mockUserStore) GetByUsername(u string) (*User, error) { return nil, nil }
func (m *mockUserStore) RecordFailedAttempt(id int) (*time.Time, error) { return nil, nil }
func (m *mockUserStore) ResetFailedAttempts(id int) error { return nil }

func TestLoginController_RejectsInvalidPIN(t *testing.T) {
    store := &mockUserStore{
        verifyPINFunc: func(id int, pin string) (bool, error) {
            return false, nil // simulate wrong PIN
        },
    }

    controller := NewAuthController(store)
    req := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"pin":"0000"}`))
    rec := httptest.NewRecorder()

    controller.ServeHTTP(rec, req)

    if rec.Code != http.StatusUnauthorized {
        t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
    }
}
```

#### Why

+ Testing handler logic (cookie issuance, HTTP headers, error codes, rate limits) without needing a running database.
+ Tests execute in sub-millisecond time with zero database lock contention.

---

## 4. In-Memory HTTP Testing: `net/http/httptest`

Never spin up a live TCP server on port 4907 to test HTTP endpoints. Use Go's built-in `net/http/httptest` package.

```go
func TestStreamController_HandlesMissingFile(t *testing.T) {
    router := SetupRoutes(&mockStore{})

    req := httptest.NewRequest("GET", "/stream/movie/9999", nil)
    rec := httptest.NewRecorder()

    router.ServeHTTP(rec, req)

    if rec.Code != http.StatusNotFound {
        t.Errorf("expected 404 Not Found, got %d", rec.Code)
    }
}
```

+ Tests the full HTTP pipeline: routing, context injection, query parameters, and middleware.
+ Completely memory-bound and network-free.

---

## 5. Pure Function Testing (Zero-Dependency Units)

Core business logic is isolated into pure functions that take inputs and return outputs without any external dependencies:

1. **Filename Regex Sanitizer:**
   `CleanFilename(raw string) (title string, year int, season int, episode int)`
   Tested with table-driven tests against dozens of messy real-world filenames.
2. **HTTP Range Parser:**
   `ParseByteRange(rangeHeader string, fileSize int64) (start int64, length int64, err error)`
   Tested with boundary conditions (open-ended ranges, invalid offsets).
3. **EPUB Identifier & CFI Validator:**
   `ValidateCFI(cfi string) bool`
   Tested against valid and malformed Canonical Fragment Identifiers.

---

## 6. TDD Workflow & Conventions

When implementing features in Flan Media Server, adhere to the red-green-refactor cycle:

1. **Red:** Write a failing test in a `_test.go` file next to the code (e.g. `scanner_test.go`).
2. **Green:** Write the simplest implementation that satisfies the test.
3. **Refactor:** Clean up code, remove duplication, and optimize while keeping tests green.

### Running Tests

```bash
# Run all unit tests with execution timings
go test -v ./...

# Run tests with race condition detection
go test -race ./...
```

---

### Related Documentation

+ [Master System Specifications](docs/design.md)
+ [System Directories & Package Guide](docs/directories.md)
+ [Database Schema & Query Mocking](docs/database.md)
+ [Scraping Engine & Scanner Design](docs/scraper.md)
