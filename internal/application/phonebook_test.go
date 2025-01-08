package application

import (
	"testing"

	"github.com/Businge931/practice-interfaces/internal/domain"
	"github.com/Businge931/practice-interfaces/internal/ports"
)

// MockDatabase implements ports.Database interface for testing
type MockDatabase struct {
	createFunc func(location string, data map[string]interface{}) (bool, string)
	readFunc   func(location string) (bool, string, map[string]interface{})
	updateFunc func(location string, data map[string]interface{}) (bool, string)
	deleteFunc func(location string) (bool, string)
}

func (m *MockDatabase) Create(location string, data map[string]interface{}) (bool, string) {
	return m.createFunc(location, data)
}

func (m *MockDatabase) Read(location string) (bool, string, map[string]interface{}) {
	return m.readFunc(location)
}

func (m *MockDatabase) Update(location string, data map[string]interface{}) (bool, string) {
	return m.updateFunc(location, data)
}

func (m *MockDatabase) Delete(location string) (bool, string) {
	return m.deleteFunc(location)
}

type phonebookTestCase struct {
	name string
	db   ports.Database
	args struct {
		location string
		contact  domain.Contact
	}
	want    domain.Contact
	wantErr bool
	errMsg  string
}

func TestPhonebookService_AddContact(t *testing.T) {
	tests := []phonebookTestCase{
		{
			name: "successful add contact",
			db: &MockDatabase{
				createFunc: func(location string, data map[string]interface{}) (bool, string) {
					return true, ""
				},
			},
			args: struct {
				location string
				contact  domain.Contact
			}{
				location: "contacts/john.json",
				contact: domain.Contact{
					Name:    "John Doe",
					Phone:   "123-456-7890",
					Email:   "john@example.com",
					Address: "123 Main St",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid contact - empty name",
			db: &MockDatabase{
				createFunc: func(location string, data map[string]interface{}) (bool, string) {
					return false, "Invalid contact"
				},
			},
			args: struct {
				location string
				contact  domain.Contact
			}{
				location: "contacts/invalid.json",
				contact: domain.Contact{
					Phone:   "123-456-7890",
					Email:   "invalid@example.com",
					Address: "123 Main St",
				},
			},
			wantErr: true,
			errMsg:  "Name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPhonebookService(tt.db)
			success, msg := s.AddContact(tt.args.location, tt.args.contact)

			if tt.wantErr {
				if success {
					t.Errorf("Expected error but got success")
				}
			} else {
				if !success {
					t.Errorf("Expected success but got error: %s", msg)
				}
			}
		})
	}
}

func TestPhonebookService_GetContact(t *testing.T) {
	tests := []phonebookTestCase{
		{
			name: "successful get contact",
			db: &MockDatabase{
				readFunc: func(location string) (bool, string, map[string]interface{}) {
					return true, "", map[string]interface{}{
						"name":    "John Doe",
						"phone":   "123-456-7890",
						"email":   "john@example.com",
						"address": "123 Main St",
					}
				},
			},
			args: struct {
				location string
				contact  domain.Contact
			}{
				location: "contacts/john.json",
			},
			want: domain.Contact{
				Name:    "John Doe",
				Phone:   "123-456-7890",
				Email:   "john@example.com",
				Address: "123 Main St",
			},
			wantErr: false,
		},
		{
			name: "contact not found",
			db: &MockDatabase{
				readFunc: func(location string) (bool, string, map[string]interface{}) {
					return false, "File does not exist", nil
				},
			},
			args: struct {
				location string
				contact  domain.Contact
			}{
				location: "contacts/nonexistent.json",
			},
			wantErr: true,
			errMsg:  "File does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewPhonebookService(tt.db)
			success, msg, contact := s.GetContact(tt.args.location)

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
				if contact != tt.want {
					t.Errorf("Expected contact %+v but got %+v", tt.want, contact)
				}
			}
		})
	}
}

func TestPhonebookService_ValidateContact(t *testing.T) {
	type validateTestCase struct {
		name    string
		contact domain.Contact
		wantErr bool
		errMsg  string
	}

	tests := []validateTestCase{
		{
			name: "valid contact",
			contact: domain.Contact{
				Name:    "John Doe",
				Phone:   "123-456-7890",
				Email:   "john@example.com",
				Address: "123 Main St",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			contact: domain.Contact{
				Phone:   "123-456-7890",
				Email:   "john@example.com",
				Address: "123 Main St",
			},
			wantErr: true,
			errMsg:  "invalid contact: Name is required",
		},
		{
			name: "empty phone",
			contact: domain.Contact{
				Name:    "John Doe",
				Email:   "john@example.com",
				Address: "123 Main St",
			},
			wantErr: true,
			errMsg:  "invalid contact: Phone is required",
		},
	}

	s := NewPhonebookService(&MockDatabase{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.ValidateContact(tt.contact)
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Expected error message %q but got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}
