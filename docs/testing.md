# Testing Strategy & TDD Guidelines

Flan Media Server enforces a fast, hermetic testing suite that runs in sub-100ms times without spinning up live TCP servers, creating real files on disk, or requiring a running database.

---

## 1. Core Testing Philosophy

| Principle | Rule | Implementation |
| :--- | :--- | :--- |
| **Zero Disk I/O** | Unit tests must never touch the physical filesystem. | Scanners and parsers accept Go's standard `io/fs.FS` interface. Tests run against in-memory virtual filesystems via `testing/fstest.MapFS`. |
| **Zero DB Coupling in Handlers** | HTTP controllers must never depend directly on `*sql.DB`. | Controllers accept narrow consumer-driven interfaces (e.g. `UserStore`, `MovieStore`). Unit tests mock these in memory. |
| **Network-Free Handlers** | Never spin up live TCP servers on port 4907 for testing. | Endpoints and middlewares are verified using standard `net/http/httptest`. |
| **Sub-Second Execution** | The entire suite must execute in milliseconds. | Allows continuous test-driven development (`TDD`) without friction. |

---

## 2. Decoupled Filesystem: In-Memory Virtual FS

Production code consumes `io/fs.FS`, while tests use `testing/fstest.MapFS`:

```go
// internal/scraper/scanner_test.go
func TestScanDirectory_DetectsMediaAndLocalArtwork(t *testing.T) {
    mockFS := fstest.MapFS{
        "The Matrix (1999).mp4": &fstest.MapFile{Data: []byte("fake video data")},
        "poster.jpg":            &fstest.MapFile{Data: []byte("fake cover art")},
    }

    candidates, err := ScanDirectory(mockFS, "movies")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(candidates) != 1 || candidates[0].Title != "The Matrix" {
        t.Errorf("unexpected candidates: %+v", candidates)
    }
}
```

---

## 3. Decoupled Handlers: Consumer-Driven Interfaces

Controllers accept small interfaces containing only the methods they invoke:

```go
// Controller Interface
type UserStore interface {
    VerifyPIN(userID int, pin string) (bool, error)
}

// Unit Test Mock
type mockUserStore struct {
    verifyPINFunc func(userID int, pin string) (bool, error)
}

func (m *mockUserStore) VerifyPIN(id int, pin string) (bool, error) {
    return m.verifyPINFunc(id, pin)
}

func TestLoginController_RejectsInvalidPIN(t *testing.T) {
    store := &mockUserStore{
        verifyPINFunc: func(id int, pin string) (bool, error) { return false, nil },
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

---

## 4. Pure Function Testing

Core utilities are implemented as pure, zero-dependency functions:

* **Filename Sanitizer:** `CleanFilename(raw string) (title string, year int, season int, episode int)` — tested with table-driven tests against messy real-world filenames.
* **HTTP Range Parser:** `ParseByteRange(header string, size int64) (start int64, length int64, err error)` — tested against boundary and overflow conditions.
* **EPUB Identifier & CFI Validator:** `ValidateCFI(cfi string) bool` — tested against valid and malformed CFIs.

---

## 5. Running Tests

```bash
# Run all unit tests
go test -v ./...

# Run tests with race condition detection
go test -race ./...
```

---

## 6. Related Documentation

* [Master System Architecture](design.md)
* [Directory Structure](directories.md)
* [Database Schema](database.md)
