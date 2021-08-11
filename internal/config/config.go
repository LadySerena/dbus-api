package config

import (
	"dbus-api/pkg/auth"
	"dbus-api/pkg/dbus"
	"dbus-api/pkg/server"
	"fmt"
	"net/http"
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

type App struct {
	ServiceName     string
	AuthFilePath    string
	TLSConfig       TLS
	ListenerAddress string
}

type TLS struct {
	TLSEnabled bool
	CertPath   string
	KeyPath    string
}

func NewApp() (*App, error) {
	viper.SetEnvPrefix(prefix)
	viper.SetDefault(listenAddressKey, listenerDefault)
	viper.SetDefault(tlsEnabledKey, tlsEnabledDefault)
	viper.AutomaticEnv()

	serviceName := viper.GetString(serviceKey)
	if serviceName == "" {
		return nil, newConfigError(serviceKey)
	}
	authFilePath := viper.GetString(authFileKey)
	if authFilePath == "" {
		return nil, newConfigError(authFileKey)
	}

	listener := viper.GetString(listenAddressKey)

	tlsEnabled := viper.GetBool(tlsEnabledKey)

	// only error if tls is enabled and the path is blank
	tlsCertPath := viper.GetString(tlsCertPathKey)
	if tlsEnabled && tlsCertPath == "" {
		return nil, newTLSConfigError(tlsCertPathKey)
	}

	tlsPrivateKeyPath := viper.GetString(tlsPrivateKeyPathKey)
	// only error if tls is enabled and the path is blank
	if tlsEnabled && tlsPrivateKeyPath == "" {
		return nil, newTLSConfigError(tlsPrivateKeyPathKey)
	}

	return &App{
		ServiceName:  serviceName,
		AuthFilePath: authFilePath,
		TLSConfig: TLS{
			TLSEnabled: tlsEnabled,
			CertPath:   tlsCertPath,
			KeyPath:    tlsPrivateKeyPath,
		},
		ListenerAddress: listener,
	}, nil
}

func (a App) Run() error {
	client, clientCreateErr := dbus.NewClient()
	if clientCreateErr != nil {
		return clientCreateErr
	}
	defer client.Close()
	db, dbErr := auth.NewDatabase(a.AuthFilePath)
	if dbErr != nil {
		return dbErr
	}
	appServer := server.NewConfig(client, a.ServiceName)
	router := mux.NewRouter()

	router.Use(db.BasicAuthMiddleware)
	router.HandleFunc("/service", appServer.PostService).Methods(http.MethodPost)
	router.HandleFunc("/service", appServer.GetService).Methods(http.MethodGet)

	if a.TLSConfig.TLSEnabled {
		return http.ListenAndServeTLS(a.ListenerAddress, a.TLSConfig.CertPath, a.TLSConfig.KeyPath, router)
	} else {
		return http.ListenAndServe(a.ListenerAddress, router)
	}
}

type configError struct {
	message string
}

func newConfigError(key string) configError {
	config := fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(key))
	return configError{message: fmt.Sprintf("missing config: %s, configure this option via setting the environment variable: %s", config, config)}
}

func newTLSConfigError(key string) configError {
	config := fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(key))
	return configError{message: fmt.Sprintf("missing tls config: %s, configure this option via setting the environment variable: %s", config, config)}
}

func (e configError) Error() string {
	return e.message
}
