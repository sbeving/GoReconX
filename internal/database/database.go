package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Initialize initializes the database connection and creates tables
func Initialize() (*sql.DB, error) {
	// Ensure data directory exists
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	// Open database connection
	dbPath := filepath.Join(dataDir, "gorconx.db")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create tables
	if err := createTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

// createTables creates all necessary database tables
func createTables(db *sql.DB) error {
	schemas := []string{
		createSessionsTable,
		createSessionResultsTable,
		createModuleConfigsTable,
		createAPIKeysTable,
		createScansTable,
		createScanResultsTable,
		createReportsTable,
		createAuditLogsTable,
	}

	for _, schema := range schemas {
		if _, err := db.Exec(schema); err != nil {
			return err
		}
	}

	return nil
}

// Table schemas
const createSessionsTable = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    target TEXT NOT NULL,
    status TEXT DEFAULT 'created',
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    completed_at INTEGER,
    options TEXT, -- JSON
    metadata TEXT  -- JSON
);`

const createSessionResultsTable = `
CREATE TABLE IF NOT EXISTS session_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT NOT NULL,
    module_name TEXT NOT NULL,
    result_type TEXT NOT NULL,
    data TEXT NOT NULL, -- JSON
    metadata TEXT, -- JSON
    created_at INTEGER NOT NULL,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);`

const createModuleConfigsTable = `
CREATE TABLE IF NOT EXISTS module_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    module_name TEXT UNIQUE NOT NULL,
    config TEXT NOT NULL, -- JSON
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);`

const createAPIKeysTable = `
CREATE TABLE IF NOT EXISTS api_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    service_name TEXT UNIQUE NOT NULL,
    encrypted_key TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL,
    last_used INTEGER
);`

const createScansTable = `
CREATE TABLE IF NOT EXISTS scans (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    module_name TEXT NOT NULL,
    target TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    progress REAL DEFAULT 0.0,
    started_at INTEGER,
    completed_at INTEGER,
    error_message TEXT,
    options TEXT, -- JSON
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);`

const createScanResultsTable = `
CREATE TABLE IF NOT EXISTS scan_results (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scan_id TEXT NOT NULL,
    result_type TEXT NOT NULL,
    data TEXT NOT NULL, -- JSON
    metadata TEXT, -- JSON
    created_at INTEGER NOT NULL,
    FOREIGN KEY (scan_id) REFERENCES scans(id) ON DELETE CASCADE
);`

const createReportsTable = `
CREATE TABLE IF NOT EXISTS reports (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    name TEXT NOT NULL,
    format TEXT NOT NULL, -- html, pdf, json
    file_path TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    metadata TEXT, -- JSON
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);`

const createAuditLogsTable = `
CREATE TABLE IF NOT EXISTS audit_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    action TEXT NOT NULL,
    user_id TEXT,
    session_id TEXT,
    module_name TEXT,
    target TEXT,
    details TEXT, -- JSON
    ip_address TEXT,
    user_agent TEXT,
    created_at INTEGER NOT NULL
);`
