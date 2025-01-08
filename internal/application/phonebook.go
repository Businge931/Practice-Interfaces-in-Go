package application

import (
	"github.com/Businge931/practice-interfaces/internal/domain"
	"github.com/Businge931/practice-interfaces/internal/ports"
)

type PhonebookService struct {
	db ports.Database
}

func NewPhonebookService(db ports.Database) *PhonebookService {
	return &PhonebookService{db: db}
}

func (s *PhonebookService) AddContact(location string, contact domain.Contact) (bool, string) {
	// Validate the contact
	if err := s.ValidateContact(contact); err != nil {
		return false, err.Error()
	}

	// Convert contact to a map for storage
	contactData := map[string]interface{}{
		"name":    contact.Name,
		"phone":   contact.Phone,
		"email":   contact.Email,
		"address": contact.Address,
	}

	// Call the database's Create method
	return s.db.Create(location, contactData)
}

func (s *PhonebookService) GetContact(id string) (bool, string, domain.Contact) {
	// Call the database's Read method
	success, message, data := s.db.Read(id)
	if !success {
		return false, message, domain.Contact{}
	}

	// Convert the map to a Contact struct
	contact := domain.Contact{
		Name:    data["name"].(string),
		Phone:   data["phone"].(string),
		Email:   data["email"].(string),
		Address: data["address"].(string),
	}

	return true, "", contact
}

func (s *PhonebookService) UpdateContact(id string, contact domain.Contact) (bool, string) {
	// Validate the contact
	if err := s.ValidateContact(contact); err != nil {
		return false, err.Error()
	}

	// Convert contact to a map for storage
	contactData := map[string]interface{}{
		"name":    contact.Name,
		"phone":   contact.Phone,
		"email":   contact.Email,
		"address": contact.Address,
	}

	// Call the database's Update method
	return s.db.Update(id, contactData)
}

func (s *PhonebookService) DeleteContact(id string) (bool, string) {
	// Call the database's Delete method
	return s.db.Delete(id)
}

func (s *PhonebookService) ValidateContact(contact domain.Contact) error {
	if contact.Name == ""  {
		return domain.ErrInvalidContactName
	}
	if contact.Phone==""{
		return domain.ErrInvalidContactNumber
	}
	return nil
}
