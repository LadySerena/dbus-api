package main

import (
	"dbus-api/internal/config"
	"log"
)

func main() {
	app, appErr := config.NewApp()
	if appErr != nil {
		log.Fatalln(appErr.Error())
	}
	log.Println(app.Run())
}
