package main

import (
	"net/http"
	"html/template"
	"time"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var authenticated = false
var loggedUser = User{}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(authenticated)
	if r.Method != "GET" {
		http.Error(w, "Method is not allowed.", http.StatusBadRequest)
	}

	var funcMap = template.FuncMap{
		"titleF": func() string {
			return "munhasir"
		},

	}

	template, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info := make(map[string]string)
	info["title"] = "munhasir"
	info["content"] = "A platform to keeping diaries for those who are cautious(or paranoid)."

	err = template.Execute(w, info)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/register.html")
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	// password2 := r.FormValue("password2")
	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("users")

	result := User{}
	
	err := collection.Find(bson.M{"username":username}).Select(bson.M{"username":username}).One(&result)
	
	if err == nil{
		http.Error(w, "already registered username", http.StatusInternalServerError)
		return
	}

	createdAt := time.Now()
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	checkInternalServerError(err, w)

	newUser := User{Username:username, Password:string(hashedPassword), Datetime:createdAt}
	err = collection.Insert(newUser)

	if err != nil{
		http.Error(w, "error: " + err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/login", http.StatusMovedPermanently)	
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/login.html")
		return
	}
	
	username := r.FormValue("username")
	password := r.FormValue("password")
	result := User{}

	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("users")
	err := collection.Find(bson.M{"username":username}).One(&result)
	if err != nil{
		http.Error(w, "user doesn't exist", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))

	if err != nil {
		http.Error(w, "password is not true", http.StatusInternalServerError)
		return
	} else {
		authenticated = true
		loggedUser = result
		http.Redirect(w, r, "/",  http.StatusMovedPermanently)
	}

}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/post.html")
		return
	}
	uncryptedText := r.FormValue("text")
	unhashedKey := r.FormValue("key")
	user := loggedUser
	createdAt := time.Now()

	hashedKey := hash(unhashedKey)

	encryptedText := encrypt(hashedKey, uncryptedText)

	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("entries")

	newEntry := Entry{User:user, Day:createdAt, EncryptedText:encryptedText}
	err := collection.Insert(newEntry)

	if err != nil{
		http.Error(w, "error: " + err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/list", http.StatusMovedPermanently)
}


func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not allowed.", http.StatusBadRequest)
	}

	user := loggedUser

	results := []Entry{}

	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("entries")
	err := collection.Find(bson.M{"user":user}).All(&results)

	if err != nil{
		http.Error(w, "error: " + err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(results)

}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	// decode and return entry via post
}
