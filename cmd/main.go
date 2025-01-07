package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/Businge931/practice-interfaces/internal/adoptors/database"
	"github.com/Businge931/practice-interfaces/internal/application"
	"github.com/Businge931/practice-interfaces/internal/domain"
)

func main() {
	// Connect to PostgreSQL Database
	pgDb, err := database.NewPostgresDatabase("postgres://user:pass@localhost:5432/phonebook?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgDb.Close()

	// Initialize phonebook service
	phonebook := application.NewPhonebookService(pgDb)

	// Test all CRUD operations
	
	// 1. Create a contact
	contact := domain.Contact{
		Name:    "John Doe",
		Phone:   "123-456-7890",
		Email:   "johndoe@example.com",
		Address: "123 Main St",
	}
	success, msg := phonebook.AddContact("contacts/johndoe", contact)
	if !success {
		log.Printf("Error adding contact: %v\n", msg)
	} else {
		log.Println("Contact added successfully")
	}

	// 2. Read the contact
	success, msg, retrievedContact := phonebook.GetContact("contacts/johndoe")
	if success {
		log.Printf("Retrieved contact: %+v\n", retrievedContact)
	} else {
		log.Printf("Error retrieving contact: %v\n", msg)
	}

	// 3. Update the contact
	updatedContact := domain.Contact{
		Name:    "John Doe Jr",
		Phone:   "999-999-9999",
		Email:   "john.jr@example.com",
		Address: "456 Oak St",
	}
	success, msg = phonebook.UpdateContact("contacts/johndoe", updatedContact)
	if success {
		log.Println("Contact updated successfully")
	} else {
		log.Printf("Error updating contact: %v\n", msg)
	}

	// 4. Read the updated contact
	success, msg, retrievedContact = phonebook.GetContact("contacts/johndoe")
	if success {
		log.Printf("Retrieved updated contact: %+v\n", retrievedContact)
	} else {
		log.Printf("Error retrieving contact: %v\n", msg)
	}

	// // 5. Delete the contact
	// success, msg = phonebook.DeleteContact("contacts/johndoe")
	// if success {
	// 	log.Println("Contact deleted successfully")
	// } else {
	// 	log.Printf("Error deleting contact: %v\n", msg)
	// }

	// // 6. Try to read the deleted contact (should fail)
	// success, msg, _ = phonebook.GetContact("contacts/johndoe")
	// if !success {
	// 	log.Printf("As expected, contact not found: %v\n", msg)
	// } else {
	// 	log.Println("Error: Contact still exists after deletion!")
	// }
}
