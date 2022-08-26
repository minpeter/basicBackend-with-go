package main

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const dbname = "PersonDb"
const collectionName = "Person"

func indexRoute(res *fiber.Ctx) error {
	collection, err := getMongoDbCollection(dbname, collectionName)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	filter := bson.M{}
	curr, err := collection.Find(context.Background(), filter)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	defer curr.Close(context.Background())
	var result []bson.M
	curr.All(context.Background(), &result)

	json, _ := json.Marshal(result)
	return res.Status(200).Send(json)
}

func addPerson(res *fiber.Ctx) error {
	collection, err := getMongoDbCollection(dbname, collectionName)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	var newParson Person
	json.Unmarshal([]byte(res.Body()), &newParson)
	curr, err := collection.InsertOne(context.Background(), newParson)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	response, _ := json.Marshal(curr)
	return res.Status(200).Send(response)
}

func updatePerson(res *fiber.Ctx) error {
	collection, err := getMongoDbCollection(dbname, collectionName)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	id := res.Params("id")
	objId, _ := primitive.ObjectIDFromHex(id)
	var updatePerson Person
	json.Unmarshal([]byte(res.Body()), &updatePerson)
	filter := bson.M{
		"_id": objId,
	}
	update := bson.M{
		"$set": updatePerson,
	}
	curr, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	response, _ := json.Marshal(curr)
	return res.Status(200).Send(response)

}

func deletePerson(res *fiber.Ctx) error {
	collectionName, err := getMongoDbCollection(dbname, collectionName)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	id := res.Params("id")
	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": objId,
	}
	curr, err := collectionName.DeleteOne(context.Background(), filter)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	response, _ := json.Marshal(curr)
	return res.Status(200).Send(response)
}

func createPerson(res *fiber.Ctx) error {
	collection, err := getMongoDbCollection(dbname, "Assignment")
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	var newAssignment Assignment
	json.Unmarshal([]byte(res.Body()), &newAssignment)
	personCollection, err := getMongoDbCollection(dbname, collectionName)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	id := res.Params("id")
	newAssignment.Person = id
	objId, _ := primitive.ObjectIDFromHex(id)
	curr, _ := collection.InsertOne(context.Background(), newAssignment)
	filter := bson.M{
		"_id": objId,
	}
	var temp Person
	personCollection.FindOne(context.Background(), filter).Decode(&temp)
	temp.Assigment = append(temp.Assigment, curr.InsertedID.(primitive.ObjectID).Hex())
	update := bson.M{
		"$set": temp,
	}
	result, _ := personCollection.UpdateOne(context.Background(), filter, update)
	response, _ := json.Marshal(result)
	return res.Status(200).Send(response)
}

func getPersonAssignment(res *fiber.Ctx) error {
	collection, err := getMongoDbCollection(dbname, "Assignment")
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	id := res.Params("id")
	filter := bson.M{
		"person": id,
	}
	curr, err := collection.Find(context.Background(), filter)
	if err != nil {
		return res.Status(400).SendString("There is some Problem ! Please Try again")
	}
	defer curr.Close(context.Background())
	var result []bson.M
	curr.All(context.Background(), &result)
	json, _ := json.Marshal(result)
	return res.Status(200).Send(json)
}

func main() {
	app := fiber.New()
	app.Get("/", indexRoute)
	app.Post("/create", addPerson)
	app.Put("/update/:id", updatePerson)
	app.Delete("/delete/:id", deletePerson)
	app.Post("/assignment/:id", createPerson)
	app.Get("/assignment/:id", getPersonAssignment)
	app.Listen(":8000")
}
