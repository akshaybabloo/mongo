package mongodb

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var client *Client

type data struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func init() {
	client = NewMongoClient(
		"mongodb://root:example@localhost:27017/?retryWrites=true&w=majority",
		"test",
		context.Background())
}

func TestClient_Add(t *testing.T) {
	testData := data{
		ID:   "1231",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)
}

func TestClient_AddMany(t *testing.T) {
	var testData = []interface{}{
		data{
			ID:   "111",
			Name: "Akshay",
		},
		data{
			ID:   "222",
			Name: "Raj",
		},
	}

	done, err := client.AddMany("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedIDs)
}

func TestClient_DeleteMany(t *testing.T) {
	type data struct {
		ID   string `bson:"_id"`
		Name string `bson:"name"`
	}

	var testData = []interface{}{
		data{
			ID:   "1",
			Name: "Akshay",
		},
		data{
			ID:   "2",
			Name: "Raj",
		},
	}

	done, err := client.AddMany("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedIDs)

	deleted, err := client.DeleteMany("test_collection", bson.M{"_id": bson.M{"$in": bson.A{"1", "2"}}})
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("Deleted items: %d", deleted.DeletedCount)
}

func TestClient_Get(t *testing.T) {
	testData := data{
		ID:   "2",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var decodeData data
	data, err := client.Get("test_collection", "2")
	if err != nil {
		panic("data not found")
	}
	err = data.Decode(&decodeData)
	if err != nil {
		t.Errorf("No data found.")
	}
	t.Logf("%v", decodeData)
}

func TestClient_GetCustom(t *testing.T) {
	// testData := data{
	// 	ID:   "2",
	// 	Name: "Akshay",
	// }
	//
	// done, err := client.Add("test_collection", testData)
	// if err != nil {
	// 	t.Errorf("Unable to add data. %s", err)
	// }
	// t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var decodeData data
	data, err := client.GetCustom("test_collection", bson.M{"_id": "2"})
	if err != nil {
		t.Errorf("No data found.")
		return
	}
	err = data.Decode(&decodeData)
	if err != nil {
		t.Errorf("No data found.")
	}
	t.Logf("%v", decodeData)

}

func TestClient_GetAll(t *testing.T) {
	testData := data{
		ID:   "123",
		Name: "Akshay",
	}

	_, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}

	// Actual test
	var result []data
	err = client.GetAll("test_collection", "1", &result)
	if err != nil {
		t.Errorf("No data found.")
	}
	t.Logf("%v", result)
}

func TestClient_GetAllCustom(t *testing.T) {

	// Actual test
	var result []data
	err := client.GetAllCustom("test_collection", bson.M{"_id": "1"}, &result)
	if err != nil {
		t.Errorf("No data found.")
	}
	t.Logf("%v", result)
}

// func TestClient_Update(t *testing.T) {
// 	testData := data{
// 		ID:   "3",
// 		Name: "Akshay",
// 	}
//
// 	done, err := client.Add("test_collection", testData)
// 	if err != nil {
// 		t.Errorf("Unable to add data. %s", err)
// 	}
// 	t.Logf("The ID is %s", done.InsertedID)
//
// 	// Actual test
// 	data := data{
// 		Name: "Gollahalli",
// 	}
//
// 	update, err := client.Update("test_collection", "3", data)
// 	if err != nil {
// 		t.Errorf("Unable to update data. %s", err)
// 	}
// 	t.Logf("The ID is %d", update.ModifiedCount)
// }

func TestClient_Delete(t *testing.T) {
	testData := data{
		ID:   "4",
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

func TestClient_DeleteCustom(t *testing.T) {
	testData := data{
		ID:   "4",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	deleted, err := client.DeleteCustom("test_collection", bson.M{"_id": 4})
	if err != nil {
		t.Errorf("Unable to delete data. %s", err)
	}
	t.Logf("Number deleted %d", deleted.DeletedCount)
}

func TestClient_Collection(t *testing.T) {
	collection, client, ctx, err := client.Collection("test_collection")
	if err != nil {
		t.Errorf("something went wrong. %s", err)
	}
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			return
		}
	}(client, ctx)
	if collection.Name() != "test_collection" {
		t.Errorf("Collection name incorrect")
	}
}

func TestClient_DB(t *testing.T) {
	db, _ := client.DB()
	if db.Name() != "test" {
		t.Errorf("Database name incorrect")
	}
}

func TestClient_DeleteDatabase(t *testing.T) {
	err := client.DeleteDatabase()
	if err != nil {
		t.Errorf("Unable to delete database. %s", err)
	}
	t.Logf("Database deleted successfully")
}
