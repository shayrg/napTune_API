package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"github.com/shayrg/napTune_API/people"
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
	//Listen and serve
	port := "8080"
	fmt.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}