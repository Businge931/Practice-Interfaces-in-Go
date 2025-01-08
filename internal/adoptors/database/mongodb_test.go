package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// "go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoTestDeps struct {
	db *MongoDatabase
}

type MongoCRUDArgs struct {
	location string
	data     map[string]interface{}
}

type MongoCRUDTest struct {
	name     string
	deps     MongoTestDeps
	args     MongoCRUDArgs
	before   func(*testing.T, MongoTestDeps)
	after    func(*testing.T, MongoTestDeps)
	expected map[string]interface{}
	wantErr  string
}

func setupMongoTest(t *testing.T) *MongoDatabase {
	db, err := NewMongoDatabase(
		"mongodb://user:pass@localhost:27017/phonebook",
		"phonebook_test",
		"contacts",
	)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func cleanupMongoTest(t *testing.T, db *MongoDatabase) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := db.collection.Drop(ctx)
	if err != nil {
		t.Errorf("Failed to cleanup test collection: %v", err)
	}
	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close database connection: %v", err)
	}
}

func TestMongoDatabase_Create(t *testing.T) {
	db := setupMongoTest(t)
	defer cleanupMongoTest(t, db)

	tests := []MongoCRUDTest{
		{
			name: "Create Valid Contact",
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
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
		},
		{
			name: "Create Duplicate Contact",
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
				location: "contacts/test1",
				data: map[string]interface{}{
					"name": "Jane Doe",
				},
			},
			before: func(t *testing.T, deps MongoTestDeps) {
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

func TestMongoDatabase_Read(t *testing.T) {
	db := setupMongoTest(t)
	defer cleanupMongoTest(t, db)

	tests := []MongoCRUDTest{
		{
			name: "Read Existing Contact",
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
				location: "contacts/test1",
			},
			before: func(t *testing.T, deps MongoTestDeps) {
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
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
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

func TestMongoDatabase_Update(t *testing.T) {
	db := setupMongoTest(t)
	defer cleanupMongoTest(t, db)

	tests := []MongoCRUDTest{
		{
			name: "Update Existing Contact",
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
				location: "contacts/test1",
				data: map[string]interface{}{
					"name":    "John Updated",
					"phone":   "999-999-9999",
					"email":   "john.updated@example.com",
					"address": "456 Oak St",
				},
			},
			before: func(t *testing.T, deps MongoTestDeps) {
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
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
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

func TestMongoDatabase_Delete(t *testing.T) {
	db := setupMongoTest(t)
	defer cleanupMongoTest(t, db)

	tests := []MongoCRUDTest{
		{
			name: "Delete Existing Contact",
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
				location: "contacts/test1",
			},
			before: func(t *testing.T, deps MongoTestDeps) {
				data := map[string]interface{}{
					"name": "John Doe",
				}
				success, _ := deps.db.Create("contacts/test1", data)
				assert.True(t, success, "Failed to create test contact")
			},
			after: func(t *testing.T, deps MongoTestDeps) {
				success, _, _ := deps.db.Read("contacts/test1")
				assert.False(t, success, "Contact should not exist after deletion")
			},
		},
		{
			name: "Delete Non-existent Contact",
			deps: MongoTestDeps{db: db},
			args: MongoCRUDArgs{
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

func TestMongoDatabase_Concurrency(t *testing.T) {
	db := setupMongoTest(t)
	defer cleanupMongoTest(t, db)

	// Create initial contact
	success, err := db.Create("contacts/concurrent", map[string]interface{}{
		"name":    "Concurrent Test",
		"counter": int64(0),
	})
	assert.True(t, success, "Failed to create initial contact: %v", err)

	// Test concurrent updates
	const numGoroutines = 10
	done := make(chan bool)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			success, _, data := db.Read("contacts/concurrent")
			assert.True(t, success, "Failed to read contact")

			counter := data["counter"].(int64)
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
	assert.Equal(t, int64(numGoroutines), data["counter"].(int64), "Counter value mismatch")
}
