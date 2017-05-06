package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func directToHttps(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.URL.Scheme == "https" || strings.HasPrefix(r.Proto, "HTTPS") || r.Header.Get("X-Forwarded-Proto") == "https" {
		next(w, r)
	} else {
		next(w, r)
		target := "https://" + r.Host + r.URL.Path

		http.Redirect(w, r, target,
			http.StatusTemporaryRedirect)
	}
}

var (
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)

func initRSAKeys() {
	signBytes, err := ioutil.ReadFile("app.rsa")

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)

	if err != nil {
		log.Fatal("Error reading private key")
		return
	}

	verifyBytes, err := ioutil.ReadFile("app.rsa.pub")

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)

	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
}

func JsonResponse(response interface{}, w http.ResponseWriter) {
	w.Header()["Content-Type"] = []string{"application/json"}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	w.WriteHeader(http.StatusOK)

	enc.Encode(response)
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//validate token
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return VerifyKey, nil
	})

	if err == nil {
		if token.Valid {
			w.Header().Set("token", token.Raw)
			next(w, r)
		} else {
			JsonResponse("this token is not a real token", w)
			return
		}
	} else {
		JsonResponse("this token is not authorized for this content", w)
		return
	}

}

func getUserByToken(token string) User {
	result := TokenUserPair{}

	session := connect()
	defer session.Close()

	end := bson.Now()
	start := end.Add(-20 * time.Minute)

	collection := session.DB("munhasir").C("usertoken")
	err := collection.Find(bson.M{"timestamp": bson.M{"$gte": start, "$lte": end}, "token": token}).One(&result)

	if err != nil {
		fmt.Println("there is no user token pair for given token")
	}
	return result.User
}

func getEntryById(id string) Entry {
	session := connect()
	defer session.Close()

	collection := session.DB("munhasir").C("entries")

	result := Entry{}
	err := collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)

	if err != nil {
		return Entry{}
	} else {
		return result
	}
}