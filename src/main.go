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
	http.Handle("/",
		http.StripPrefix("", http.FileServer(http.Dir("./statics"))),
	)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/register", registerHandler)
	http.Handle("/api/list", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(listHandler)),
	))
	http.Handle("/api/post", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(postHandler)),
	))
	http.Handle("/api/entry", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(entryHandler)),
	))
	http.Handle("/api/decrypt", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(decryptHandler)),
	))

	// http.Handle("/api/token", negroni.New(
	// 	negroni.HandlerFunc(ValidateTokenMiddleware),
	// 	negroni.Wrap(http.HandlerFunc(isTokenValid)),
	// ))

	

	if err := http.ListenAndServe(":"+port, nil);err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


