package mongo_test

import (
	"fmt"

	"github.com/akshaybabloo/mongo"
)

func ExampleNewMongoDbClient_Add() {

	type data struct {
		Id   string    `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoDbClient{
		ConnectionUrl: "mongodb://localhost:27017/?retryWrites=true&w=majority",
		DatabaseName:  "test",
	}

	testData := data{
		Id:   "1",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		panic(err)
	}
	fmt.Println("The ID is:", done.InsertedID)
}

func ExampleNewMongoDbClient_Delete() {
	client := mongo.NewMongoDbClient{
		ConnectionUrl: "mongodb://localhost:27017/?retryWrites=true&w=majority",
		DatabaseName:  "test",
	}

	deleted, err := client.Delete("test_collection", "1")
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted items:", deleted.DeletedCount)
}

func ExampleNewMongoDbClient_Update() {
	type data struct {
		Name string `bson:"name"`
	}

	client := mongo.NewMongoDbClient{
		ConnectionUrl: "mongodb://localhost:27017/?retryWrites=true&w=majority",
		DatabaseName:  "test",
	}

	testData := data{
		Name: "Akshay",
	}

	updated, err := client.Update("test_collection", "1", testData)
	if err != nil {
		panic(err)
	}
	fmt.Println("Modified items:", updated.ModifiedCount)
}

func ExampleNewMongoDbClient_Get() {

	type data struct {
		Id   int    `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoDbClient{
		ConnectionUrl: "mongodb://localhost:27017/?retryWrites=true&w=majority",
		DatabaseName:  "test",
	}

	var decodeData data
	output := client.Get("test_collection", "2").Decode(&decodeData)
	if output != nil {
		panic("No data found.")
	}
	fmt.Println(decodeData)
}
