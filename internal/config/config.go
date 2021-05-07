package config

import (
	"errors"
	"fmt"
	"strings"

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
	// initialize new viper and set defaults
	configViper := viper.New()
	configViper.SetEnvPrefix(prefix)
	configViper.SetDefault(listenAddressKey, listenerDefault)
	configViper.SetDefault(tlsEnabledKey, tlsEnabledDefault)
	configViper.AutomaticEnv()

	serviceName := configViper.GetString(serviceKey)
	if serviceName == "" {
		return nil, errors.New(fmt.Sprintf("you must provide a service name via setting the environment variable: %s_%s", strings.ToUpper(prefix), strings.ToUpper(serviceKey)))
	}

	return nil, errors.New("not implemented yet")
}
