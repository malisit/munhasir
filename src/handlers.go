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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCredentials

	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	session := connect()
	defer session.Close()
	result := User{}
	collection := session.DB("munhasir").C("users")
	err = collection.Find(bson.M{"username":user.Username}).One(&result)
	if err != nil{
		JsonResponse("user doesn't exist", w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		JsonResponse("password is not true", w)
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

	// create a user-token pair

	collection = session.DB("munhasir").C("usertoken")
	newPair := TokenUserPair{User:result, Token:tokenString, Timestamp: bson.Now()}
	err = collection.Insert(newPair)

	if err != nil{
		http.Error(w, "everything is something happened: " + err.Error(), http.StatusBadRequest)
		return
	}

	// return token string
	JsonResponse(tokenString, w)


}

func indexHandler(w http.ResponseWriter, r *http.Request) {
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

// func isTokenValid(w http.ResponseWriter, r *http.Request) {
// 	JsonResponse("valid",w)
// }

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		JsonResponse("wrong method", w)
		return
	}
	var user UserCredentials

	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}

	username := user.Username
	password := user.Password

	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("users")

	result := User{}

	err = collection.Find(bson.M{"username":username}).Select(bson.M{"username":username}).One(&result)

	if err == nil{
		JsonResponse("already registered username", w)
		return
	}

	createdAt := time.Now()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	checkInternalServerError(err, w)

	newUser := User{Username:username, Password:string(hashedPassword), Datetime:createdAt}
	err = collection.Insert(newUser)

	if err != nil{
		JsonResponse("db error", w)
		return
	}
	JsonResponse("success", w)
}



func postHandler(w http.ResponseWriter, r *http.Request) {
	var entry EntryPost

	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		JsonResponse("error in request",w)
		return
	}


	uncryptedText := entry.Text
	unhashedKey := entry.Key

	token := w.Header().Get("token")

	// get user by token
	user := getUserByToken(token)

	createdAt := time.Now()

	hashedKey := hash(unhashedKey)

	encryptedText := encrypt(hashedKey, uncryptedText)

	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("entries")

	newEntry := Entry{User:user, Day:createdAt, EncryptedText:encryptedText}
	err = collection.Insert(newEntry)

	if err != nil{
		JsonResponse("db error",w)
		return
	}
	JsonResponse("success", w)
}


func listHandler(w http.ResponseWriter, r *http.Request) {

	token := w.Header().Get("token")

	// get user by token
	user := getUserByToken(token)

	results := []Entry{}

	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("entries")
	err := collection.Find(bson.M{"user":user}).All(&results)

	if err != nil{
		http.Error(w, "error: " + err.Error(), http.StatusBadRequest)
		return
	}
	
	JsonResponse(results, w)
}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	var entry EntryPost

	//decode request into UserCredentials struct
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error in request")
		return
	}


	encryptedText := entry.Text
	unhashedKey := entry.Key

	hashedKey := hash(unhashedKey)
	decryptedText := decrypt(hashedKey, encryptedText)

	JsonResponse(decryptedText, w)
	
}
