package playlist

import (
	"net/http"
	"database/sql"
	"github.com/shayrg/napTune_API/global"
	"encoding/json"
	"github.com/gorilla/mux"
)

type playlist struct {
	Id		string	`json:"id"`
	Name 	string	`json:"name"`
}
type playlistSong struct {
	PlaylistId			string 	`json:"playlist_id"`
	Id 					string 	`json:"id"`
	Name				string 	`json:"name"`
	ArtistFirstName 	string 	`json:"first_name"`
	ArtistLastName 		string 	`json:"last_name"`
	Length 				string 	`json:"length"`
	Location			string	`json:"location"`
	Position 			string	`json:"position"`
}
type playlists []playlist
type playlistSongs []playlistSong

func GetAllPlaylists(w http.ResponseWriter, r *http.Request) {
	plylists := buildPlaylistList(getPlaylists("", ""))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plylists)
}

func GetPlaylistsById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistId := vars["playlistId"]
	plylists := buildPlaylistList(getPlaylists("where id = ?", playlistId))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(plylists)
}

func GetPlaylistSongs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playlistId := vars["playlistId"]
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	stmt, err := db.Prepare("SELECT playlistId, songs.id, name, people.firstName, people.lastName, length, location, " +
		"position FROM songs join people on artistId = people.id join playlistSongs on songId = songs.id " +
		"where playlistId = ?")
	global.CheckErr(err)
	rows, err := stmt.Query(playlistId)
	global.CheckErr(err)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(buildPlaylistSongs(rows))
}

func getPlaylists(whereStatement string, value string) *sql.Rows{
	var rows *sql.Rows
	db, err := sql.Open("mysql", global.DbString)
	global.CheckErr(err)
	selectStatement := "SELECT id, name from playlists "
	if whereStatement != "" {
		stmt, err := db.Prepare(selectStatement + whereStatement)
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

func buildPlaylistList(rows *sql.Rows) playlists {
	var plylists playlists
	for rows.Next() {
		var id string
		var name string
		err := rows.Scan(&id, &name)
		global.CheckErr(err)
		plylist := playlist{
			Id:              id,
			Name:            name,
		}
		plylists = append(plylists, plylist)
	}
	return plylists
}

func buildPlaylistSongs(rows *sql.Rows) playlistSongs {
	var songList playlistSongs
	for rows.Next() {
		var playListId string
		var id string
		var name string
		var firstName string
		var lastName string
		var length string
		var location string
		var position string
		err := rows.Scan(&playListId, &id, &name, &firstName, &lastName, &length, &location, &position)
		global.CheckErr(err)
		sng := playlistSong{
			PlaylistId:		 playListId,
			Id:              id,
			Name:            name,
			ArtistFirstName: firstName,
			ArtistLastName:  lastName,
			Length:          length,
			Location:        location,
			Position:		 position,
		}
		songList = append(songList, sng)
	}
	return songList
}