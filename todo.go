package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database   string
	collection string
)

const (
	// environment variables
	mongoDBConnectionStringEnvVarName = "MONGODB_CONNECTION_STRING"
	mongoDBDatabaseEnvVarName         = "MONGODB_DATABASE"
	mongoDBCollectionEnvVarName       = "MONGODB_COLLECTION"

	// status
	statusPending   = "pending"
	statusCompleted = "completed"
	listAllCriteria = "all"
	statusAttribute = "status"

	// flags (commands)
	createFlag = "create"
	listFlag   = "list"
	updateFlag = "update"
	deleteFlag = "delete"

	// help text
	createHelp = "create a todo: enter description. e.g. todo -create \"get milk\""
	listHelp   = "list all, pending or completed todos. e.g. todo -list <criteria> (criteria can be all, pending or completed"
	updateHelp = "update a todo: enter todo ID and new status e.g. todo -update <id>,<new status> e.g. todo -update 1,completed"
	deleteHelp = "delete a todo: enter todo ID e.g. todo -delete 42"
)

func main() {
	todoDescription := flag.String("create", "", createHelp)
	listCriteria := flag.String("list", "", listHelp)
	updateInfo := flag.String("update", "", updateHelp)
	deleteTodo := flag.String("delete", "", deleteHelp)

	flag.Parse()

	if len(os.Args) > 3 {
		log.Fatalf("incorrect usage. please use 'todo --help'")
	}
	if *todoDescription != "" {
		create(*todoDescription)
	}

	if *listCriteria != "" {
		list(*listCriteria)
	}

	if *updateInfo != "" {
		if !strings.Contains(*updateInfo, ",") {
			log.Fatalf("invalid update info. please use 'todo --help'")
		}
		todoid := strings.Split(*updateInfo, ",")[0]
		newStatus := strings.Split(*updateInfo, ",")[1]
		update(todoid, newStatus)
	}

	if *deleteTodo != "" {
		delete(*deleteTodo)
	}
}

// connects to MongoDB
func connect() *mongo.Client {
	mongoDBConnectionString := os.Getenv(mongoDBConnectionStringEnvVarName)
	if mongoDBConnectionString == "" {
		log.Fatal("missing environment variable: ", mongoDBConnectionStringEnvVarName)
	}

	database = os.Getenv(mongoDBDatabaseEnvVarName)
	if database == "" {
		log.Fatal("missing environment variable: ", mongoDBDatabaseEnvVarName)
	}

	collection = os.Getenv(mongoDBCollectionEnvVarName)
	if collection == "" {
		log.Fatal("missing environment variable: ", mongoDBCollectionEnvVarName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(true)
	c, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatalf("unable to initialize connection %v", err)
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("unable to connect %v", err)
	}
	return c
}

// creates a todo
func create(desc string) {
	c := connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	todoCollection := c.Database(database).Collection(collection)
	r, err := todoCollection.InsertOne(ctx, Todo{Description: desc, Status: statusPending})
	if err != nil {
		log.Fatalf("failed to add todo %v", err)
	}
	fmt.Println("added todo", r.InsertedID)
}

// lists todos
func list(status string) {

	var filter interface{}
	switch status {
	case listAllCriteria:
		filter = bson.D{}
	case statusCompleted:
		filter = bson.D{{statusAttribute, statusCompleted}}
	case statusPending:
		filter = bson.D{{statusAttribute, statusPending}}
	default:
		log.Fatal("invalid criteria for listing todo(s)")
	}

	c := connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	todoCollection := c.Database(database).Collection(collection)
	rs, err := todoCollection.Find(ctx, filter)
	if err != nil {
		log.Fatalf("failed to list todo(s) %v", err)
	}
	var todos []Todo
	err = rs.All(ctx, &todos)
	if err != nil {
		log.Fatalf("failed to list todo(s) %v", err)
	}
	if len(todos) == 0 {
		fmt.Println("no todos found")
		return
	}

	todoTable := [][]string{}

	for _, todo := range todos {
		s, _ := todo.ID.MarshalJSON()
		todoTable = append(todoTable, []string{string(s), todo.Description, todo.Status})
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Description", "Status"})

	for _, v := range todoTable {
		table.Append(v)
	}
	table.Render()
}

// updates a todo
func update(todoid, newStatus string) {
	c := connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	todoCollection := c.Database(database).Collection(collection)
	oid, err := primitive.ObjectIDFromHex(todoid)
	if err != nil {
		log.Fatalf("failed to update todo %v", err)
	}
	filter := bson.D{{"_id", oid}}
	update := bson.D{{"$set", bson.D{{statusAttribute, newStatus}}}}
	_, err = todoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("failed to update todo %v", err)
	}
}

// deletes a todo
func delete(todoid string) {
	c := connect()
	ctx := context.Background()
	defer c.Disconnect(ctx)

	todoCollection := c.Database(database).Collection(collection)
	oid, err := primitive.ObjectIDFromHex(todoid)
	if err != nil {
		log.Fatalf("invalid todo ID %v", err)
	}
	filter := bson.D{{"_id", oid}}
	_, err = todoCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf("failed to delete todo %v", err)
	}
}

// Todo represents a todo
type Todo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Description string             `bson:"description"`
	Status      string             `bson:"status"`
}
