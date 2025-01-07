package database

import (
	"fmt"
	"sync"
	"testing"
)

type memoryTestCase struct {
	name     string
	setup    func(t *testing.T, db *InMemoryDatabase) map[string]interface{}
	teardown func(t *testing.T, db *InMemoryDatabase)
	args     struct {
		location string
		data     map[string]interface{}
	}
	want    map[string]interface{}
	wantErr bool
	errMsg  string
}

func TestInMemoryDatabase_Create(t *testing.T) {
	tests := []memoryTestCase{
		{
			name: "successful create",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/john.json",
				data: map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				},
			},
			wantErr: false,
		},
		{
			name: "location already exists",
			setup: func(t *testing.T, db *InMemoryDatabase) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "Existing Contact",
					"phone": "000-000-0000",
				}
				db.Create("contacts/existing.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/existing.json",
				data: map[string]interface{}{
					"name":  "New Contact",
					"phone": "111-111-1111",
				},
			},
			wantErr: true,
			errMsg:  "Location already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDatabase()
			if tt.setup != nil {
				tt.setup(t, db)
			}

			success, msg := db.Create(tt.args.location, tt.args.data)
			if tt.wantErr {
				if success {
					t.Errorf("Expected error but got success")
				}
				if msg != tt.errMsg {
					t.Errorf("Expected error message %q but got %q", tt.errMsg, msg)
				}
			} else {
				if !success {
					t.Errorf("Expected success but got error: %s", msg)
				}
				// Verify data was stored correctly
				if data, exists := db.store[tt.args.location]; !exists {
					t.Error("Data was not stored")
				} else {
					for k, v := range tt.args.data {
						if data[k] != v {
							t.Errorf("Expected %v for key %s but got %v", v, k, data[k])
						}
					}
				}
			}

			if tt.teardown != nil {
				tt.teardown(t, db)
			}
		})
	}
}

func TestInMemoryDatabase_Read(t *testing.T) {
	tests := []memoryTestCase{
		{
			name: "successful read",
			setup: func(t *testing.T, db *InMemoryDatabase) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				}
				db.Create("contacts/readable.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/readable.json",
			},
			want:    map[string]interface{}{"name": "John Doe", "phone": "123-456-7890"},
			wantErr: false,
		},
		{
			name: "location not found",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/nonexistent.json",
			},
			wantErr: true,
			errMsg:  "Location does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDatabase()
			if tt.setup != nil {
				tt.want = tt.setup(t, db)
			}

			success, msg, data := db.Read(tt.args.location)
			if tt.wantErr {
				if success {
					t.Errorf("Expected error but got success")
				}
				if msg != tt.errMsg {
					t.Errorf("Expected error message %q but got %q", tt.errMsg, msg)
				}
			} else {
				if !success {
					t.Errorf("Expected success but got error: %s", msg)
				}
				for k, v := range tt.want {
					if data[k] != v {
						t.Errorf("Expected %v for key %s but got %v", v, k, data[k])
					}
				}
			}

			if tt.teardown != nil {
				tt.teardown(t, db)
			}
		})
	}
}

func TestInMemoryDatabase_Update(t *testing.T) {
	tests := []memoryTestCase{
		{
			name: "successful update",
			setup: func(t *testing.T, db *InMemoryDatabase) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				}
				db.Create("contacts/updatable.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/updatable.json",
				data: map[string]interface{}{
					"name":  "John Doe Updated",
					"phone": "999-999-9999",
				},
			},
			wantErr: false,
		},
		{
			name: "update nonexistent location",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/nonexistent.json",
				data: map[string]interface{}{
					"name":  "New Contact",
					"phone": "111-111-1111",
				},
			},
			wantErr: true,
			errMsg:  "Location does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDatabase()
			if tt.setup != nil {
				tt.setup(t, db)
			}

			success, msg := db.Update(tt.args.location, tt.args.data)
			if tt.wantErr {
				if success {
					t.Errorf("Expected error but got success")
				}
				if msg != tt.errMsg {
					t.Errorf("Expected error message %q but got %q", tt.errMsg, msg)
				}
			} else {
				if !success {
					t.Errorf("Expected success but got error: %s", msg)
				}
				// Verify data was updated correctly
				if data, exists := db.store[tt.args.location]; !exists {
					t.Error("Data was not stored")
				} else {
					for k, v := range tt.args.data {
						if data[k] != v {
							t.Errorf("Expected %v for key %s but got %v", v, k, data[k])
						}
					}
				}
			}

			if tt.teardown != nil {
				tt.teardown(t, db)
			}
		})
	}
}

func TestInMemoryDatabase_Delete(t *testing.T) {
	tests := []memoryTestCase{
		{
			name: "successful delete",
			setup: func(t *testing.T, db *InMemoryDatabase) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				}
				db.Create("contacts/deletable.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/deletable.json",
			},
			wantErr: false,
		},
		{
			name: "delete nonexistent location",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "contacts/nonexistent.json",
			},
			wantErr: true,
			errMsg:  "Location does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := NewInMemoryDatabase()
			if tt.setup != nil {
				tt.setup(t, db)
			}

			success, msg := db.Delete(tt.args.location)
			if tt.wantErr {
				if success {
					t.Errorf("Expected error but got success")
				}
				if msg != tt.errMsg {
					t.Errorf("Expected error message %q but got %q", tt.errMsg, msg)
				}
			} else {
				if !success {
					t.Errorf("Expected success but got error: %s", msg)
				}
				// Verify data was actually deleted
				if _, exists := db.store[tt.args.location]; exists {
					t.Error("Data was not deleted")
				}
			}

			if tt.teardown != nil {
				tt.teardown(t, db)
			}
		})
	}
}

func TestInMemoryDatabase_Concurrency(t *testing.T) {
	db := NewInMemoryDatabase()
	const numOperations = 100
	var wg sync.WaitGroup

	// Test concurrent reads and writes
	wg.Add(3)

	// Concurrent writes
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			db.Create(fmt.Sprintf("location%d", i), map[string]interface{}{"value": i})
		}
	}()

	// Concurrent reads
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			db.Read(fmt.Sprintf("location%d", i))
		}
	}()

	// Concurrent updates
	go func() {
		defer wg.Done()
		for i := 0; i < numOperations; i++ {
			db.Update(fmt.Sprintf("location%d", i), map[string]interface{}{"value": i * 2})
		}
	}()

	wg.Wait()
}
