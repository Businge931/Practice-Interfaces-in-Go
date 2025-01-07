package database

// import (
// 	"context"
// 	"fmt"

// 	"cloud.google.com/go/firestore"
// 	"google.golang.org/api/option"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// type FirestoreDatabase struct {
// 	client     *firestore.Client
// 	collection string
// }

// // NewFirestoreDatabase creates a new Firestore database client
// // credentialsFile is the path to the service account JSON file
// func NewFirestoreDatabase(projectID, collection, credentialsFile string) (*FirestoreDatabase, error) {
// 	ctx := context.Background()
	
// 	// Create Firestore client
// 	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsFile))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create Firestore client: %v", err)
// 	}

// 	return &FirestoreDatabase{
// 		client:     client,
// 		collection: collection,
// 	}, nil
// }

// func (f *FirestoreDatabase) Create(location string, data map[string]interface{}) (bool, string) {
// 	ctx := context.Background()

// 	// Check if document already exists
// 	_, err := f.client.Collection(f.collection).Doc(location).Get(ctx)
// 	if err == nil {
// 		return false, "Location already exists"
// 	}
// 	if status.Code(err) != codes.NotFound {
// 		return false, fmt.Sprintf("Error checking document existence: %v", err)
// 	}

// 	// Create the document
// 	_, err = f.client.Collection(f.collection).Doc(location).Set(ctx, data)
// 	if err != nil {
// 		return false, fmt.Sprintf("Error creating document: %v", err)
// 	}

// 	return true, ""
// }

// func (f *FirestoreDatabase) Read(location string) (bool, string, map[string]interface{}) {
// 	ctx := context.Background()

// 	// Get the document
// 	doc, err := f.client.Collection(f.collection).Doc(location).Get(ctx)
// 	if status.Code(err) == codes.NotFound {
// 		return false, "Location does not exist", nil
// 	}
// 	if err != nil {
// 		return false, fmt.Sprintf("Error reading document: %v", err), nil
// 	}

// 	// Convert the data
// 	data := make(map[string]interface{})
// 	for k, v := range doc.Data() {
// 		data[k] = v
// 	}

// 	return true, "", data
// }

// func (f *FirestoreDatabase) Update(location string, data map[string]interface{}) (bool, string) {
// 	ctx := context.Background()

// 	// Check if document exists
// 	_, err := f.client.Collection(f.collection).Doc(location).Get(ctx)
// 	if status.Code(err) == codes.NotFound {
// 		return false, "Location does not exist"
// 	}
// 	if err != nil {
// 		return false, fmt.Sprintf("Error checking document existence: %v", err)
// 	}

// 	// Update the document
// 	_, err = f.client.Collection(f.collection).Doc(location).Set(ctx, data)
// 	if err != nil {
// 		return false, fmt.Sprintf("Error updating document: %v", err)
// 	}

// 	return true, ""
// }

// func (f *FirestoreDatabase) Delete(location string) (bool, string) {
// 	ctx := context.Background()

// 	// Check if document exists
// 	_, err := f.client.Collection(f.collection).Doc(location).Get(ctx)
// 	if status.Code(err) == codes.NotFound {
// 		return false, "Location does not exist"
// 	}
// 	if err != nil {
// 		return false, fmt.Sprintf("Error checking document existence: %v", err)
// 	}

// 	// Delete the document
// 	_, err = f.client.Collection(f.collection).Doc(location).Delete(ctx)
// 	if err != nil {
// 		return false, fmt.Sprintf("Error deleting document: %v", err)
// 	}

// 	return true, ""
// }

// func (f *FirestoreDatabase) Close() error {
// 	return f.client.Close()
// }
