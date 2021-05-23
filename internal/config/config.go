package config

import (
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
		return nil, newConfigError(serviceName)
	}
	authFilePath := configViper.GetString(authFileKey)
	if authFilePath == "" {
		return nil, newConfigError(authFilePath)
	}

	listener := configViper.GetString(listenAddressKey)

	if listener == "" {
		listener = listenerDefault
	}
	tlsEnabled := configViper.GetBool(tlsEnabledKey)

	// only error if tls is enabled and the path is blank
	tlsCertPath := configViper.GetString(tlsCertPathKey)
	if tlsEnabled && tlsCertPath == "" {
		return nil, newTLSConfigError(tlsCertPathKey)
	}

	tlsPrivateKeyPath := configViper.GetString(tlsPrivateKeyPathKey)
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

type configError struct {
	message string
}

func newConfigError(key string) *configError {
	config := fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(key))
	return &configError{message: fmt.Sprintf("missing config: %s, configure this option via setting the environment variable: %s", config, config)}
}

func newTLSConfigError(key string) *configError {
	config := fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(key))
	return &configError{message: fmt.Sprintf("missing tls config: %s, configure this option via setting the environment variable: %s", config, config)}
}

func (e configError) Error() string {
	return e.message
}
