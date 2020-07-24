package main

import (
	"dbus-api/pkg/auth"
	"dbus-api/pkg/dbus"
	"dbus-api/pkg/server"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	client, clientCreateErr := dbus.NewClient()
	if clientCreateErr != nil {
		fmt.Printf("could not create the dbus client due to: %s\n", clientCreateErr.Error())
		return
	}
	defer client.Close()
	fmt.Println(os.Getwd())
	db, dbErr := auth.NewDatabase("./secrets/decrypted/authUsers")
	if dbErr != nil {
		fmt.Printf("could not create database due to: %s\n", dbErr.Error())
		return
	}
	config := server.NewConfig(client, "docker.service")
	router := mux.NewRouter()

	router.Use(db.BasicAuthMiddleware)
	router.HandleFunc("/service", config.PostService).Methods(http.MethodPost)
	router.HandleFunc("/service", config.GetService).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", router))
}
