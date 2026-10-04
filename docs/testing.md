# Testing Strategy & Guidelines

Flan Media Server uses a fast, hermetic testing suite that runs in-memory without spawning live network servers, writing temporary files to physical disks, or requiring an external database process.

---

## 1. Testing Principles

* **Zero Disk I/O:** Unit tests do not touch the physical filesystem. Scanners and file readers accept Go's standard `io/fs.FS` interface. Tests run against in-memory virtual filesystems using `testing/fstest.MapFS`.
* **Decoupled Handlers:** Controllers never depend directly on concrete database handles (`*sql.DB`). They interact with service interfaces that can be mocked in memory during testing.
* **Network-Free Handlers:** HTTP endpoints and middleware filters are tested using Go's standard `net/http/httptest` package without opening network ports.
* **Isolated Database Tests:** Repository integration tests run against ephemeral in-memory SQLite instances (`file::memory:?cache=shared`), executing full migrations and SQL queries without leaving artifacts on disk.

---

## 2. In-Memory Virtual Filesystem Testing

Production scanning logic consumes `io/fs.FS`, allowing tests to simulate arbitrary folder layouts with `testing/fstest.MapFS`:

```go
// internal/service/scanner_service_test.go
func TestScanDirectory_DetectsContainerAndFiles(t *testing.T) {
    mockFS := fstest.MapFS{
        "Breaking Bad/S01E01.mp4": &fstest.MapFile{Data: []byte("fake video data")},
        "Breaking Bad/S01E02.mp4": &fstest.MapFile{Data: []byte("fake video data")},
        "Breaking Bad/poster.jpg": &fstest.MapFile{Data: []byte("fake cover art")},
    }

    candidates, err := ScanDirectory(mockFS, "video")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(candidates) != 1 || candidates[0].Title != "Breaking Bad" {
        t.Errorf("unexpected candidates: %+v", candidates)
    }
    if len(candidates[0].Files) != 2 {
        t.Errorf("expected 2 episodes, got %d", len(candidates[0].Files))
    }
}
```

---

## 3. Controller Unit Testing with Mocks

Controllers accept service interfaces, allowing testing HTTP status codes, headers, and payload serialization with `net/http/httptest`:

```go
type mockAuthService struct {
    verifyPINFunc func(userID int64, pin string) (bool, error)
}

func (m *mockAuthService) VerifyPIN(userID int64, pin string) (bool, error) {
    return m.verifyPINFunc(userID, pin)
}

func TestLoginController_RejectsInvalidPIN(t *testing.T) {
    authService := &mockAuthService{
        verifyPINFunc: func(userID int64, pin string) (bool, error) { 
            return false, nil 
        },
    }
    controller := NewAuthController(authService)
    
    req := httptest.NewRequest("POST", "/api/login", strings.NewReader(`{"user_id":1,"pin":"0000"}`))
    rec := httptest.NewRecorder()

    controller.handleLogin(rec, req)

    if rec.Code != http.StatusUnauthorized {
        t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
    }
}
```

---

## 4. Pure Function & Cryptographic Testing

Core algorithms and security utilities are implemented as pure, zero-dependency functions:

* **Filename Sanitization:** `CleanFilename(raw string) (title string, year int, season int, episode int)` — tested with table-driven test cases against real-world media release formats.
* **Signed URL Verification:** `VerifySignedURL(secret []byte, fileID int64, userID int64, exp int64, sigHex string) bool` — tested against expired timestamps, modified user IDs, tampered signatures, and constant-time execution.
* **Reconciliation Diff Engine:** `DiffCatalog(disk []DiscoveredItem, db []StoredItem) ReconciliationPlan` — tested against new files, missing files, drive unmounts, and metadata lock preservation.

---

## 5. Executing Tests

```bash
# Run all unit tests
go test -v ./...

# Run tests with race detector enabled
go test -race ./...
```

---

## 6. Related Documentation

* [Master System Architecture](design.md)
* [Directory Structure & Architecture](directories.md)
* [Database Schema & Queries](database.md)
