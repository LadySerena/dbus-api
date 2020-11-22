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
	"github.com/spf13/viper"
)

const (
	serviceKey = "SERVICE_NAME"
)

func main() {

	bindErr := viper.BindEnv(serviceKey)
	if bindErr != nil {
		log.Fatal("you must provide a service name")
	}

	serviceName := viper.GetString(serviceKey)
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
	config := server.NewConfig(client, serviceName)
	router := mux.NewRouter()

	router.Use(db.BasicAuthMiddleware)
	router.HandleFunc("/service", config.PostService).Methods(http.MethodPost)
	router.HandleFunc("/service", config.GetService).Methods(http.MethodGet)
	log.Print(http.ListenAndServe(":8080", router))
}
