package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Run() {
	router := mux.NewRouter()
	router.HandleFunc("/create", CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/get/{id}", SearchUser).Methods(http.MethodGet)
	router.HandleFunc("/get", ShowUsers).Methods(http.MethodGet)
	router.HandleFunc("/alter/{id}", AlterUser).Methods(http.MethodPut)
	router.HandleFunc("/delete/{id}", DeleteUser).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":5000", router))
}
