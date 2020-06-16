package mongo

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

var client NewMongoDbClient

type data struct {
	Id   string `bson:"id"`
	Name string `bson:"name"`
}

func init() {
	client = NewMongoDbClient{
		ConnectionUrl: "mongodb://localhost:27017/?retryWrites=true&w=majority",
		DatabaseName:  "test",
	}
}

func TestNewMongoDbClient_Add(t *testing.T) {
	testData := data{
		Id:   "1",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)
}

func TestNewMongoDbClient_AddMany(t *testing.T) {
	var testData = []interface{}{
		data{
			Id:   "111",
			Name: "Akshay",
		},
		data{
			Id:   "222",
			Name: "Raj",
		},
	}

	done, err := client.AddMany("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedIDs)
}

func TestNewMongoDbClient_Get(t *testing.T) {
	testData := data{
		Id:   "2",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var decodeData data
	data := client.Get("test_collection", "2").Decode(&decodeData)
	t.Logf("%v", decodeData)
	if data != nil {
		t.Errorf("No data found.")
	}
}

func TestNewMongoDbClient_GetCustom(t *testing.T) {
	testData := data{
		Id:   "2",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var decodeData data
	data := client.GetCustom("test_collection", bson.M{"id": "2"}).Decode(&decodeData)
	t.Logf("%v", decodeData)
	if data != nil {
		t.Errorf("No data found.")
	}
}

func TestNewMongoDbClient_GetAll(t *testing.T) {
	testData := data{
		Id:   "123",
		Name: "Akshay",
	}

	_, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	//t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var result []data
	err = client.GetAll("test_collection", "1", &result)
	if err != nil {
		t.Errorf("No data found.")
	}
	t.Logf("%v", result)
}

func TestNewMongoDbClient_GetAllCustom(t *testing.T) {
	testData := data{
		Id:   "123",
		Name: "Akshay",
	}

	_, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	//t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var result []data
	err = client.GetAllCustom("test_collection", bson.M{"id": "1"}, &result)
	if err != nil {
		t.Errorf("No data found.")
	}
	t.Logf("%v", result)
}

func TestNewMongoDbClient_Update(t *testing.T) {
	testData := data{
		Id:   "3",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	data := data{
		Name: "Gollahalli",
	}

	update, err := client.Update("test_collection", "3", data)
	if err != nil {
		t.Errorf("Unable to update data. %s", err)
	}
	t.Logf("The ID is %d", update.ModifiedCount)
}

func TestNewMongoDbClient_Delete(t *testing.T) {
	testData := data{
		Id:   "4",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	deleted, err := client.Delete("test_collection", "4")
	if err != nil {
		t.Errorf("Unable to delete data. %s", err)
	}
	t.Logf("Number deleted %d", deleted.DeletedCount)
}

func TestNewMongoDbClient_Collection(t *testing.T) {
	collection, client, ctx := client.Collection("test_collection")
	defer client.Disconnect(ctx)
	if collection.Name() != "test_collection" {
		t.Errorf("Collection name incorrect")
	}
}

func TestNewMongoDbClient_DB(t *testing.T) {
	db := client.DB()
	if db.Name() != "test" {
		t.Errorf("Database name incorrect")
	}
}
