package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"github.com/shayrg/napTune_API/people"
	"github.com/shayrg/napTune_API/songs"
	"github.com/shayrg/napTune_API/playlist"
)

func main(){
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").
		Path("/authenticate").
		Name("authenticate").
		Handler(http.HandlerFunc(people.Authenticate))
	router.Methods("POST").
		Path("/authorize").
		Name("authorize").
		HandlerFunc(http.HandlerFunc(people.Authorize))
	router.Methods("POST").
		Path("/logout").
		Name("logout").
		HandlerFunc(http.HandlerFunc(people.Logout))
	router.Methods("POST").
		Path("/new_user").
		Name("new_user").
		HandlerFunc(http.HandlerFunc(people.NewUser))
	router.Methods("GET").
		Path("/artists").
		Name("artists").
		HandlerFunc(http.HandlerFunc(people.GetAllArtists))
	router.Methods("GET").
		Path("/artists/{artistId}").
		Name("artistsById").
		HandlerFunc(http.HandlerFunc(people.GetArtistById))
	router.Methods("GET").
		Path("/artists/{artistId}/songs").
		Name("songsByArtistId").
		HandlerFunc(http.HandlerFunc(songs.GetSongsByArtist))
	router.Methods("GET").
		Path("/songs").
		Name("songs").
		HandlerFunc(http.HandlerFunc(songs.GetAllSongs))
	router.Methods("GET").
		Path("/songs/{songId}").
		Name("songsById").
		HandlerFunc(http.HandlerFunc(songs.GetSongById))
	router.Methods("GET").
		Path("/playlists").
		Name("playlists").
		HandlerFunc(http.HandlerFunc(playlist.GetAllPlaylists))
	router.Methods("GET").
		Path("/playlists/{playlistId}").
		Name("playlistById").
		HandlerFunc(http.HandlerFunc(playlist.GetPlaylistsById))
	router.Methods("GET").
		Path("/playlists/{playlistId}/songs").
		Name("getSongsById").
		HandlerFunc(http.HandlerFunc(playlist.GetPlaylistSongs))
	//Listen and serve
	port := "8080"
	fmt.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}