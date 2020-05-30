package mongo

import "testing"

var client NewMongoDbClient

type data struct {
	Id   int    `bson:"id"`
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
		Id:   1,
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)
}

func TestNewMongoDbClient_Get(t *testing.T) {
	testData := data{
		Id:   2,
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var decodeData data
	data := client.Get("test_collection", 2).Decode(&decodeData)
	t.Logf("%v", decodeData)
	if data != nil {
		t.Errorf("No data found.")
	}
}

func TestNewMongoDbClient_Update(t *testing.T) {
	testData := data{
		Id:   3,
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

	update, err := client.Update("test_collection", 3, data)
	if err != nil {
		t.Errorf("Unable to update data. %s", err)
	}
	t.Logf("The ID is %d", update.ModifiedCount)
}

func TestNewMongoDbClient_Delete(t *testing.T) {
	testData := data{
		Id:   4,
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	deleted, err := client.Delete("test_collection", 4)
	if err != nil {
		t.Errorf("Unable to delete data. %s", err)
	}
	t.Logf("Number deleted %d", deleted.DeletedCount)
}

func TestNewMongoDbClient_Collection(t *testing.T) {
	collection := client.Collection("test_collection")
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
