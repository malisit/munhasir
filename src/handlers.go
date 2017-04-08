package main

import (
	"net/http"
	"html/template"
	"time"
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
)

var authenticated = false
var loggedUser = User{}

func denemeHandler(w http.ResponseWriter, r *http.Request) {
}


func newLoginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	fmt.Println(user)

	session := connect()
	defer session.Close()
	result := User{}
	collection := session.DB("munhasir").C("users")
	err = collection.Find(bson.M{"username":user.Username}).One(&result)
	if err != nil{
		http.Error(w, "user doesn't exist", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		http.Error(w, "password is not true", http.StatusInternalServerError)
		return
	}


	//create a rsa 256 signer
	signer := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := make(jwt.MapClaims)
	//set claims
	claims["iss"] = "admin"
	claims["exp"] = time.Now().Add(time.Minute * 20).Unix()
	claims["CustomUserInfo"] = struct {
		Name	string
		Role	string
	}{user.Username, "Member"}
	signer.Claims = claims
	tokenString, err := signer.SignedString(SignKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		log.Printf("Error signing token: %v\n", err)
	}

	//create a token instance using the token string
	JsonResponse(tokenString, w)


}

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
	fmt.Println("pass: ",password)
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

	template, err := template.New("list.html").ParseFiles("templates/list.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = template.Execute(w, results)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	// decode and return entry via post
	if r.Method != "POST" {
		http.Redirect(w, r, "/list", http.StatusMovedPermanently)
	} else {
		encryptedText := r.FormValue("text")
		unhashedKey := r.FormValue("key")
		hashedKey := hash(unhashedKey)
		decryptedText := decrypt(hashedKey, encryptedText)

		result := make(map[string]string)
		result["text"] = decryptedText

		template, err := template.New("entry.html").ParseFiles("templates/entry.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = template.Execute(w, result)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
