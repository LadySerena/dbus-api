package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	serviceName  = "test-service.service"
	authFilePath = "test.auth"
	tlsEnabled   = false
)

func Test_NewApp_BaseConfig(t *testing.T) {
	setupTestConfig()
	app, createErr := NewApp()
	assert.NoError(t, createErr)
	assert.Equal(t, serviceName, app.ServiceName)
	assert.Equal(t, authFilePath, app.AuthFilePath)
	assert.Equal(t, tlsEnabled, app.TLSConfig.TLSEnabled)
	assert.Equal(t, listenerDefault, app.ListenerAddress)
	assert.Equal(t, tlsEnabledDefault, app.TLSConfig.TLSEnabled)
	t.Cleanup(viper.Reset)
}

func Test_NewApp_CustomListerConfig(t *testing.T) {
	listener := "localhost:8081"
	setupTestConfig()
	viper.Set(listenAddressKey, listener)
	app, createErr := NewApp()
	assert.NoError(t, createErr)
	assert.Equal(t, serviceName, app.ServiceName)
	assert.Equal(t, authFilePath, app.AuthFilePath)
	assert.Equal(t, tlsEnabled, app.TLSConfig.TLSEnabled)
	assert.Equal(t, listener, app.ListenerAddress)
	t.Cleanup(viper.Reset)
}

func Test_NewApp_NoService(t *testing.T) {
	setupTestConfig()
	viper.Set(serviceKey, "") // override the service component
	app, createErr := NewApp()
	assert.Nil(t, app)
	assert.Error(t, createErr)
	assert.ErrorIs(t, createErr, newConfigError(serviceKey))
	t.Cleanup(viper.Reset)
}

func Test_NewApp_AuthFile(t *testing.T) {
	setupTestConfig()
	viper.Set(authFileKey, "") // override the auth file component
	app, createErr := NewApp()
	assert.Nil(t, app)
	assert.Error(t, createErr)
	assert.ErrorIs(t, createErr, newConfigError(authFileKey))
	t.Cleanup(viper.Reset)
}

func Test_NewApp_TLSHappy(t *testing.T) {
	certPath := "cert.crt"
	keyPath := "key.key"
	setupTestConfig()
	viper.Set(tlsEnabledKey, true)
	viper.Set(tlsCertPathKey, certPath)
	viper.Set(tlsPrivateKeyPathKey, keyPath)
	app, createErr := NewApp()
	assert.NoError(t, createErr)
	assert.Equal(t, serviceName, app.ServiceName)
	assert.Equal(t, authFilePath, app.AuthFilePath)
	assert.True(t, app.TLSConfig.TLSEnabled)
	assert.Equal(t, certPath, app.TLSConfig.CertPath)
	assert.Equal(t, keyPath, app.TLSConfig.KeyPath)
	t.Cleanup(viper.Reset)
}

func Test_NewApp_TLSMissingCertPath(t *testing.T) {
	keyPath := "key.key"
	setupTestConfig()
	viper.Set(tlsEnabledKey, true)
	viper.Set(tlsPrivateKeyPathKey, keyPath)
	app, createErr := NewApp()
	assert.Error(t, createErr)
	assert.ErrorIs(t, createErr, newTLSConfigError(tlsCertPathKey))
	assert.Nil(t, app)
	t.Cleanup(viper.Reset)
}

func Test_NewApp_TLSMissingKeyPath(t *testing.T) {
	certPath := "cert.crt"
	setupTestConfig()
	viper.Set(tlsEnabledKey, true)
	viper.Set(tlsCertPathKey, certPath)
	app, createErr := NewApp()
	assert.Error(t, createErr)
	assert.ErrorIs(t, createErr, newTLSConfigError(tlsPrivateKeyPathKey))
	assert.Nil(t, app)
	t.Cleanup(viper.Reset)
}

func Test_ConfigError(t *testing.T) {
	config := fmt.Sprintf("%s_%s", strings.ToUpper(prefix), strings.ToUpper(serviceKey))
	err := newConfigError(serviceKey)
	assert.Equal(t, err.Error(), fmt.Sprintf("missing config: %s, configure this option via setting the environment variable: %s", config, config))
}

func setupTestConfig() {
	viper.Set(serviceKey, serviceName)
	viper.Set(authFileKey, authFilePath)
	viper.Set(tlsEnabledKey, tlsEnabled)
}
