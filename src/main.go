package main

import (
	"log"
	"net/http"
	"os"
	"github.com/codegangsta/negroni"
	//"gopkg.in/mgo.v2/bson"
)

func main() {
	initRSAKeys()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// route
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login2", loginHandler)
	http.HandleFunc("/login", newLoginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/post", postHandler)
	http.Handle("/list", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(listHandler)),
	))
	http.HandleFunc("/entry", entryHandler)

	http.Handle("/statics/",
		http.StripPrefix("/statics/", http.FileServer(http.Dir("./statics"))),
	)

	if err := http.ListenAndServe(":"+port, nil);err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


