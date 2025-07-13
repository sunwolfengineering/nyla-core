package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestMigrationRunner(t *testing.T) {
	// Create temporary database file
	dbPath := "test_migrate.db"
	defer os.Remove(dbPath)
	defer os.Remove(dbPath + "-shm")
	defer os.Remove(dbPath + "-wal")
	
	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()
	
	// Create migration runner
	runner := NewMigrationRunner(db)
	
	// Create temporary migrations directory
	tempDir, err := os.MkdirTemp("", "migrations")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test migration files
	migration1 := `
		CREATE TABLE test_table1 (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
		INSERT INTO schema_migrations (version) VALUES (1);
	`
	
	migration2 := `
		CREATE TABLE test_table2 (
			id INTEGER PRIMARY KEY,
			description TEXT
		);
		INSERT INTO schema_migrations (version) VALUES (2);
	`
	
	err = os.WriteFile(filepath.Join(tempDir, "001_create_test_table1.sql"), []byte(migration1), 0644)
	require.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(tempDir, "002_create_test_table2.sql"), []byte(migration2), 0644)
	require.NoError(t, err)
	
	// Run migrations
	err = runner.Run(tempDir)
	require.NoError(t, err)
	
	// Verify migrations table was created
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
	
	// Verify test tables were created
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_table1'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='test_table2'").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
	
	// Run migrations again (should be no-op)
	err = runner.Run(tempDir)
	require.NoError(t, err)
	
	// Should still have 2 migrations
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestParseMigrationFileName(t *testing.T) {
	runner := &MigrationRunner{}
	
	tests := []struct {
		fileName     string
		expectedVersion int
		expectedName string
		shouldError  bool
	}{
		{"001_initial_schema.sql", 1, "initial_schema", false},
		{"00002_create_users.sql", 2, "create_users", false},
		{"123_add_indexes.sql", 123, "add_indexes", false},
		{"invalid.sql", 0, "", true},
		{"abc_invalid.sql", 0, "", true},
	}
	
	for _, test := range tests {
		migration, err := runner.parseMigrationFileName(test.fileName, "/fake/path/"+test.fileName)
		
		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedVersion, migration.Version)
			assert.Equal(t, test.expectedName, migration.Name)
		}
	}
}

func TestGetCurrentVersion(t *testing.T) {
	// Create temporary database file
	dbPath := "test_version.db"
	defer os.Remove(dbPath)
	defer os.Remove(dbPath + "-shm")
	defer os.Remove(dbPath + "-wal")
	
	// Open database connection
	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)
	defer db.Close()
	
	runner := NewMigrationRunner(db)
	
	// Should return 0 when no migrations table exists
	version, err := runner.getCurrentVersion()
	require.Error(t, err) // Should error when table doesn't exist
	
	// Create migrations table
	err = runner.createMigrationsTable()
	require.NoError(t, err)
	
	// Should still return 0 when table is empty
	version, err = runner.getCurrentVersion()
	require.NoError(t, err)
	assert.Equal(t, 0, version)
	
	// Insert some migrations
	_, err = db.Exec("INSERT INTO schema_migrations (version) VALUES (1), (3), (2)")
	require.NoError(t, err)
	
	// Should return highest version
	version, err = runner.getCurrentVersion()
	require.NoError(t, err)
	assert.Equal(t, 3, version)
}
