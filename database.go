package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nu7hatch/gouuid"
)

const dbString = "root:mysql@tcp(localhost:3306)/napTune?charset=utf8&parseTime=True"

func GetUser(user User) User {
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
	db, err := sql.Open("mysql", dbString)
	checkErr(err)
	stmt, err := db.Prepare("select * from users where " + selectStatement + " = ?")
	checkErr(err)
	rows, err := stmt.Query(selectValue)
	checkErr(err)
	db.Close()
	return BuildUser(rows)
}
func AuthenticateUser(loginUser User) User{
	dbUser := GetUser(loginUser)
	if loginUser.Password == dbUser.Password {
		return setToken(dbUser)
	} else {
		return User{}
	}
}
func setToken(user User) User{
	user.Token = generateToken()
	//user.Expiration = time.Now().Add(2*time.Hour).String()
	db, err := sql.Open("mysql", dbString)
	checkErr(err)
	stmt, err := db.Prepare("update users set token = ?, expiration = NOW() where id = ?")
	checkErr(err)
	_, err = stmt.Exec(user.Token, user.Id)
	checkErr(err)
	_, err = db.Exec("update users set expiration = ADDTIME(expiration, '02:00:00')")
	//checkErr(err)
	return GetUser(user)
}
func generateToken() string {
	token, err := uuid.NewV4()
	checkErr(err)
	return token.String()
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}