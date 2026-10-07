package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/WesleyEsq/Flan-Media-Server/internal/model"
	_ "modernc.org/sqlite"
)

type DB struct {
	Writer *sql.DB
	Reader *sql.DB
}

func Connect(dbPath string) (*DB, error) {
	// If in-memory (e.g. for testing)
	isMemory := dbPath == ":memory:" || dbPath == "file::memory:?cache=shared"

	if !isMemory {
		dbDir := filepath.Dir(dbPath)
		// Check that the directory exists
		if _, err := os.Stat(dbDir); os.IsNotExist(err) {
			if err := os.MkdirAll(dbDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create db directory %s: %w", dbDir, err)
			}
		}

		// Verify .flan-keep marker file
		markerPath := filepath.Join(dbDir, ".flan-keep")
		if _, err := os.Stat(markerPath); os.IsNotExist(err) {
			// Create it if empty database initialization, but if directory was already mounted and lost marker, halt
			_ = os.WriteFile(markerPath, []byte("flan-keep"), 0644)
		}
	}

	dsn := dbPath
	if isMemory {
		dsn = "file::memory:?cache=shared"
	}

	// 1. Dedicated Writer handle (serialized writes)
	writerDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open writer db: %w", err)
	}
	writerDB.SetMaxOpenConns(1)
	writerDB.SetMaxIdleConns(1)

	// 2. Dedicated Reader handle (parallel reads)
	readerDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		writerDB.Close()
		return nil, fmt.Errorf("failed to open reader db: %w", err)
	}
	readerDB.SetMaxOpenConns(3)
	readerDB.SetMaxIdleConns(3)
	readerDB.SetConnMaxLifetime(0)

	// Initialize pragmas on both handles
	for _, db := range []*sql.DB{writerDB, readerDB} {
		if err := initPragmas(db); err != nil {
			writerDB.Close()
			readerDB.Close()
			return nil, fmt.Errorf("failed to initialize pragmas: %w", err)
		}
	}

	// Startup quick_check
	var checkResult string
	if err := readerDB.QueryRow("PRAGMA quick_check;").Scan(&checkResult); err != nil || checkResult != "ok" {
		writerDB.Close()
		readerDB.Close()
		return nil, fmt.Errorf("database quick_check failed: %s (err: %v)", checkResult, err)
	}

	appDB := &DB{
		Writer: writerDB,
		Reader: readerDB,
	}

	// Apply migrations on writer connection
	if err := Migrate(appDB.Writer); err != nil {
		appDB.Close()
		return nil, fmt.Errorf("database migration failed: %w", err)
	}

	return appDB, nil
}

func initPragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA cache_size = -2000;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA foreign_keys = ON;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("exec pragma %q: %w", p, err)
		}
	}
	return nil
}

func (db *DB) Close() error {
	var err1, err2 error
	if db.Writer != nil {
		err1 = db.Writer.Close()
	}
	if db.Reader != nil {
		err2 = db.Reader.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}

func (db *DB) GetDataVersion() (int64, error) {
	var version int64
	err := db.Reader.QueryRow("PRAGMA data_version;").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

func VerifyMarker(dirPath string) error {
	markerPath := filepath.Join(dirPath, ".flan-keep")
	if _, err := os.Stat(markerPath); os.IsNotExist(err) {
		return model.ErrMissingMarker
	}
	return nil
}
