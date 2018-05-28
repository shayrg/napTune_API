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

type user struct {
	Id			string `json:"id"`
	FirstName 	string `json:"first_name"`
	LastName 	string `json:"last_name"`
	Email 		string `json:"email"`
	Password	string `json:"password"`
	Token 		string `json:"token"`
	Expiration	time.Time  `json:"expiration"`
	Roll 		string `json:"roll"`
}
type authorized struct {
	Authorized bool
}
const hoursOffset = 2

func Authenticate(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var usr user
	err := decoder.Decode(&usr)
	global.CheckErr(err)
	defer r.Body.Close()
	dbUser := authenticateUser(usr)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	//Email must but set
	if usr.Email != "" && dbUser.Id != ""{
		dbUser.Password = ""
		dbUser.Id = "-1"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dbUser)
	} else {
		auth := authorized{
			Authorized: false,
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(auth)
	}
}

func Authorize(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	auth := authorized{
		Authorized: false,
	}
	var usr user
	var dbUser user
	err := decoder.Decode(&usr)
	global.CheckErr(err)
	defer r.Body.Close()
	//Token must be set
	if usr.Token != "" {
		dbUser = getUser(usr)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if time.Now().Before(dbUser.Expiration) && dbUser.Roll == usr.Roll {
		auth.Authorized = true
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	json.NewEncoder(w).Encode(auth)

}

func Logout(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var usr user
	err := decoder.Decode(&usr)
	global.CheckErr(err)
	defer r.Body.Close()
	auth := authorized{
		Authorized: false,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if logoutUser(usr) {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
	json.NewEncoder(w).Encode(auth)
}

func NewUser(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	var usr user
	err := decoder.Decode(&usr)
	global.CheckErr(err)
	defer r.Body.Close()
	success := false
	//Check email, name, and password are set
	if usr.Email != "" && usr.Password != "" && usr.FirstName != "" && usr.LastName != "" {
		//Check email is not used
		usr.Token = ""
		dbUser := getUser(usr)
		if dbUser.Id == "" {
			success = true
		}
	}
	if success {
		//Add user to database
		db, err := sql.Open("mysql", global.DbString)
		global.CheckErr(err)
		stmt, err := db.Prepare("insert into users (" +
			"firstName, lastName, email, password, role, expiration) values (?,?,?,?,?,?)")
		global.CheckErr(err)
		_, err = stmt.Exec(usr.FirstName, usr.LastName, usr.Email, usr.Password, "user",
			time.Now().UTC().UTC().Format("2006-01-02 15:04:05"))
		global.CheckErr(err)
		db.Close()
		usr = getUser(usr)
		usr = setToken(usr, hoursOffset)
		if usr.Token == "" {
			success = false
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if success {
		w.WriteHeader(http.StatusOK)
		usr.Id = "-1"
		usr.Password = ""
		json.NewEncoder(w).Encode(usr)
	} else {
		auth := authorized{
			Authorized: false,
		}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(auth)
	}
}

func buildUser(rows *sql.Rows) user {
	var usr user
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
		usr = user {
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
	return usr
}

func getUser(usr user) user {
	var selectStatement string
	var selectValue string
	//Token
	if usr.Token != "" {
		selectStatement = "token"
		selectValue = usr.Token
		//Id
	} else if usr.Id != "" {
		selectStatement = "id"
		selectValue = usr.Id
		//Email
	} else if usr.Email != "" {
		selectStatement = "email"
		selectValue = usr.Email
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

func authenticateUser(loginUser user) user{
	dbUser := getUser(loginUser)
	if loginUser.Password == dbUser.Password {
		return setToken(dbUser, hoursOffset)
	} else {
		return user{}
	}
}

func logoutUser(usr user) bool{
	success := false
	//Token must be set
	if usr.Token != "" {
		dbUser := getUser(usr)
		dbUser = setToken(dbUser,0)
		//If expiration is in the past
		if time.Now().After(dbUser.Expiration) && dbUser.Id != "" {
			success = true
		}
	}
	return success
}

func setToken(usr user, offset int) user{
	expirationOffset := time.Duration(offset) * time.Hour
	expirationDate := time.Now().UTC().Add(expirationOffset).Format("2006-01-02 15:04:05")
	usr.Token = generateToken()
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	stmt, err := db.Prepare("update users set token = ?, expiration = ? where id = ?")
	global.CheckErr(err)
	_, err = stmt.Exec(usr.Token, expirationDate, usr.Id)
	global.CheckErr(err)
	db.Close()
	return getUser(usr)
}

func generateToken() string {
	token, err := uuid.NewV4()
	global.CheckErr(err)
	return token.String()
}