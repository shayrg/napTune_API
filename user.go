package main

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"time"
)

type User struct {
	Id			string `json:"id"`
	FirstName 	string `json:"first_name"`
	LastName 	string `json:"last_name"`
	Email 		string `json:"email"`
	Password	string `json:"password"`
	Token 		string `json:"token"`
	Expiration	time.Time  `json:"expiration"`
	Roll 		string `json:"roll"`
}

func Authenticate(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var user User
	err := decoder.Decode(&user)
	checkErr(err)
	defer r.Body.Close()
	user = AuthenticateUser(user)
	Authorize(user.Token,user.Roll)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func Authorize(token string, roll string) bool{
	authorized := false
	user := User {
		Token: token,
	}
	user = GetUser(user)
	//Convert to local
	user.Expiration = user.Expiration.Local()
	//Adjust for local
	user.Expiration = user.Expiration.Add(5*time.Hour)
	if time.Now().Before(user.Expiration) && user.Roll == roll {
		authorized = true
	}
	return authorized
}
func BuildUser(rows *sql.Rows) User {
	var user User
	for rows.Next() {
		var id			string
		var firstName 	string
		var lastName 	string
		var email 		string
		var password	string
		var token 		string
		var expiration	time.Time
		var roll 		string
		err := rows.Scan(&id, &firstName, &lastName, &email, &password, &token, &expiration, &roll)
		checkErr(err)
		user = User {
			Id: id,
			FirstName: firstName,
			LastName: lastName,
			Email: email,
			Password: password,
			Token: token,
			Expiration: expiration,
			Roll: roll,
		}
	}
	return user
}