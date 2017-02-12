package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Day struct {
	Id             	int64   		`bson:"id"`
	DayDate    		time.Time   	`bson:"day_date"`
}

type User struct {
	Id      		bson.ObjectId	`bson:"_id,omitempty"`
	Username 		string			`bson:"username"`
	Password		string 			`bson:"password"`
	Datetime		time.Time 		`bson:"registration_datetime"`
}

type Entry struct {
	Id      		bson.ObjectId	`bson:"_id,omitempty"`
	User 			User			`bson:"user"`
	Day 			time.Time		`bson:"day"`
	EncryptedText	string			`bson:"encrypted_text"`
}

type KeyVIPair struct {		
	KeyMd5			string 			`bson:"keymd5"`
	VI 				string 			`bson:"vi"`
}