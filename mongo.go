// Package mongo is a simple wrapper for MongoDb Driver, this package uses "_id" instead of "_id" to find or add a document.
//
// It is important to know that you will have to index id field for optimum performance.
//
// In general, you wouldn't need this package at all, if you rely more on "_id" and simple access to MongoDB API than this module will help you.
//
// Example:
//
// 	import "github.com/akshaybabloo/mongo"
//
// 	type data struct {
// 		ID   int    `bson:"_id"`
// 		Name string `bson:"name"`
// 	}
//
// 	func main() {
// 		client := NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")
//
// 		testData := data{
// 			ID:   1,
// 			Name: "Akshay",
// 		}
//
// 		done, err := client.Add("test_collection", testData)
// 		if err != nil {
// 			panic(err)
// 		}
// 		print(done.InsertedID)
// 	}
//
//
package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client takes in the
type Client struct {
	// ConnectionUrl which connects to MongoDB atlas or local deployment
	ConnectionUrl string

	// DatabaseName with database name
	DatabaseName string
}

// NewMongoClient returns Client and it's associated functions
func NewMongoClient(connectionURL string, databaseName string) *Client {
	return &Client{
		ConnectionUrl: connectionURL,
		DatabaseName:  databaseName,
	}
}

// Add can be used to add document to MongoDB
func (connectionDetails *Client) Add(collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// AddMany can be used to add multiple documents to MongoDB
func (connectionDetails *Client) AddMany(collectionName string, data []interface{}) (*mongo.InsertManyResult, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// Update can be used to update values by its ID
func (connectionDetails *Client) Update(collectionName string, id string, data interface{}) (*mongo.UpdateResult, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	updateResult, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{"$set", data}})
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

// Delete deletes a document by ID only.
func (connectionDetails *Client) Delete(collectionName string, id string) (*mongo.DeleteResult, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// Get finds one document based on "_id"
func (connectionDetails *Client) Get(collectionName string, id string) (*mongo.SingleResult, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	findOne := collection.FindOne(ctx, bson.M{"_id": id})

	return findOne, nil
}

// GetCustom finds one document by a filter - bson.M{}, bson.A{}, or bson.D{}
func (connectionDetails *Client) GetCustom(collectionName string, filter interface{}) (*mongo.SingleResult, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	findOne := collection.FindOne(ctx, filter)

	return findOne, nil
}

// GetAll finds all documents by "_id".
//
// The 'result' parameter needs to be a pointer.
func (connectionDetails *Client) GetAll(collectionName string, id string, result interface{}) error {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	find, err := collection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if err = find.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// GetAllCustom finds all documents by filter - bson.M{}, bson.A{}, or bson.D{}.
//
// The 'result' parameter needs to be a pointer.
func (connectionDetails *Client) GetAllCustom(collectionName string, filter interface{}, result interface{}) error {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	find, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}

	if err = find.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// Collection returns mongo.Collection
//
// Note: Do not forget to do - defer Client.Disconnect(ctx)
func (connectionDetails *Client) Collection(collectionName string) (*mongo.Collection, *mongo.Client, context.Context, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, nil, nil, err
	}
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	return collection, client, ctx, nil
}

// DB returns mongo.Database
func (connectionDetails *Client) DB() (*mongo.Database, error) {
	client, ctx, err := connectionDetails.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	db := client.Database(connectionDetails.DatabaseName)

	return db, nil
}

// RawClient returns mongo.Client
func (connectionDetails *Client) RawClient() (*mongo.Client, context.Context, error) {
	return connectionDetails.client()
}

func (connectionDetails *Client) client() (*mongo.Client, context.Context, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionDetails.ConnectionUrl))
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
