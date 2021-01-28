# dbus-api

REST API to interface with systemd on servers. The `resources/etc/polkit-1/rules.d/10-dbus-api.rules` contains a sample
polkit rule that will allow the api to manage the service listed in
the `action.lookup("unit") == <your service name here>`.

TODO add usage information

TODO poll service status to add "rate limiting"

# configuration

The api uses [viper](https://github.com/spf13/viper) and is configured via environment variables

- DBUS_API_SERVICE_NAME
    - the name of the systemd service that you want to manage. Note this must include the `.service`
- DBUS_API_AUTH_FILE
    - path to the httpd formatted basic auth file
- DBUS_API_LISTEN_ADDRESS
    - which ip address and port to listen on defaults to localhost:8080
- DBUS_API_TLS_ENABLED
    - true to enable tls or false to disable it (when deploying this in a networked setting please use tls)
- DBUS_API_TLS_CERT_PATH
    - required if `DBUS_API_TLS_ENABLED=true` path to the public key of the certificate
- DBUS_API_TLS_KEY_PATH
    - required if `DBUS_API_TLS_ENABLED=true` path to the private key of the certificate