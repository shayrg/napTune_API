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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func Authorize(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	authorized := false
	var user User
	err := decoder.Decode(&user)
	checkErr(err)
	defer r.Body.Close()
	var dbUser = GetUser(user)
	//Convert to local
	dbUser.Expiration = dbUser.Expiration.Local()
	//Adjust for local
	dbUser.Expiration = dbUser.Expiration.Add(5*time.Hour)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if time.Now().Before(dbUser.Expiration) && dbUser.Roll == user.Roll {
		authorized = true
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(authorized)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(authorized)
	}
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