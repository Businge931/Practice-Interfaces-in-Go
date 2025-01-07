package main

import (
	// "os"
	// "path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/Businge931/practice-interfaces/internal/adoptors/database"
	"github.com/Businge931/practice-interfaces/internal/application"
	"github.com/Businge931/practice-interfaces/internal/domain"
)

func main() {
	// Choose one of the following database implementations:

	// // 1. Filesystem Database
	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	log.Fatalf("Failed to get home directory: %v", err)
	// }
	// fsDb := database.NewFileSystemDatabase(filepath.Join(homeDir, "data"))

	// // 2. In-Memory Database
	// memDb := database.NewInMemoryDatabase()

	// 3. PostgreSQL Database
	pgDb, err := database.NewPostgresDatabase("postgres://user:pass@localhost:5432/phonebook?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgDb.Close()

	// // 4. MongoDB Database
	// // Using authentication credentials from docker-compose
	// mongoDb, err := database.NewMongoDatabase(
	// 	"mongodb://user:pass@localhost:27017/phonebook",
	// 	"phonebook",
	// 	"contacts",
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to connect to MongoDB: %v", err)
	// }
	// defer mongoDb.Close()

	// // 5. Google Cloud Firestore Database
	// firestoreDb, err := database.NewFirestoreDatabase(
	// 	"your-project-id",           // To be replace with my GCP project ID
	// 	"contacts",                  // Collection name
	// 	"path/to/credentials.json",  // To be replace with path to my service account key file
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to connect to Firestore: %v", err)
	// }
	// defer firestoreDb.Close()

	// Use one of the database implementations
	db := pgDb // or memDb, pgDb, mongoDb, firestoreDb

	// Initialize application layer
	phonebook := application.NewPhonebookService(db)

	// Example usage
	contact := domain.Contact{
		Name:    "John Doe",
		Phone:   "123-456-7890",
		Email:   "johndoe@example.com",
		Address: "123 Main St",
	}

	// Create a new contact
	success, msg := phonebook.AddContact("contacts/johndoe.json", contact)
	if !success {
		log.Printf("Error adding contact: %v\n", msg)
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
}
