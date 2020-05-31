/*
mongo is a simple wrapper for MongoDb Driver, this package uses "id" instead of "_id" to find or add a document.

It is important to know that you will have to index id field for optimum performance.

In general you would't need this package at all, if you rely more on "id" and simple access to MongoDB API then this module will help you.

Example:

	import "github.com/akshaybabloo/mongo"

	type data struct {
		Id   int    `bson:"id"`
		Name string `bson:"name"`
	}

	func main() {
		client := mongo.NewMongoDbClient{
			ConnectionUrl: "mongodb://localhost:27017/?retryWrites=true&w=majority",
			DatabaseName:  "test",
		}

		testData := data{
			Id:   1,
			Name: "Akshay",
		}

		done, err := client.Add("test_collection", testData)
		if err != nil {
			panic(err)
		}
		print(done.InsertedID)
	}

*/
package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDb implements MongoDb's CRUD operations
type MongoDb interface {
	Add(collectionName string, data interface{}) (*mongo.InsertOneResult, error)
	Update(collectionName string, id int, data interface{}) (*mongo.UpdateResult, error)
	Delete(collectionName string, id int) (*mongo.DeleteResult, error)
	Get(collectionName string, id int) *mongo.SingleResult
	Collection(collectionName string) *mongo.Collection
	DB() *mongo.Database
	client() (*mongo.Client, context.Context)
}

// NewMongoDbClient takes in the
type NewMongoDbClient struct {
	// ConnectionUrl which connects to MongoDB atlas or local deployment
	ConnectionUrl string

	// DatabaseName with database name
	DatabaseName string
}

// Add can be used to add document to MongoDB
func (connectionDetails NewMongoDbClient) Add(collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// Update can be used to update values by it's ID
func (connectionDetails NewMongoDbClient) Update(collectionName string, id int, data interface{}) (*mongo.UpdateResult, error) {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	updateResult, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.D{{"$set", data}})
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

// Delete deletes a document by ID only.
func (connectionDetails NewMongoDbClient) Delete(collectionName string, id int) (*mongo.DeleteResult, error) {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// Get finds one document based on "id" and not "_id"
func (connectionDetails NewMongoDbClient) Get(collectionName string, id int) *mongo.SingleResult {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	findOne := collection.FindOne(ctx, bson.M{"id": id})

	return findOne
}

// Collection returns mongo.Collection
func (connectionDetails NewMongoDbClient) Collection(collectionName string) *mongo.Collection {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	return collection
}

// DB returns mongo.Database
func (connectionDetails NewMongoDbClient) DB() *mongo.Database {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	return db
}

func (connectionDetails NewMongoDbClient) client() (*mongo.Client, context.Context) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionDetails.ConnectionUrl))
	if err != nil {
		panic(err)
	}

	return client, ctx
}
