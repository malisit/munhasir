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

	newUser := User{username, string(hashedPassword), createdAt}
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
		http.Redirect(w, r, "/",  http.StatusMovedPermanently)
	}

}