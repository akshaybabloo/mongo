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
	AddMany(collectionName string, data []interface{}) (*mongo.InsertManyResult, error)
	Update(collectionName string, id string, data interface{}) (*mongo.UpdateResult, error)
	Delete(collectionName string, id string) (*mongo.DeleteResult, error)
	Get(collectionName string, id string) *mongo.SingleResult
	GetCustom(collectionName string, id string) *mongo.SingleResult
	GetAll(collectionName string, id string, result interface{}) error
	GetAllCustom(collectionName string, id string, result interface{}) error
	Collection(collectionName string) (*mongo.Collection, *mongo.Client, context.Context)
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

// AddMany can be used to add multiple documents to MongoDB
func (connectionDetails NewMongoDbClient) AddMany(collectionName string, data []interface{}) (*mongo.InsertManyResult, error) {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.InsertMany(ctx, data)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// Update can be used to update values by it's ID
func (connectionDetails NewMongoDbClient) Update(collectionName string, id string, data interface{}) (*mongo.UpdateResult, error) {
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
func (connectionDetails NewMongoDbClient) Delete(collectionName string, id string) (*mongo.DeleteResult, error) {
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
func (connectionDetails NewMongoDbClient) Get(collectionName string, id string) *mongo.SingleResult {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	findOne := collection.FindOne(ctx, bson.M{"id": id})

	return findOne
}

// Get finds one document by a filter - bson.M{}, bson.A{}, or bson.D{}
func (connectionDetails NewMongoDbClient) GetCustom(collectionName string, filter interface{}) *mongo.SingleResult {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	findOne := collection.FindOne(ctx, filter)

	return findOne
}

// GetAll finds all documents by "id" and not "_id".
//
//The 'result' parameter needs to be a pointer.
func (connectionDetails NewMongoDbClient) GetAll(collectionName string, id string, result interface{}) error {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	find, err := collection.Find(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if err = find.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// GetAll finds all documents by filter - bson.M{}, bson.A{}, or bson.D{}.
//
//The 'result' parameter needs to be a pointer.
func (connectionDetails NewMongoDbClient) GetAllCustom(collectionName string, filter interface{}, result interface{}) error {
	client, ctx := connectionDetails.client()
	defer client.Disconnect(ctx)
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
// Note: Do not forget to do - defer client.Disconnect(ctx)
func (connectionDetails NewMongoDbClient) Collection(collectionName string) (*mongo.Collection, *mongo.Client, context.Context) {
	client, ctx := connectionDetails.client()
	db := client.Database(connectionDetails.DatabaseName)

	collection := db.Collection(collectionName)
	return collection, client, ctx
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
