package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDatabase struct {
	client     *mongo.Client
	collection *mongo.Collection
}

type MongoDocument struct {
	Location string                 `bson:"_id"`
	Data     map[string]interface{} `bson:"data"`
}

func NewMongoDatabase(uri, database, collection string) (*MongoDatabase, error) {
	ctx := context.Background()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	coll := client.Database(database).Collection(collection)

	// Create unique index on location
	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %v", err)
	}

	return &MongoDatabase{
		client:     client,
		collection: coll,
	}, nil
}

func (m *MongoDatabase) Create(location string, data map[string]interface{}) (bool, string) {
	ctx := context.Background()

	doc := MongoDocument{
		Location: location,
		Data:     data,
	}

	_, err := m.collection.InsertOne(ctx, doc)
	if mongo.IsDuplicateKeyError(err) {
		return false, "Location already exists"
	}
	if err != nil {
		return false, fmt.Sprintf("Error creating document: %v", err)
	}

	return true, ""
}

func (m *MongoDatabase) Read(location string) (bool, string, map[string]interface{}) {
	ctx := context.Background()

	var doc MongoDocument
	err := m.collection.FindOne(ctx, bson.M{"_id": location}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return false, "Location does not exist", nil
	}
	if err != nil {
		return false, fmt.Sprintf("Error reading document: %v", err), nil
	}

	return true, "", doc.Data
}

func (m *MongoDatabase) Update(location string, data map[string]interface{}) (bool, string) {
	ctx := context.Background()

	result, err := m.collection.UpdateOne(
		ctx,
		bson.M{"_id": location},
		bson.M{"$set": bson.M{"data": data}},
	)
	if err != nil {
		return false, fmt.Sprintf("Error updating document: %v", err)
	}
	if result.MatchedCount == 0 {
		return false, "Location does not exist"
	}

	return true, ""
}

func (m *MongoDatabase) Delete(location string) (bool, string) {
	ctx := context.Background()

	result, err := m.collection.DeleteOne(ctx, bson.M{"_id": location})
	if err != nil {
		return false, fmt.Sprintf("Error deleting document: %v", err)
	}
	if result.DeletedCount == 0 {
		return false, "Location does not exist"
	}

	return true, ""
}

func (m *MongoDatabase) Close() error {
	return m.client.Disconnect(context.Background())
}
