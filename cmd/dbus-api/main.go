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
	serviceKey       = "SERVICE_NAME"
	authFileKey      = "AUTH_FILE"
	listenAddressKey = "LISTEN_ADDRESS"
	listenerDefault  = "localhost:8080"
)

func main() {
	viper.SetDefault(listenAddressKey, listenerDefault)
	bindErr := viper.BindEnv(serviceKey)
	if bindErr != nil {
		log.Fatalf("error getting service name from environment for %s: %s", serviceKey, bindErr.Error())
		return
	}
	serviceName := viper.GetString(serviceKey)
	if serviceName == "" {
		log.Fatalf("you must provide a service name via setting the environment variable: %s", serviceKey)
		return
	}
	authBindErr := viper.BindEnv(authFileKey)
	if authBindErr != nil {
		log.Fatalf("error getting auth file path from environment for %s: %s", authFileKey, authBindErr.Error())
		return
	}
	authFilePath := viper.GetString(authFileKey)
	if authFilePath == "" {
		log.Fatalf("you must provide a path to httpd formatted basic auth file via setting the environment variable: %s", authFileKey)
		return
	}
	listenBindErr := viper.BindEnv(listenAddressKey)
	if listenBindErr != nil {
		log.Printf("error getting listener address from environment for %s: %s", listenAddressKey, listenBindErr.Error())
		log.Printf("using default value of %s for %s", listenerDefault, listenAddressKey)
		return
	}
	listenerAddress := viper.GetString(listenAddressKey)
	if listenerAddress == listenerDefault {
		log.Printf("using default lister address of %s", listenerDefault)
	}

	client, clientCreateErr := dbus.NewClient()
	if clientCreateErr != nil {
		fmt.Printf("could not create the dbus client due to: %s\n", clientCreateErr.Error())
		return
	}
	defer client.Close()
	fmt.Println(os.Getwd())
	db, dbErr := auth.NewDatabase(authFilePath)
	if dbErr != nil {
		fmt.Printf("could not create database due to: %s\n", dbErr.Error())
		return
	}
	config := server.NewConfig(client, serviceName)
	router := mux.NewRouter()

	router.Use(db.BasicAuthMiddleware)
	router.HandleFunc("/service", config.PostService).Methods(http.MethodPost)
	router.HandleFunc("/service", config.GetService).Methods(http.MethodGet)
	log.Print(http.ListenAndServe(listenerDefault, router))
}
