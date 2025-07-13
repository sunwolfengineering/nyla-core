package storage

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db *sql.DB
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db *sql.DB) *MigrationRunner {
	return &MigrationRunner{db: db}
}

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	Path    string
}

// Run executes all pending migrations
func (m *MigrationRunner) Run(migrationsPath string) error {
	log.Printf("Running migrations from %s", migrationsPath)
	
	// Ensure migrations table exists
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	// Get current version
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	
	log.Printf("Current database version: %d", currentVersion)
	
	// Find migration files
	migrations, err := m.findMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to find migrations: %w", err)
	}
	
	// Filter migrations that need to be applied
	pendingMigrations := make([]Migration, 0)
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}
	
	if len(pendingMigrations) == 0 {
		log.Println("No pending migrations")
		return nil
	}
	
	log.Printf("Found %d pending migrations", len(pendingMigrations))
	
	// Apply migrations
	for _, migration := range pendingMigrations {
		if err := m.applyMigration(migration); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}
		log.Printf("Applied migration %d: %s", migration.Version, migration.Name)
	}
	
	log.Println("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the schema_migrations table if it doesn't exist
func (m *MigrationRunner) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`
	
	_, err := m.db.Exec(query)
	return err
}

// getCurrentVersion returns the current schema version
func (m *MigrationRunner) getCurrentVersion() (int, error) {
	var version int
	query := "SELECT COALESCE(MAX(version), 0) FROM schema_migrations"
	
	err := m.db.QueryRow(query).Scan(&version)
	if err != nil {
		return 0, err
	}
	
	return version, nil
}

// findMigrations discovers migration files in the given directory
func (m *MigrationRunner) findMigrations(migrationsPath string) ([]Migration, error) {
	var migrations []Migration
	
	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		
		if d.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}
		
		fileName := d.Name()
		migration, err := m.parseMigrationFileName(fileName, path)
		if err != nil {
			log.Printf("Skipping invalid migration file %s: %v", fileName, err)
			return nil
		}
		
		migrations = append(migrations, migration)
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	
	return migrations, nil
}

// parseMigrationFileName extracts version and name from migration file name
// Expected format: 001_initial_schema.sql or 00001_create_events.sql
func (m *MigrationRunner) parseMigrationFileName(fileName string, path string) (Migration, error) {
	// Remove .sql extension
	nameWithoutExt := strings.TrimSuffix(fileName, ".sql")
	
	// Find the first underscore to separate version from name
	underscoreIdx := strings.Index(nameWithoutExt, "_")
	if underscoreIdx == -1 {
		return Migration{}, fmt.Errorf("invalid migration file format: %s", fileName)
	}
	
	versionStr := nameWithoutExt[:underscoreIdx]
	name := nameWithoutExt[underscoreIdx+1:]
	
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return Migration{}, fmt.Errorf("invalid version number in %s: %w", fileName, err)
	}
	
	return Migration{
		Version: version,
		Name:    name,
		Path:    path,
	}, nil
}

// applyMigration executes a single migration file
func (m *MigrationRunner) applyMigration(migration Migration) error {
	// Read migration file
	content, err := os.ReadFile(migration.Path)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}
	
	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Execute migration SQL
	_, err = tx.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}
	
	// Record migration in schema_migrations table
	_, err = tx.Exec(
		"INSERT INTO schema_migrations (version) VALUES (?) ON CONFLICT DO NOTHING",
		migration.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}
	
	return nil
}

// Status returns the current migration status
func (m *MigrationRunner) Status(migrationsPath string) error {
	currentVersion, err := m.getCurrentVersion()
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}
	
	migrations, err := m.findMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to find migrations: %w", err)
	}
	
	fmt.Printf("Current database version: %d\n", currentVersion)
	fmt.Printf("Available migrations:\n")
	
	for _, migration := range migrations {
		status := "pending"
		if migration.Version <= currentVersion {
			status = "applied"
		}
		fmt.Printf("  %03d %s [%s]\n", migration.Version, migration.Name, status)
	}
	
	return nil
}
