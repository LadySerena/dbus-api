package main

import (
	"dbus-api/pkg/auth"
	"dbus-api/pkg/dbus"
	"dbus-api/pkg/server"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
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
	viper.SetEnvPrefix(prefix)
	viper.SetDefault(listenAddressKey, listenerDefault)
	viper.SetDefault(tlsEnabledKey, tlsEnabledDefault)
	viper.AutomaticEnv()
	serviceName := viper.GetString(serviceKey)
	if serviceName == "" {
		log.Fatalf("you must provide a service name via setting the environment variable: %s_%s", strings.ToUpper(prefix), strings.ToUpper(serviceKey))
		return
	}
	authFilePath := viper.GetString(authFileKey)
	if authFilePath == "" {
		log.Fatalf("you must provide a path to httpd formatted basic auth file via setting the environment variable: %s_%s", strings.ToUpper(prefix), strings.ToUpper(authFileKey))
		return
	}
	tlsEnabled := viper.GetBool(tlsEnabledKey)

	tlsCertPath := viper.GetString(tlsCertPathKey)
	if tlsCertPath == "" && tlsEnabled {
		log.Fatalf("you must provide the path to the tls cert (public key) via setting the environment variable %s_%s", strings.ToUpper(prefix), strings.ToUpper(tlsCertPathKey))
		return
	}

	tlsKeyPath := viper.GetString(tlsPrivateKeyPathKey)
	if tlsKeyPath == "" && tlsEnabled {
		log.Fatalf("you must provide the path to the tls private key via setting the environment variable %s_%s", strings.ToUpper(prefix), strings.ToUpper(tlsPrivateKeyPathKey))
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

	if tlsEnabled {
		log.Print(http.ListenAndServeTLS(listenerAddress, tlsCertPath, tlsKeyPath, router))
	} else {
		log.Print(http.ListenAndServe(listenerAddress, router))
	}
}
