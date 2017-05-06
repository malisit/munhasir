package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"os"
)

func connect() (session *mgo.Session) {
	connectURL := "localhost"
	session, err := mgo.Dial(connectURL)
	if err != nil {
		fmt.Printf("Cannot connect to MongoDB, %v\n", err)
		os.Exit(1)
	}

	session.SetSafe(&mgo.Safe{})

	return session
}