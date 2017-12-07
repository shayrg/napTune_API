package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"log"
)

func main(){
	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").
		Path("/authenticate").
		Name("authenticate").
		Handler(http.HandlerFunc(Authenticate))
	//Listen and serve
	port := "8080"
	fmt.Println("Listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}