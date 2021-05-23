package main

import (
	"dbus-api/internal/config"
	"log"
)

const (
	prefix               = "dbus_api"
	serviceKey           = "service_name"
	authFileKey          = "auth_file"
	listenAddressKey     = "listen_address"
	listenerDefault      = "localhost:8080"
	tlsEnabledKey        = "tls_enabled"
	tlsEnabledDefault    = false
	tlsCertPathKey       = "tls_cert_path"
	tlsPrivateKeyPathKey = "tls_key_path"
)

func main() {
	app, appErr := config.NewApp()
	if appErr != nil {
		log.Fatalln(appErr.Error())
	}
	log.Println(app.Run())
}
