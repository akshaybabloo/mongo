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

func TestClient_Update(t *testing.T) {
	testData := data{
		ID:   "update-id-1",
		Name: "OriginalName",
	}
	ctx := context.Background()
	_, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		t.Fatalf("Unable to add data: %s", err)
	}

	// Update only the Name field (not the whole struct)
	updateData := bson.M{"name": "UpdatedName"}
	updateResult, err := client.Update(ctx, "test_collection", "update-id-1", updateData)
	if err != nil {
		t.Fatalf("Unable to update data: %s", err)
	}
	if updateResult.ModifiedCount != 1 {
		t.Errorf("Expected 1 document to be updated, got %d", updateResult.ModifiedCount)
	}

	// Verify update
	var got data
	res, err := client.Get(ctx, "test_collection", "update-id-1")
	if err != nil {
		t.Fatalf("Unable to get data: %s", err)
	}
	err = res.Decode(&got)
	if err != nil {
		t.Fatalf("Unable to decode data: %s", err)
	}
	if got.Name != "UpdatedName" {
		t.Errorf("Expected Name to be 'UpdatedName', got '%s'", got.Name)
	}
}

func TestClient_UpdateCustom(t *testing.T) {
	testData := data{
		ID:   "updatecustom-id-1",
		Name: "OriginalName",
	}
	ctx := context.Background()
	_, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		t.Fatalf("Unable to add data: %s", err)
	}

	updateData := bson.M{"name": "CustomUpdatedName"}
	updateResult, err := client.UpdateCustom(ctx, "test_collection", bson.M{"_id": "updatecustom-id-1"}, updateData)
	if err != nil {
		t.Fatalf("Unable to update data: %s", err)
	}
	if updateResult.ModifiedCount != 1 {
		t.Errorf("Expected 1 document to be updated, got %d", updateResult.ModifiedCount)
	}

	// Verify update
	var got data
	res, err := client.Get(ctx, "test_collection", "updatecustom-id-1")
	if err != nil {
		t.Fatalf("Unable to get data: %s", err)
	}
	err = res.Decode(&got)
	if err != nil {
		t.Fatalf("Unable to decode data: %s", err)
	}
	if got.Name != "CustomUpdatedName" {
		t.Errorf("Expected Name to be 'CustomUpdatedName', got '%s'", got.Name)
	}
}

func TestClient_UpdateMany(t *testing.T) {
	testData := []interface{}{
		data{ID: "updatemany-id-1", Name: "Name1"},
		data{ID: "updatemany-id-2", Name: "Name2"},
	}
	ctx := context.Background()
	_, err := client.AddMany(ctx, "test_collection", testData)
	if err != nil {
		t.Fatalf("Unable to add data: %s", err)
	}

	updateData := bson.M{"name": "BulkUpdated"}
	updateResult, err := client.UpdateMany(ctx, "test_collection", bson.M{"_id": bson.M{"$in": []string{"updatemany-id-1", "updatemany-id-2"}}}, updateData)
	if err != nil {
		t.Fatalf("Unable to update many: %s", err)
	}
	if updateResult.ModifiedCount != 2 {
		t.Errorf("Expected 2 documents to be updated, got %d", updateResult.ModifiedCount)
	}

	// Verify updates
	var results []data
	err = client.FindAll(ctx, "test_collection", bson.M{"_id": bson.M{"$in": []string{"updatemany-id-1", "updatemany-id-2"}}}, &results)
	if err != nil {
		t.Fatalf("Unable to find updated documents: %s", err)
	}
	for _, d := range results {
		if d.Name != "BulkUpdated" {
			t.Errorf("Expected Name to be 'BulkUpdated', got '%s' for ID '%s'", d.Name, d.ID)
		}
	}
}

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

func TestClient_FindByID(t *testing.T) {
	// Add a single test document
	testData := data{ID: "findbyid-test-unique", Name: "TestData"}
	ctx := context.Background()
	_, err := client.Add(ctx, "test_collection", testData)
	if err != nil {
		t.Fatalf("Unable to add test data: %s", err)
	}

	// Test FindByID - should find the single document
	var results []data
	err = client.FindByID(ctx, "test_collection", "findbyid-test-unique", &results)
	if err != nil {
		t.Fatalf("Unable to find by ID: %s", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 document, got %d", len(results))
	}

	if len(results) > 0 {
		if results[0].ID != "findbyid-test-unique" {
			t.Errorf("Expected ID 'findbyid-test-unique', got '%s'", results[0].ID)
		}
		if results[0].Name != "TestData" {
			t.Errorf("Expected Name 'TestData', got '%s'", results[0].Name)
		}
	}

	// Test with non-existent ID
	var emptyResults []data
	err = client.FindByID(ctx, "test_collection", "non-existent-id", &emptyResults)
	if err != nil {
		t.Fatalf("FindByID should not error for non-existent ID: %s", err)
	}
	if len(emptyResults) != 0 {
		t.Errorf("Expected 0 documents for non-existent ID, got %d", len(emptyResults))
	}
}

func TestClient_Ping(t *testing.T) {
	ctx := context.Background()
	err := client.Ping(ctx)
	if err != nil {
		t.Fatalf("Ping failed: %s", err)
	}
	t.Logf("Ping successful")
}

func TestClient_IsConnected(t *testing.T) {
	// Test that client is connected
	connected := client.IsConnected()
	if !connected {
		t.Errorf("Expected client to be connected, but it's not")
	}
	t.Logf("Client connection status: %v", connected)

	// Test after closing and reconnecting
	err := client.Close()
	if err != nil {
		t.Fatalf("Unable to close client: %s", err)
	}

	connected = client.IsConnected()
	if connected {
		t.Errorf("Expected client to be disconnected after Close(), but it's still connected")
	}

	// Reconnect by performing an operation (this will trigger reconnection)
	ctx := context.Background()
	err = client.Ping(ctx)
	if err != nil {
		t.Fatalf("Unable to reconnect: %s", err)
	}

	connected = client.IsConnected()
	if !connected {
		t.Errorf("Expected client to be connected after Ping(), but it's not")
	}
	t.Logf("Client reconnected successfully")
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

func TestClient_Aggregate(t *testing.T) {
	// Insert sample data
	testData := []interface{}{
		data{ID: "agg-1", Name: "Alice"},
		data{ID: "agg-2", Name: "Bob"},
		data{ID: "agg-3", Name: "Alice"},
	}
	ctx := context.Background()
	_, err := client.AddMany(ctx, "test_collection", testData)
	if err != nil {
		t.Fatalf("Unable to add data for aggregation: %s", err)
	}

	// Aggregation pipeline: group by name and count
	pipeline := bson.A{
		bson.M{"$group": bson.M{"_id": "$name", "count": bson.M{"$sum": 1}}},
	}

	type aggResult struct {
		ID    string `bson:"_id"`
		Count int32  `bson:"count"`
	}
	var results []aggResult

	err = client.Aggregate(ctx, "test_collection", pipeline, &results)
	if err != nil {
		t.Fatalf("Aggregate failed: %s", err)
	}

	// Check that both "Alice" and "Bob" are present with correct counts
	foundAlice, foundBob := false, false
	for _, r := range results {
		if r.ID == "Alice" && r.Count == 2 {
			foundAlice = true
		}
		if r.ID == "Bob" && r.Count == 1 {
			foundBob = true
		}
	}
	if !foundAlice || !foundBob {
		t.Errorf("Aggregate results incorrect: %+v", results)
	}

	// Test error on nil pipeline
	err = client.Aggregate(ctx, "test_collection", nil, &results)
	if err == nil {
		t.Errorf("Expected error for nil pipeline, got nil")
	}

	// Test error on nil result
	err = client.Aggregate(ctx, "test_collection", pipeline, nil)
	if err == nil {
		t.Errorf("Expected error for nil result, got nil")
	}
}
