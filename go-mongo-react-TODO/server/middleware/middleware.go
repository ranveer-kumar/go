package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ranveer/go/go-mongo-react-TODO/server/models"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

const connectionString = "mongodb://localhost:27017"
const dbName = "go"
const collName = "todolist"

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)

	client, error := mongo.Connect(context.TODO(), clientOptions)

	if error != nil {
		log.Fatal(error)
	}
	err := client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected to MongoDB..")
	collection = client.Database(dbName).Collection(collName)
	log.Println("Collection instace created..")
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, `product service is good`)
	r.Body.Close()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
func GetAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllTask()
	json.NewEncoder(w).Encode(payload)
}

// CreateTask create task route
func CreateTask(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var task models.ToDoList
	_ = json.NewDecoder(request.Body).Decode(&task)
	// fmt.Println(task, r.Body)
	insertOneTask(task)

	json.NewEncoder(w).Encode(task)
}

// DeleteTask delete one task route
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	deleteOneTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
	// json.NewEncoder(w).Encode("Task not found")

}

// Delete all Task route
func DeleteAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// params := mux.Vars(r)
	deleteAllTask()
	// json.NewEncoder(w).Encode(params["id"])
	// json.NewEncoder(w).Encode("Task not found")

}

// Complete one task route
func TaskComplete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	taskComplete(params["id"])
	json.NewEncoder(w).Encode(params["id"])
	// json.NewEncoder(w).Encode("Task not found")

}

// UndoTask undo the complete task route
func UndoTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)
	undoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

// get all task from the DB and return it
func getAllTask() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

func insertOneTask(task models.ToDoList) {
	insertResult, err := collection.InsertOne(context.Background(), task)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Inserted a Single Record ", insertResult.InsertedID)
}

// delete one task from the DB, delete by ID
func deleteOneTask(task string) {
	log.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Deleted Document", d.DeletedCount)
}

// delete all task from the DB
func deleteAllTask() {
	log.Println("Deleting all tasks from db..")
	// id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{}
	d, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("all documents Deleted", d.DeletedCount)
}

// complete one task from the DB, delete by ID
func taskComplete(task string) {
	log.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("modified count: ", result.ModifiedCount)
}

// task undo method, update task's status to false
func undoTask(task string) {
	log.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": false}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("modified count: ", result.ModifiedCount)
}
