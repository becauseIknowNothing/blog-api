package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

type Author struct {
	AuthorName string `json:"authorname" bson:"authorname"`
	AuthorImg  string `json:"authorimg" bson:"authorimg"`
}

type Blog struct {
	Title        string `json:"title" bson:"title"`
	AuthorInfo   Author `json:"author" bson:"author"`
	Body         string `json:"body" bson:"body`
	LastModified string `json:"lastmodified" bson:"lastmodified"`
	// ModifiedAtDate time.Date
}

func GetBlogsListEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var blogs []Blog
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
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
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	blog.LastModified = time.Now().String()
	result, _ := collection.InsertOne(ctx, blog)
	fmt.Println(result)
}

func ReadBlogEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	vars := mux.Vars(req)
	title, _ := vars["title"]
	findoptions := options.Find()
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findoptions)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	var blogs []Blog
	for cur.Next(context.TODO()) {
		var blog Blog
		err := cur.Decode(&blog)
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(strings.ToLower(blog.Title), strings.ToLower(title)) {
			blogs = append(blogs, blog)
		}
	}
	json.NewEncoder(res).Encode(blogs)
}

func UpdateBlogEndPoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	vars := mux.Vars(req)
	title, _ := vars["title"]
	var blog Blog
	_ = json.NewDecoder(req.Body).Decode(&blog)
	filter := bson.D{{"title", title}}
	update := bson.M{"$set": bson.M{"title": blog.Title, "body": blog.Body, "lastmodified": time.Now().String()}}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}

func DeleteBlogEndPoint(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	title, _ := vars["title"]
	_, err := collection.DeleteOne(context.TODO(), bson.M{"title": title})
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
}

func CloseDB() {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}

func TestingEndPoint(res http.ResponseWriter, req *http.Request) {
	json.NewEncoder(res).Encode(Blog{})
}
func main() {
	fmt.Println("API Started")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	defer CloseDB()
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("BCH").Collection("blogs")
	router := mux.NewRouter()
	router.HandleFunc("/", TestingEndPoint).Methods("GET")
	router.HandleFunc("/blogs", GetBlogsListEndPoint).Methods("GET")
	router.HandleFunc("/createblog", CreateBlogEndPoint).Methods("POST")
	router.HandleFunc("/readblog/{title}", ReadBlogEndPoint).Methods("GET")
	router.HandleFunc("/updateblog/{title}", UpdateBlogEndPoint).Methods("POST")
	router.HandleFunc("/deleteblog/{title}", DeleteBlogEndPoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
