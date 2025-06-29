// Package mongodb is an enhanced wrapper for MongoDB Driver with better connection management,
// error handling, and performance optimizations.
//
// This package uses custom "id" field instead of MongoDB's default "_id" to find or add a document.
//
// It is important to know that you will have to index the id field for optimum performance.
//
// Example:
//
//	import "github.com/akshaybabloo/mongodb/v6"
//
//	type data struct {
//		ID   int    `bson:"_id"`
//		Name string `bson:"name"`
//	}
//
//	func main() {
//		client, err := mongodb.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")
//		if err != nil {
//			panic(err)
//		}
//		defer client.Close()
//
//		testData := data{
//			ID:   1,
//			Name: "Akshay",
//		}
//
//		ctx := context.Background()
//		result, err := client.Add(ctx, "test_collection", testData)
//		if err != nil {
//			panic(err)
//		}
//		fmt.Println(result.InsertedID)
//	}
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Client wraps MongoDB client with simplified operations and improved connection management
type Client struct {
	// ConnectionUrl which connects to MongoDB atlas or local deployment
	ConnectionUrl string
	// DatabaseName with database name
	DatabaseName string
	// client holds the MongoDB client instance
	client *mongo.Client
	// mutex for thread-safe operations
	mutex sync.RWMutex
	// connected tracks connection state
	connected bool
}

// NewMongoClient creates a new MongoDB client and establishes connection
func NewMongoClient(connectionURL string, databaseName string) (*Client, error) {
	if connectionURL == "" {
		return nil, errors.New("connection URL cannot be empty")
	}
	if databaseName == "" {
		return nil, errors.New("database name cannot be empty")
	}

	c := &Client{
		ConnectionUrl: connectionURL,
		DatabaseName:  databaseName,
	}

	if err := c.connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return c, nil
}

// connect establishes connection to MongoDB
func (c *Client) connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.connected {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(c.ConnectionUrl))
	if err != nil {
		return err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			return err
		}
		return err
	}

	c.client = client
	c.connected = true
	return nil
}

// Close disconnects from MongoDB
func (c *Client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.connected || c.client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.client.Disconnect(ctx)
	c.connected = false
	c.client = nil
	return err
}

// getClient returns the MongoDB client, ensuring connection
func (c *Client) getClient() (*mongo.Client, error) {
	c.mutex.RLock()
	if c.connected && c.client != nil {
		defer c.mutex.RUnlock()
		return c.client, nil
	}
	c.mutex.RUnlock()

	// Need to reconnect
	if err := c.connect(); err != nil {
		return nil, err
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.client, nil
}

// validateParams validates common parameters
func (c *Client) validateParams(collectionName string) error {
	if collectionName == "" {
		return errors.New("collection name cannot be empty")
	}
	return nil
}

// getCollection returns a MongoDB collection
func (c *Client) getCollection(collectionName string) (*mongo.Collection, error) {
	if err := c.validateParams(collectionName); err != nil {
		return nil, err
	}

	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	db := client.Database(c.DatabaseName)
	return db.Collection(collectionName), nil
}

// Add inserts a single document to MongoDB
func (c *Client) Add(ctx context.Context, collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.InsertOne(ctx, data)
}

// AddMany inserts multiple documents to MongoDB
func (c *Client) AddMany(ctx context.Context, collectionName string, data []interface{}) (*mongo.InsertManyResult, error) {
	if len(data) == 0 {
		return nil, errors.New("data slice cannot be empty")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.InsertMany(ctx, data)
}

// Update updates a document by its ID
func (c *Client) Update(ctx context.Context, collectionName string, id string, data interface{}) (*mongo.UpdateResult, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{"$set", data}})
}

// UpdateCustom updates a document using a custom filter
func (c *Client) UpdateCustom(ctx context.Context, collectionName string, filter interface{}, data interface{}, updateOptions ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.UpdateOne(ctx, filter, bson.D{{"$set", data}}, updateOptions...)
}

// UpdateMany updates multiple documents using a filter
func (c *Client) UpdateMany(ctx context.Context, collectionName string, filter interface{}, data interface{}, updateOptions ...options.Lister[options.UpdateManyOptions]) (*mongo.UpdateResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.UpdateMany(ctx, filter, bson.D{{"$set", data}}, updateOptions...)
}

// Delete deletes a document by ID
func (c *Client) Delete(ctx context.Context, collectionName string, id string) (*mongo.DeleteResult, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.DeleteOne(ctx, bson.M{"_id": id})
}

// DeleteCustom deletes a document using a custom filter
func (c *Client) DeleteCustom(ctx context.Context, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.DeleteOne(ctx, filter)
}

// DeleteMany deletes multiple documents using a filter
func (c *Client) DeleteMany(ctx context.Context, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	return collection.DeleteMany(ctx, filter)
}

// Get finds one document by ID
func (c *Client) Get(ctx context.Context, collectionName string, id string) (*mongo.SingleResult, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	result := collection.FindOne(ctx, bson.M{"_id": id})
	return result, nil
}

// GetCustom finds one document using a custom filter
func (c *Client) GetCustom(ctx context.Context, collectionName string, filter interface{}, findOptions ...options.Lister[options.FindOneOptions]) (*mongo.SingleResult, error) {
	if filter == nil {
		return nil, errors.New("filter cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return nil, err
	}

	result := collection.FindOne(ctx, filter, findOptions...)
	return result, nil
}

// FindByID finds all documents with the same ID (renamed from GetAll for clarity)
func (c *Client) FindByID(ctx context.Context, collectionName string, id string, result interface{}) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	if result == nil {
		return errors.New("result cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return err
	}

	cursor, err := collection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

// FindAll finds all documents using a custom filter
func (c *Client) FindAll(ctx context.Context, collectionName string, filter interface{}, result interface{}, findOptions ...options.Lister[options.FindOptions]) error {
	if filter == nil {
		return errors.New("filter cannot be nil")
	}
	if result == nil {
		return errors.New("result cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return err
	}

	cursor, err := collection.Find(ctx, filter, findOptions...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

// Exists checks if a document exists by ID
func (c *Client) Exists(ctx context.Context, collectionName string, id string) (bool, error) {
	if id == "" {
		return false, errors.New("id cannot be empty")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return false, err
	}

	count, err := collection.CountDocuments(ctx, bson.M{"_id": id})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// ExistsCustom checks if a document exists using a custom filter
func (c *Client) ExistsCustom(ctx context.Context, collectionName string, filter interface{}) (bool, error) {
	if filter == nil {
		return false, errors.New("filter cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return false, err
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Aggregate performs an aggregation operation on a collection
func (c *Client) Aggregate(ctx context.Context, collectionName string, pipeline interface{}, result interface{}, aggregateOptions ...options.Lister[options.AggregateOptions]) error {
	if pipeline == nil {
		return errors.New("pipeline cannot be nil")
	}
	if result == nil {
		return errors.New("result cannot be nil")
	}

	collection, err := c.getCollection(collectionName)
	if err != nil {
		return err
	}

	cursor, err := collection.Aggregate(ctx, pipeline, aggregateOptions...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	return cursor.All(ctx, result)
}

// Collection returns a MongoDB collection
// Note: The client connection is managed internally, no need to manually disconnect
func (c *Client) Collection(collectionName string) (*mongo.Collection, error) {
	return c.getCollection(collectionName)
}

// Database returns the MongoDB database
func (c *Client) Database() (*mongo.Database, error) {
	client, err := c.getClient()
	if err != nil {
		return nil, err
	}

	return client.Database(c.DatabaseName), nil
}

// RawClient returns the underlying MongoDB client
// Note: Do not disconnect this client manually, use Close() method instead
func (c *Client) RawClient() (*mongo.Client, error) {
	return c.getClient()
}

// DropDatabase drops the entire database
func (c *Client) DropDatabase(ctx context.Context) error {
	db, err := c.Database()
	if err != nil {
		return err
	}

	return db.Drop(ctx)
}

// Ping tests the connection to MongoDB
func (c *Client) Ping(ctx context.Context) error {
	client, err := c.getClient()
	if err != nil {
		return err
	}

	return client.Ping(ctx, nil)
}

// IsConnected returns the connection status
func (c *Client) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.connected
}
