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
	http.HandleFunc("/api/user/login", loginHandler)
	http.HandleFunc("/api/user/register", registerHandler)
	http.Handle("/api/user/password", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(changePassword)),
	))
	http.Handle("/api/user/delete", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(deleteUser)),
	))
	http.Handle("/api/entry/list", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(listHandler)),
	))
	http.Handle("/api/entry/post", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(postHandler)),
	))
	http.Handle("/api/entry/view", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(entryHandler)),
	))
	http.Handle("/api/entry/decrypt", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(decryptHandler)),
	))
	http.Handle("/api/entry/delete", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(deleteHandler)),
	))
	http.Handle("/api/entry/edit", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(editHandler)),
	))

	// http.Handle("/api/token", negroni.New(
	// 	negroni.HandlerFunc(ValidateTokenMiddleware),
	// 	negroni.Wrap(http.HandlerFunc(isTokenValid)),
	// ))

	

	if err := http.ListenAndServe(":"+port, nil);err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


