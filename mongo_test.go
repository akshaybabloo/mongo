package mongo

import "testing"

var client NewMongoDbClient

type data struct {
	id   int
	name string
}

func init() {
	client = NewMongoDbClient{
		ConnectionUrl: "mongodb://localhost:27017/?retryWrites=true&w=majority",
		DatabaseName:  "test",
	}
}

func TestNewMongoDbClient_Add(t *testing.T) {
	testData := data{
		id:   1,
		name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		t.Errorf("Unable to add data. %s", err)
	}
	t.Logf("The ID is %d", done.InsertedID)
}

func TestNewMongoDbClient_Delete(t *testing.T) {

}

func TestNewMongoDbClient_Get(t *testing.T) {

}

func TestNewMongoDbClient_Update(t *testing.T) {

}

func TestNewMongoDbClient_Collection(t *testing.T) {

}

func TestNewMongoDbClient_DB(t *testing.T) {

}
