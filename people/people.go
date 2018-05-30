package people

import (
	"net/http"
	"database/sql"
	"github.com/shayrg/napTune_API/global"
	"encoding/json"
	"github.com/gorilla/mux"
)

type artist struct {
	Id			string `json:"id"`
	FirstName 	string `json:"first_name"`
	LastName 	string `json:"last_name"`
	Roll 		string `json:"roll"`
}
type artists []artist

func GetAllArtists(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	rows, err := db.Query("SELECT * FROM people where roll like 'artist'")
	global.CheckErr(err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildArtistList(rows))
}

func GetArtistById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistId := vars["artistId"]
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	stmt, err := db.Prepare("SELECT * FROM people where roll like 'artist' and id = ?")
	global.CheckErr(err)
	rows, err := stmt.Query(artistId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildArtistList(rows))
}

func buildArtistList(rows *sql.Rows) artists {
	var artList artists
	for rows.Next() {
		var id			string
		var firstName 	string
		var lastName 	string
		var roll 		string
		err := rows.Scan(&id, &firstName, &lastName, &roll)
		global.CheckErr(err)
		art := artist {
			Id: id,
			FirstName: firstName,
			LastName: lastName,
			Roll: roll,
		}
		artList = append(artList, art)
	}
	return artList
}