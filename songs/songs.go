package songs

import (
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	"github.com/shayrg/napTune_API/global"
	"encoding/json"
)

type song struct {
	Id 					string 	`json:"id"`
	Name				string 	`json:"name"`
	ArtistFirstName 	string 	`json:"first_name"`
	ArtistLastName 		string 	`json:"last_name"`
	Length 				string 	`json:"length"`
	Location			string	`json:"location"`
}
type songs []song

func GetSongsByArtist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistId := vars["artistId"]
	songList := buildSongsList(getSongs("where artistId = ?", artistId))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songList)
}

func GetAllSongs(w http.ResponseWriter, r *http.Request) {
	songList := buildSongsList(getSongs("", ""))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songList)
}

func GetSongById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	songId := vars["songId"]
	songList := buildSongsList(getSongs("where songs.id = ?", songId))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songList)
}

func getSongs(whereStatment string, value string) *sql.Rows{
	var rows *sql.Rows
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	selectStatement := "SELECT songs.id, name, people.firstName, people.lastName, length, location " +
	"FROM songs join people on artistId = people.id "
	if whereStatment != "" {
		stmt, err := db.Prepare(selectStatement + whereStatment)
		global.CheckErr(err)
		rows, err = stmt.Query(value)
		global.CheckErr(err)
	} else {
		rows, err = db.Query(selectStatement)
		global.CheckErr(err)
	}
	db.Close()
	return rows
}

func buildSongsList(rows *sql.Rows) songs {
	var sngList songs
	for rows.Next() {
		var id string
		var name string
		var firstName string
		var lastName string
		var length string
		var location string
		err := rows.Scan(&id, &name, &firstName, &lastName, &length, &location)
		global.CheckErr(err)
		sng := song{
			Id:              id,
			Name:            name,
			ArtistFirstName: firstName,
			ArtistLastName:  lastName,
			Length:          length,
			Location:        location,
		}
		sngList = append(sngList, sng)
	}
	return sngList
}