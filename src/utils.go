package main

import (
	"net/http"
	"crypto/rsa"
	"io/ioutil"
	"log"
	"fmt"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	VerifyKey *rsa.PublicKey
	SignKey *rsa.PrivateKey
)


func initRSAKeys(){
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

	json, err :=  json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	//validate token
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error){
		return VerifyKey, nil
	})



	if err == nil {

		if token.Valid{
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorised access to this resource")
	}

}
