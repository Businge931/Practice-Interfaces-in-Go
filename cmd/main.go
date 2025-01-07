package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/Businge931/practice-interfaces/internal/adoptors/database"
	"github.com/Businge931/practice-interfaces/internal/application"
	"github.com/Businge931/practice-interfaces/internal/domain"
)

func main() {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}
	// Initialize the database adapter (driven adapter)
	db := database.NewFileSystemDatabase(filepath.Join(homeDir, "data"))

	// db:=database.NewInMemoryDatabase()


	//Initialize in application layer
	phonebook := application.NewPhonebookService(db)

	// Example usage
	contact := domain.Contact{
		Name:    "John Doeson",
		Phone:   "123-456-7890",
		Email:   "johndoe@example.com",
		Address: "Bukoto",
	}

	// create a new contact
	hasAdded, phoneNumber := phonebook.AddContact("contacts/johndoe.json", contact)
	if !hasAdded {
		log.Printf("Error adding contact: %v\n", phoneNumber)
	} else {
		log.Println("Contact added successfully")
	}

	// Retrieve the contact
	success, msg, retrievedContact := phonebook.GetContact("contacts/johndoe.json")
	if success {
		log.Printf("Retrieved contact: %+v\n", retrievedContact)
	} else {
		log.Printf("Error retrieving contact: %v\n", msg)
	}

	// Update the contact
	updatedContact := domain.Contact{
		Name:    "John Doe",
		Phone:   "123-456-7890",
		Email:   "johndoe@example.com",
		Address: "Bukoto",
	}
	success, msg = phonebook.UpdateContact("contacts/johndoe.json", updatedContact)
	if success {
		log.Println("Contact updated successfully")
	} else {
		log.Printf("Error updating contact: %v\n", msg)
	}

	// Delete the contact
	success, msg = phonebook.DeleteContact("contacts/johndoe.json")
	if success {
		log.Println("Contact deleted successfully")
	} else {
		log.Printf("Error deleting contact: %v\n", msg)
	}
}
