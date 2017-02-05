package munhasir

import (
	"time"
)

type Day struct {
	Id             	int64   	`json:"id"`
	DayDate    		time.Time   `json:"day_date"`
}

type User struct {
	Username 		string		`json:"username"`
	Password		string 		`json:"password"`
	Email 			string 		`json:"email"`
	Datetime		time.Time 	`json:"registration_datetime"`
}

type Entry struct {
	User 			User		`json:"user"`
	Day 			Day 		`json:"day"`
	EncryptedText	string		`json:"encrypted_text"`
}

type KeyVIPair struct {		
	KeyMd5			string 		`json:"keymd5"`
	VI 				string 		`json:"vi"`
}