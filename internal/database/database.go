package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

// DB wraps the database connection and provides methods for data operations
type DB struct {
	*sql.DB
}

// InitDB initializes the SQLite database with required tables
func InitDB() (*DB, error) {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, "goreconx.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	dbInstance := &DB{db}
	if err := dbInstance.createTables(); err != nil {
		return nil, err
	}

	return dbInstance, nil
}

// createTables creates all necessary database tables
func (db *DB) createTables() error {
	queries := []string{
		// Projects table
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			target TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Scans table
		`CREATE TABLE IF NOT EXISTS scans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			scan_type TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			target TEXT NOT NULL,
			results TEXT,
			started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			completed_at DATETIME,
			error_message TEXT,
			FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE
		)`,
		
		// API Keys table (encrypted)
		`CREATE TABLE IF NOT EXISTS api_keys (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			service_name TEXT NOT NULL UNIQUE,
			encrypted_key TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Sessions table
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_name TEXT NOT NULL,
			session_data TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Wordlists table
		`CREATE TABLE IF NOT EXISTS wordlists (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			file_path TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Results table for structured storage
		`CREATE TABLE IF NOT EXISTS results (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			scan_id INTEGER NOT NULL,
			result_type TEXT NOT NULL,
			data TEXT NOT NULL,
			metadata TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (scan_id) REFERENCES scans (id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

// Project represents a reconnaissance project
type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Target      string `json:"target"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Scan represents a scan operation
type Scan struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"project_id"`
	ScanType     string `json:"scan_type"`
	Status       string `json:"status"`
	Target       string `json:"target"`
	Results      string `json:"results"`
	StartedAt    string `json:"started_at"`
	CompletedAt  string `json:"completed_at"`
	ErrorMessage string `json:"error_message"`
}

// StoreEncryptedAPIKey stores an API key in encrypted form
func (db *DB) StoreEncryptedAPIKey(serviceName, apiKey string) error {
	hashedKey, err := bcrypt.GenerateFromPassword([]byte(apiKey), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT OR REPLACE INTO api_keys (service_name, encrypted_key, updated_at) 
			  VALUES (?, ?, CURRENT_TIMESTAMP)`
	_, err = db.Exec(query, serviceName, string(hashedKey))
	return err
}

// CreateProject creates a new project
func (db *DB) CreateProject(name, description, target string) (*Project, error) {
	query := `INSERT INTO projects (name, description, target) VALUES (?, ?, ?)`
	result, err := db.Exec(query, name, description, target)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Project{
		ID:          int(id),
		Name:        name,
		Description: description,
		Target:      target,
	}, nil
}

// GetProjects returns all projects
func (db *DB) GetProjects() ([]*Project, error) {
	query := `SELECT id, name, description, target, created_at, updated_at FROM projects ORDER BY updated_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		p := &Project{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Target, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// CreateScan creates a new scan record
func (db *DB) CreateScan(projectID int, scanType, target string) (*Scan, error) {
	query := `INSERT INTO scans (project_id, scan_type, target) VALUES (?, ?, ?)`
	result, err := db.Exec(query, projectID, scanType, target)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Scan{
		ID:        int(id),
		ProjectID: projectID,
		ScanType:  scanType,
		Target:    target,
		Status:    "pending",
	}, nil
}

// UpdateScanStatus updates the status of a scan
func (db *DB) UpdateScanStatus(scanID int, status string, results string, errorMessage string) error {
	query := `UPDATE scans SET status = ?, results = ?, error_message = ?, 
			  completed_at = CASE WHEN ? IN ('completed', 'failed') THEN CURRENT_TIMESTAMP ELSE completed_at END 
			  WHERE id = ?`
	_, err := db.Exec(query, status, results, errorMessage, status, scanID)
	return err
}
