package database

import (
	"context"
	"database/sql"
	
	"testing"

	"github.com/stretchr/testify/assert"
	
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresTestDeps struct {
	db *PostgresDatabase
}

type PostgresCRUDArgs struct {
	location string
	data     map[string]interface{}
}

type PostgresCRUDTest struct {
	name     string
	deps     PostgresTestDeps
	args     PostgresCRUDArgs
	before   func(*testing.T, PostgresTestDeps)
	after    func(*testing.T, PostgresTestDeps)
	expected map[string]interface{}
	wantErr  string
}

func setupPostgresTest(t *testing.T) *PostgresDatabase {
	// Create test database if it doesn't exist
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN("postgres://user:pass@localhost:5432/postgres?sslmode=disable")))
	defer sqldb.Close()

	// Drop test database if it exists
	_, err := sqldb.Exec("DROP DATABASE IF EXISTS phonebook_test")
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}

	// Create test database
	_, err = sqldb.Exec("CREATE DATABASE phonebook_test")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Connect to test database
	db, err := NewPostgresDatabase("postgres://user:pass@localhost:5432/phonebook_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func cleanupPostgresTest(t *testing.T, db *PostgresDatabase) {
	ctx := context.Background()
	_, err := db.db.NewDropTable().Model((*Contact)(nil)).IfExists().Exec(ctx)
	if err != nil {
		t.Errorf("Failed to cleanup test database: %v", err)
	}
	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close database connection: %v", err)
	}
}

func TestPostgresDatabase_Create(t *testing.T) {
	db := setupPostgresTest(t)
	defer cleanupPostgresTest(t, db)

	tests := []PostgresCRUDTest{
		{
			name: "Create Valid Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/test1",
				data: map[string]interface{}{
					"name":    "John Doe",
					"phone":   "123-456-7890",
					"email":   "john@example.com",
					"address": "123 Main St",
				},
			},
			expected: map[string]interface{}{
				"name":    "John Doe",
				"phone":   "123-456-7890",
				"email":   "john@example.com",
				"address": "123 Main St",
			},
			wantErr: "",
		},
		{
			name: "Create Duplicate Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/test1",
				data: map[string]interface{}{
					"name": "Jane Doe",
				},
			},
			before: func(t *testing.T, deps PostgresTestDeps) {
				success, _ := deps.db.Create("contacts/test1", map[string]interface{}{"name": "Original"})
				assert.True(t, success, "Failed to create initial contact")
			},
			wantErr: "Location already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t, tt.deps)
			}

			success, err := tt.deps.db.Create(tt.args.location, tt.args.data)

			if tt.after != nil {
				tt.after(t, tt.deps)
			}

			if tt.wantErr != "" {
				assert.False(t, success)
				assert.Contains(t, err, tt.wantErr)
				return
			}

			assert.True(t, success)
			assert.Empty(t, err)

			// Verify created data
			_, _, data := tt.deps.db.Read(tt.args.location)
			assert.Equal(t, tt.expected, data)
		})
	}
}

func TestPostgresDatabase_Read(t *testing.T) {
	db := setupPostgresTest(t)
	defer cleanupPostgresTest(t, db)

	tests := []PostgresCRUDTest{
		{
			name: "Read Existing Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/test1",
			},
			before: func(t *testing.T, deps PostgresTestDeps) {
				data := map[string]interface{}{
					"name":    "John Doe",
					"phone":   "123-456-7890",
					"email":   "john@example.com",
					"address": "123 Main St",
				}
				success, _ := deps.db.Create("contacts/test1", data)
				assert.True(t, success, "Failed to create test contact")
			},
			expected: map[string]interface{}{
				"name":    "John Doe",
				"phone":   "123-456-7890",
				"email":   "john@example.com",
				"address": "123 Main St",
			},
		},
		{
			name: "Read Non-existent Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/nonexistent",
			},
			wantErr: "Location does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t, tt.deps)
			}

			success, err, data := tt.deps.db.Read(tt.args.location)

			if tt.after != nil {
				tt.after(t, tt.deps)
			}

			if tt.wantErr != "" {
				assert.False(t, success)
				assert.Contains(t, err, tt.wantErr)
				return
			}

			assert.True(t, success)
			assert.Empty(t, err)
			assert.Equal(t, tt.expected, data)
		})
	}
}

func TestPostgresDatabase_Update(t *testing.T) {
	db := setupPostgresTest(t)
	defer cleanupPostgresTest(t, db)

	tests := []PostgresCRUDTest{
		{
			name: "Update Existing Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/test1",
				data: map[string]interface{}{
					"name":    "John Updated",
					"phone":   "999-999-9999",
					"email":   "john.updated@example.com",
					"address": "456 Oak St",
				},
			},
			before: func(t *testing.T, deps PostgresTestDeps) {
				data := map[string]interface{}{
					"name":    "John Doe",
					"phone":   "123-456-7890",
					"email":   "john@example.com",
					"address": "123 Main St",
				}
				success, _ := deps.db.Create("contacts/test1", data)
				assert.True(t, success, "Failed to create test contact")
			},
			expected: map[string]interface{}{
				"name":    "John Updated",
				"phone":   "999-999-9999",
				"email":   "john.updated@example.com",
				"address": "456 Oak St",
			},
		},
		{
			name: "Update Non-existent Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/nonexistent",
				data: map[string]interface{}{
					"name": "Nobody",
				},
			},
			wantErr: "Location does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t, tt.deps)
			}

			success, err := tt.deps.db.Update(tt.args.location, tt.args.data)

			if tt.after != nil {
				tt.after(t, tt.deps)
			}

			if tt.wantErr != "" {
				assert.False(t, success)
				assert.Contains(t, err, tt.wantErr)
				return
			}

			assert.True(t, success)
			assert.Empty(t, err)

			// Verify updated data
			_, _, data := tt.deps.db.Read(tt.args.location)
			assert.Equal(t, tt.expected, data)
		})
	}
}

func TestPostgresDatabase_Delete(t *testing.T) {
	db := setupPostgresTest(t)
	defer cleanupPostgresTest(t, db)

	tests := []PostgresCRUDTest{
		{
			name: "Delete Existing Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/test1",
			},
			before: func(t *testing.T, deps PostgresTestDeps) {
				data := map[string]interface{}{
					"name": "John Doe",
				}
				success, _ := deps.db.Create("contacts/test1", data)
				assert.True(t, success, "Failed to create test contact")
			},
			after: func(t *testing.T, deps PostgresTestDeps) {
				success, _, _ := deps.db.Read("contacts/test1")
				assert.False(t, success, "Contact should not exist after deletion")
			},
		},
		{
			name: "Delete Non-existent Contact",
			deps: PostgresTestDeps{db: db},
			args: PostgresCRUDArgs{
				location: "contacts/nonexistent",
			},
			wantErr: "Location does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t, tt.deps)
			}

			success, err := tt.deps.db.Delete(tt.args.location)

			if tt.after != nil {
				tt.after(t, tt.deps)
			}

			if tt.wantErr != "" {
				assert.False(t, success)
				assert.Contains(t, err, tt.wantErr)
				return
			}

			assert.True(t, success)
			assert.Empty(t, err)
		})
	}
}

func TestPostgresDatabase_Concurrency(t *testing.T) {
	db := setupPostgresTest(t)
	defer cleanupPostgresTest(t, db)

	// Create initial contact
	success, err := db.Create("contacts/concurrent", map[string]interface{}{
		"name":    "Concurrent Test",
		"counter": 0,
	})
	assert.True(t, success, "Failed to create initial contact: %v", err)

	// Test concurrent updates
	const numGoroutines = 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			success, _, data := db.Read("contacts/concurrent")
			assert.True(t, success, "Failed to read contact")

			counter := data["counter"].(float64)
			data["counter"] = counter + 1

			success, _ = db.Update("contacts/concurrent", data)
			assert.True(t, success, "Failed to update contact")

			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify final counter value
	success, _, data := db.Read("contacts/concurrent")
	assert.True(t, success, "Failed to read final value")
	assert.Equal(t, float64(numGoroutines), data["counter"].(float64), "Counter value mismatch")
}
