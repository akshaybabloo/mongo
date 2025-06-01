// Package mongodb is a simple wrapper for MongoDb Driver, this package uses "_id" instead of "_id" to find or add a document.
//
// It is important to know that you will have to index id field for optimum performance.
//
// In general, you wouldn't need this package at all, if you rely more on "_id" and simple access to MongoDB API than this module will help you.
//
// Example:
//
//	import "github.com/akshaybabloo/mongodb"
//
//	type data struct {
//		ID   int    `bson:"_id"`
//		Name string `bson:"name"`
//	}
//
//	func main() {
//		client := NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")
//
//		testData := data{
//			ID:   1,
//			Name: "Akshay",
//		}
//
//		done, err := client.Add("test_collection", testData)
//		if err != nil {
//			panic(err)
//		}
//		print(done.InsertedID)
//	}
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Client takes in the
type Client struct {
	// ConnectionUrl which connects to MongoDB atlas or local deployment
	ConnectionUrl string

	// DatabaseName with database name
	DatabaseName string

	// Highly recommend using timeout Context
	Context context.Context
}

// NewMongoClient returns Client and it's associated functions
func NewMongoClient(connectionURL string, databaseName string, ctx context.Context) *Client {
	return &Client{
		ConnectionUrl: connectionURL,
		DatabaseName:  databaseName,
		Context:       ctx,
	}
}

// NewMongoClientDefault returns Client, and it's associated functions with default context
func NewMongoClientDefault(connectionURL string, databaseName string) *Client {
	return &Client{
		ConnectionUrl: connectionURL,
		DatabaseName:  databaseName,
		Context:       context.Background(),
	}
}

// Add can be used to add document to MongoDB
func (c *Client) Add(collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.InsertOne(c.Context, data)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// AddMany can be used to add multiple documents to MongoDB
func (c *Client) AddMany(collectionName string, data []interface{}) (*mongo.InsertManyResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.InsertMany(c.Context, data)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// Update can be used to update values by its ID
func (c *Client) Update(collectionName string, id string, data interface{}) (*mongo.UpdateResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	updateResult, err := collection.UpdateOne(c.Context, bson.M{"_id": id}, bson.D{{"$set", data}})
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

// UpdateCustom can be used to update values by a filter - bson.M{}, bson.A{}, or bson.D{}
func (c *Client) UpdateCustom(collectionName string, filter interface{}, data interface{}, updateOptions ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	updateResult, err := collection.UpdateOne(c.Context, filter, bson.D{{"$set", data}}, updateOptions...)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

// Delete deletes a document by ID only.
func (c *Client) Delete(collectionName string, id string) (*mongo.DeleteResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.DeleteOne(c.Context, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// DeleteCustom deletes a document by a filter - bson.M{}, bson.A{}, or bson.D{}
func (c *Client) DeleteCustom(collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.DeleteOne(c.Context, filter)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// DeleteMany deletes many documents - bson.M{}, bson.A{}, or bson.D{}
func (c *Client) DeleteMany(collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	insertResult, err := collection.DeleteMany(c.Context, filter)
	if err != nil {
		return nil, err
	}
	return insertResult, nil
}

// Get finds one document based on "_id"
func (c *Client) Get(collectionName string, id string) (*mongo.SingleResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	findOne := collection.FindOne(c.Context, bson.M{"_id": id})

	return findOne, nil
}

// GetCustom finds one document by a filter - bson.M{}, bson.A{}, or bson.D{}
func (c *Client) GetCustom(collectionName string, filter interface{}) (*mongo.SingleResult, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	findOne := collection.FindOne(c.Context, filter)

	return findOne, nil
}

// GetAll finds all documents by "_id".
//
// The 'result' parameter needs to be a pointer.
func (c *Client) GetAll(collectionName string, id string, result interface{}) error {
	client, err := c.client()
	if err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	find, err := collection.Find(c.Context, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if err = find.All(c.Context, result); err != nil {
		return err
	}

	return nil
}

// GetAllCustom finds all documents by filter - bson.M{}, bson.A{}, or bson.D{}.
//
// The 'result' parameter needs to be a pointer.
func (c *Client) GetAllCustom(collectionName string, filter interface{}, result interface{}) error {
	client, err := c.client()
	if err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	find, err := collection.Find(c.Context, filter)
	if err != nil {
		return err
	}

	if err = find.All(c.Context, result); err != nil {
		return err
	}

	return nil
}

// Collection returns mongo.Collection
//
// Note: Do not forget to do - defer Client.Disconnect(ctx)
func (c *Client) Collection(collectionName string) (*mongo.Collection, *mongo.Client, context.Context, error) {
	client, err := c.client()
	if err != nil {
		return nil, nil, nil, err
	}
	db := client.Database(c.DatabaseName)

	collection := db.Collection(collectionName)
	return collection, client, c.Context, nil
}

// DB returns mongo.Database
func (c *Client) DB() (*mongo.Database, error) {
	client, err := c.client()
	if err != nil {
		return nil, err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)

	return db, nil
}

// RawClient returns mongo.Client
func (c *Client) RawClient() (*mongo.Client, error) {
	return c.client()
}

func (c *Client) client() (*mongo.Client, error) {
	// c.Context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	client, err := mongo.Connect(options.Client().ApplyURI(c.ConnectionUrl))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) DeleteDatabase() error {
	client, err := c.client()
	if err != nil {
		return err
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(c.Context)
		if err != nil {
			return
		}
	}(client, c.Context)
	db := client.Database(c.DatabaseName)
	err = db.Drop(c.Context)
	if err != nil {
		return err
	}

	return nil

}
