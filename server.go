
package main

import (
	"fmt"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main(){
	
	http.HandleFunc("/", Handler)
	err := http.ListenAndServe(":8080", nil)
	if err!=nil {
		panic(err)
	}
	fmt.Println("Hello World")
}
