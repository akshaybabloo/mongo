package mongo_test

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/akshaybabloo/mongo/v2"
)

func ExampleClient_Add() {

	type data struct {
		ID   string `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

	testData := data{
		ID:   "1",
		Name: "Akshay",
	}

	done, err := client.Add("test_collection", testData)
	if err != nil {
		panic(err)
	}
	fmt.Println("The ID is:", done.InsertedID)
}

func ExampleClient_AddMany() {

	type data struct {
		ID   string `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

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
		panic(err)
	}
	fmt.Println("The ID is:", done.InsertedIDs)
}

func ExampleClient_Delete() {
	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

	deleted, err := client.Delete("test_collection", "1")
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted items:", deleted.DeletedCount)
}

func ExampleClient_Update() {
	type data struct {
		Name string `bson:"name"`
	}

	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

	testData := data{
		Name: "Akshay",
	}

	updated, err := client.Update("test_collection", "1", testData)
	if err != nil {
		panic(err)
	}
	fmt.Println("Modified items:", updated.ModifiedCount)
}

func ExampleClient_Get() {

	type data struct {
		ID   int    `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

	var decodeData data
	output := client.Get("test_collection", "2").Decode(&decodeData)
	if output != nil {
		panic("No data found.")
	}
	fmt.Println(decodeData)
}

func ExampleClient_GetCustom() {

	type data struct {
		ID   int    `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

	var decodeData data
	output := client.GetCustom("test_collection", bson.M{"id": "2"}).Decode(&decodeData)
	if output != nil {
		panic("No data found.")
	}
	fmt.Println(decodeData)
}

func ExampleClient_GetAll() {

	type data struct {
		ID   string `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

	var testData []data
	err := client.GetAll("test_collection", "1", &data{})
	if err != nil {
		panic(err)
	}
	fmt.Println("The ID is:", testData)
}

func ExampleClient_GetAllCustom() {

	type data struct {
		ID   string `bson:"id"`
		Name string `bson:"name"`
	}

	client := mongo.NewMongoClient("mongodb://localhost:27017/?retryWrites=true&w=majority", "test")

	var testData []data
	err := client.GetAllCustom("test_collection", bson.M{"id": "1"}, &data{})
	if err != nil {
		panic(err)
	}
	fmt.Println("The ID is:", testData)
}
