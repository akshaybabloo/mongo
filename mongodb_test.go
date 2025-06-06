package mongodb

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/v2/bson"
)

var client *Client

type data struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func init() {
	var err error
	client, err = NewMongoClient(
		"mongodb://root:example@localhost:27017/?retryWrites=true&w=majority",
		"test")
	if err != nil {
		panic(err)
	}
}

func TestClient_Add(t *testing.T) {
	testData := data{
		ID:   "1231",
		Name: "Akshay",
	}
	ctx := context.Background()
	done, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
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

	ctx := context.Background()
	done, err := client.AddMany(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedIDs)
}

func TestClient_Exists(t *testing.T) {
	testData := data{
		ID:   "1233455",
		Name: "Akshay",
	}
	ctx := context.Background()
	done, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	exists, err := client.Exists(ctx, "test_collection", "1233455")
	if err != nil {
		client.Close()
		t.Errorf("Unable to check existence. %s", err)
	}
	t.Logf("Exists: %v", exists)

	existsCustom, err := client.ExistsCustom(ctx, "test_collection", bson.M{"_id": "1233455"})
	if err != nil {
		client.Close()
		t.Errorf("Unable to check existence with custom query. %s", err)
	}
	t.Logf("Exists with custom query: %v", existsCustom)
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
	ctx := context.Background()
	done, err := client.AddMany(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedIDs)

	ctx = context.Background()
	deleted, err := client.DeleteMany(ctx, "test_collection", bson.M{"_id": bson.M{"$in": bson.A{"1", "2"}}})
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("Deleted items: %d", deleted.DeletedCount)
}

func TestClient_Get(t *testing.T) {
	testData := data{
		ID:   "2",
		Name: "Akshay",
	}

	ctx := context.Background()
	done, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	var decodeData data
	ctx = context.Background()
	data, err := client.Get(ctx, "test_collection", "2")
	if err != nil {
		client.Close()
		t.Errorf("No data found.")
	}
	err = data.Decode(&decodeData)
	if err != nil {
		client.Close()
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
	ctx := context.Background()
	data, err := client.GetCustom(ctx, "test_collection", bson.M{"_id": "2"})
	if err != nil {
		client.Close()
		t.Errorf("No data found.")
		return
	}
	err = data.Decode(&decodeData)
	if err != nil {
		client.Close()
		t.Errorf("No data found.")
	}
	t.Logf("%v", decodeData)

}

func TestClient_GetAll(t *testing.T) {
	testData := data{
		ID:   "123",
		Name: "Akshay",
	}
	ctx := context.Background()
	_, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}

	// Actual test
	var result []data
	ctx = context.Background()
	err = client.FindAll(ctx, "test_collection", bson.M{"_id": "1"}, &result)
	if err != nil {
		client.Close()
		t.Errorf("No data found.")
	}
	t.Logf("%v", result)
}

// func TestClient_GetAllCustom(t *testing.T) {
//
// 	// Actual test
// 	var result []data
// 	err := client.GetAllCustom("test_collection", bson.M{"_id": "1"}, &result)
// 	if err != nil {
// 		client.Close()
// 		t.Errorf("No data found.")
// 	}
// 	t.Logf("%v", result)
// }

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
	ctx := context.Background()
	done, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	ctx = context.Background()
	deleted, err := client.Delete(ctx, "test_collection", "4")
	if err != nil {
		client.Close()
		t.Errorf("Unable to delete data. %s", err)
	}
	t.Logf("Number deleted %d", deleted.DeletedCount)
}

func TestClient_DeleteCustom(t *testing.T) {
	testData := data{
		ID:   "4",
		Name: "Akshay",
	}
	ctx := context.Background()
	done, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		client.Close()
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %s", done.InsertedID)

	// Actual test
	ctx = context.Background()
	deleted, err := client.DeleteCustom(ctx, "test_collection", bson.M{"_id": 4})
	if err != nil {
		client.Close()
		t.Errorf("Unable to delete data. %s", err)
	}
	t.Logf("Number deleted %d", deleted.DeletedCount)
}

func TestClient_Collection(t *testing.T) {
	collection, err := client.Collection("test_collection")
	if err != nil {
		client.Close()
		t.Errorf("something went wrong. %s", err)
	}
	if collection.Name() != "test_collection" {
		client.Close()
		t.Errorf("Collection name incorrect")
	}
}

func TestClient_DB(t *testing.T) {
	db, _ := client.Database()
	if db.Name() != "test" {
		client.Close()
		t.Errorf("Database name incorrect")
	}
}

func TestClient_DeleteDatabase(t *testing.T) {
	ctx := context.Background()
	err := client.DropDatabase(ctx)
	if err != nil {
		client.Close()
		t.Errorf("Unable to delete database. %s", err)
	}
	t.Logf("Database deleted successfully")
	err = client.Close()
	if err != nil {
		t.Errorf("Unable to close client. %s", err)
	} else {
		t.Logf("Client closed successfully")
	}
}
