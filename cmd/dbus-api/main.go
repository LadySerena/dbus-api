package main

import (
	"log"

	"github.com/LadySerena/dbus-api/internal/config"
)

func main() {
	app, appErr := config.NewApp()
	if appErr != nil {
		log.Fatalln(appErr.Error())
	}
	log.Println(app.Run())
}
