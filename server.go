package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Blog struct {
	ID    primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	Title string             `json:"title, omitempty" bson:"title, omitempty"`
	Body  string             `json:"body, omitempty" bson:"body, omitempty"`
}

func GetBlogsListEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	fmt.Println("a")
	var blogs []Blog
	fmt.Println("a")
	collection := client.Database("BCH").Collection("blogs")
	fmt.Println("a")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	fmt.Println("a")
	cursor, err := collection.Find(ctx, bson.M{})
	fmt.Println("a")
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var blog Blog
		cursor.Decode(&blog)
		blogs = append(blogs, blog)
	}
	if err := cursor.Err(); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(res).Encode(blogs)
}

func CreateBlogEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var blog Blog
	_ = json.NewDecoder(req.Body).Decode(&blog)
	collection := client.Database("BCH").Collection("blogs")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, blog)
	json.NewEncoder(res).Encode(result)
}

func ReadBlogEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	vars := mux.Vars(req)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	collection := client.Database("BCH").Collection("blogs")
	var blog Blog
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err := collection.FindOne(ctx, Blog{ID: id}).Decode(&blog); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}

// func CloseDB() {
// 	err := client.Disconnect(context.TODO())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Connection to MongoDB closed.")
// }

func TestingEndPoint(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(Blog{})
}
func main() {
	fmt.Println("API Started")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/", TestingEndPoint).Methods("GET")
	router.HandleFunc("/blogs", GetBlogsListEndPoint).Methods("GET")
	router.HandleFunc("/createblog", CreateBlogEndPoint).Methods("POST")
	router.HandleFunc("/readblog/{id}", ReadBlogEndPoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
