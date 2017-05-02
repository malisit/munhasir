package main

import (
	"time"
	"gopkg.in/mgo.v2/bson"
)

type Day struct {
	Id				int64			`bson:"id"`
	DayDate			time.Time		`bson:"day_date"`
}

type User struct {
	Id				bson.ObjectId	`bson:"_id,omitempty"`
	Username		string			`bson:"username"`
	Password		string			`bson:"password"`
	Datetime		time.Time		`bson:"registration_datetime"`
}

type Entry struct {
	Id				bson.ObjectId	`bson:"_id,omitempty"`
	User			User			`bson:"user"`
	Title			string			`bson:"title"`
	Day				time.Time		`bson:"day"`
	Updated			time.Time		`bson:"updated"`
	EncryptedText	string			`bson:"encrypted_text"`
}

type UserCredentials struct {
	Username		string			`bson:"username" json:"username"`
	Password		string			`bson:"password" json:"password"`
}

type TokenUserPair struct {
	User			User			`bson:"user"`
	Token			string			`bson:"token"`
	Timestamp		time.Time		`bson:"timestamp"`
}

type OneWayStruct struct {
	One 			string			`json:"one"`
}

type TwoWayStruct struct {
	One 			string			`json:"one"`
	Two 			string			`json:"two"`
}

type ThreeWayStruct struct {
	One 			string			`json:"one"`
	Two 			string			`json:"two"`
	Three 			string			`json:"three"`
}

type FourWayStruct struct {
	One 			string			`json:"one"`
	Two 			string			`json:"two"`
	Three 			string			`json:"three"`
	Four 			string			`json:"four"`
}