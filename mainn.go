package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(w, request)
		return
	}
	w.Write([]byte("Display home page"))
}
func showSnippet(w http.ResponseWriter, request *http.Request) {
	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	fmt.Println(id)
	if err != nil || id < 1 {
		http.NotFound(w, request)
		return
	}

	w.Write([]byte(fmt.Sprintf("Showing snippet with ID: %d", id)))
}
func createSnippet(w http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		w.Header().Set("Allow", "POST")
		//w.WriteHeader(405)
		//w.Write([]byte("Method not allowed"))
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Write([]byte("Create snippet"))
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/create/snippet", createSnippet)
	log.Println("Starting server on port :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
