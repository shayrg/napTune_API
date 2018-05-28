package user

import (
	"net/http"
	"encoding/json"
	"database/sql"
	"time"
	"github.com/nu7hatch/gouuid"
	"github.com/shayrg/napTune_API/global"
	_ "github.com/go-sql-driver/mysql"
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
	global.CheckErr(err)
	defer r.Body.Close()
	user = authenticateUser(user)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func Authorize(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	authorized := false
	var user User
	err := decoder.Decode(&user)
	global.CheckErr(err)
	defer r.Body.Close()
	var dbUser = getUser(user)
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

func buildUser(rows *sql.Rows) User {
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
		global.CheckErr(err)
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

func getUser(user User) User {
	var selectStatement string
	var selectValue string
	//Token
	if user.Token != "" {
		selectStatement = "token"
		selectValue = user.Token
		//Id
	} else if user.Id != "" {
		selectStatement = "id"
		selectValue = user.Id
		//Email
	} else if user.Email != "" {
		selectStatement = "email"
		selectValue = user.Email
		//Fail case
	} else {
		selectStatement = "id"
		selectValue = "-1"
	}
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	stmt, err := db.Prepare("select * from users where " + selectStatement + " = ?")
	global.CheckErr(err)
	rows, err := stmt.Query(selectValue)
	global.CheckErr(err)
	db.Close()
	return buildUser(rows)
}
func authenticateUser(loginUser User) User{
	dbUser := getUser(loginUser)
	if loginUser.Password == dbUser.Password {
		return setToken(dbUser)
	} else {
		return User{}
	}
}
func setToken(user User) User{
	expirationOffset := time.Hour * 2
	expirationDate := time.Now().UTC().Add(expirationOffset).Format("2006-01-02 15:04:05")
	user.Token = generateToken()
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	stmt, err := db.Prepare("update users set token = ?, expiration = ? where id = ?")
	global.CheckErr(err)
	_, err = stmt.Exec(user.Token, expirationDate, user.Id)
	global.CheckErr(err)
	return getUser(user)
}
func generateToken() string {
	token, err := uuid.NewV4()
	global.CheckErr(err)
	return token.String()
}