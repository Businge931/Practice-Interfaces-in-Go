package database

import (
	"os"
	"path/filepath"
	"testing"
)

type testCase struct {
	name     string
	setup    func(t *testing.T, baseDir string) map[string]interface{}
	teardown func(t *testing.T, baseDir string)
	args     struct {
		location string
		data     map[string]interface{}
	}
	want    map[string]interface{}
	wantErr bool
	errMsg  string
}

func TestFileSystemDatabase_Create(t *testing.T) {
	baseDir := t.TempDir()
	db := NewFileSystemDatabase(baseDir)

	tests := []testCase{
		{
			name: "successful create",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/contact.json",
				data: map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				},
			},
			wantErr: false,
		},
		{
			name: "file already exists",
			setup: func(t *testing.T, baseDir string) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "Existing Contact",
					"phone": "000-000-0000",
				}
				db.Create("test/existing.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/existing.json",
				data: map[string]interface{}{
					"name":  "New Contact",
					"phone": "111-111-1111",
				},
			},
			wantErr: true,
			errMsg:  "File already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, baseDir)
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
			}

			if tt.teardown != nil {
				tt.teardown(t, baseDir)
			}
		})
	}
}

func TestFileSystemDatabase_Read(t *testing.T) {
	baseDir := t.TempDir()
	db := NewFileSystemDatabase(baseDir)

	tests := []testCase{
		{
			name: "successful read",
			setup: func(t *testing.T, baseDir string) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				}
				db.Create("test/readable.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/readable.json",
			},
			want:    map[string]interface{}{"name": "John Doe", "phone": "123-456-7890"},
			wantErr: false,
		},
		{
			name: "file not found",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/nonexistent.json",
			},
			wantErr: true,
			errMsg:  "File does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.want = tt.setup(t, baseDir)
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
				tt.teardown(t, baseDir)
			}
		})
	}
}

func TestFileSystemDatabase_Update(t *testing.T) {
	baseDir := t.TempDir()
	db := NewFileSystemDatabase(baseDir)

	tests := []testCase{
		{
			name: "successful update",
			setup: func(t *testing.T, baseDir string) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				}
				db.Create("test/updatable.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/updatable.json",
				data: map[string]interface{}{
					"name":  "John Doe Updated",
					"phone": "999-999-9999",
				},
			},
			wantErr: false,
		},
		{
			name: "update nonexistent file",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/nonexistent.json",
				data: map[string]interface{}{
					"name":  "New Contact",
					"phone": "111-111-1111",
				},
			},
			wantErr: true,
			errMsg:  "File does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, baseDir)
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
			}

			if tt.teardown != nil {
				tt.teardown(t, baseDir)
			}
		})
	}
}

func TestFileSystemDatabase_Delete(t *testing.T) {
	baseDir := t.TempDir()
	db := NewFileSystemDatabase(baseDir)

	tests := []testCase{
		{
			name: "successful delete",
			setup: func(t *testing.T, baseDir string) map[string]interface{} {
				data := map[string]interface{}{
					"name":  "John Doe",
					"phone": "123-456-7890",
				}
				db.Create("test/deletable.json", data)
				return data
			},
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/deletable.json",
			},
			wantErr: false,
		},
		{
			name: "delete nonexistent file",
			args: struct {
				location string
				data     map[string]interface{}
			}{
				location: "test/nonexistent.json",
			},
			wantErr: true,
			errMsg:  "File does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(t, baseDir)
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
				// Verify file is actually deleted
				if _, err := os.Stat(filepath.Join(baseDir, tt.args.location)); !os.IsNotExist(err) {
					t.Errorf("File was not deleted")
				}
			}

			if tt.teardown != nil {
				tt.teardown(t, baseDir)
			}
		})
	}
}
