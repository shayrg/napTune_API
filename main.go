package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
	"github.com/shayrg/napTune_API/user"
)

func main(){
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").
		Path("/authenticate").
		Name("authenticate").
		Handler(http.HandlerFunc(user.Authenticate))
	router.Methods("POST").
		Path("/authorize").
		Name("authorize").
		HandlerFunc(http.HandlerFunc(user.Authorize))
	router.Methods("POST").
		Path("/logout").
		Name("logout").
		HandlerFunc(http.HandlerFunc(user.Logout))
	router.Methods("POST").
		Path("/new_user").
		Name("new_user").
		HandlerFunc(http.HandlerFunc(user.NewUser))
	//Listen and serve
	port := "8080"
	fmt.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}