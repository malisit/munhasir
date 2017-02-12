package main

import (
	"log"
	"net/http"
	"os"
	//"gopkg.in/mgo.v2/bson"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	
	// route
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)

	if authenticated {
		http.HandleFunc("/post", postHandler)
		http.HandleFunc("/list", listHandler)
		http.HandleFunc("/entry", entryHandler)

	}
	
	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)
	http.ListenAndServe(":" + port, nil)
}