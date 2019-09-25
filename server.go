
package main

import (
	//"fmt"
	"encoding/json"
	"strconv"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	username string `json:"id, omitempty"`
	email string `json:"id, omitempty"`
	password string `json:"id, omitempty"`
	blogs []Blog `json:"blogs, omitempty"`
}

type Blog struct {
	BlogID string	`json:"id, omitempty"`
	Title string	`json:"title, omitempty"`
	Body string	`json:"body, omitempty"`
}
var client *mongo.Client

var blogs []Blog

/*func Handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}*/
func GetBlogsEndPoint(w http.ResponseWriter, req *http.Request){
	json.NewEncoder(w).Encode(blogs)
}

func GetBlogEndPoint(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	for _, blog := range blogs {
		if blog.BlogID == vars["id"]{
			json.NewEncoder(w).Encode(blog)
			return
		}
		//json.NewEncoder(w).Encode(x)
	}
	json.NewEncoder(w).Encode(Blog{BlogID : "null", Title : "null", Body : "null"})
	//fmt.Printf("%v", blogs)
}

func CreateBlogEndPoint(w http.ResponseWriter, req *http.Request){
	//fmt.Printf("%s", req.Body)
	var blog Blog
	_ = json.NewDecoder(req.Body).Decode(&blog)
	blog.BlogID = strconv.Itoa(len(blogs)+1)
	blogs = append(blogs,blog) 
	json.NewEncoder(w).Encode(blog)

}

func UpdateBlogEndPoint(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	var newblog Blog
	_ = json.NewDecoder(req.Body).Decode(&newblog)
	//fmt.Printf("%s %s %s", newblog.BlogID, newblog.Title, newblog.Body)
	for idx, blog := range blogs{
		if(blog.BlogID==vars["id"]){
			blog.Title = newblog.Title
			blog.Body = newblog.Body
			var blogscopy []Blog
			blogscopy = append(blogs[:idx], blog)
			blogs = append(blogscopy,blogs[idx+1:]...)
			return
		}
	}
	json.NewEncoder(w).Encode(blogs)
}
func DeleteBlogEndPoint(w http.ResponseWriter, req *http.Request){
	vars := mux.Vars(req)
	for i, blog := range blogs {
		if blog.BlogID==vars["id"]{
			blogs = append(blogs[:i], blogs[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(blogs)
}

func main(){
	blogs = append(blogs, Blog{BlogID : "1", Title : "Backend", Body : "Working with go"})
	blogs = append(blogs, Blog{BlogID : "2", Title : " Frontend", Body : "FrontEnd is ewwww"})
	r := mux.NewRouter()
	r.HandleFunc("/blogs", GetBlogsEndPoint).Methods("GET")
	r.HandleFunc("/blogs/{id}", GetBlogEndPoint).Methods("GET")
	r.HandleFunc("/blogs/create", CreateBlogEndPoint).Methods("POST")
	r.HandleFunc("/blogs/{id}", DeleteBlogEndPoint).Methods("DELETE")
	r.HandleFunc("/blogs/update/{id}", UpdateBlogEndPoint).Methods("POST")
	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", nil) ;err!=nil {
		log.Fatal(err)
	}
}
