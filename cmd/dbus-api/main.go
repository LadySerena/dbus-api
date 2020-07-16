package main

import (
	"dbus-api/pkg/dbus"
	"dbus-api/pkg/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	client, clientCreateErr := dbus.NewClient()
	if clientCreateErr != nil {
		fmt.Printf("could not create the dbus client due to: %s\n", clientCreateErr.Error())
		return
	}
	defer client.Close()
	config := server.NewConfig(client, "docker.service")
	http.HandleFunc("/service", config.GetService)
	http.HandleFunc("/meow/service", config.PostService)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
