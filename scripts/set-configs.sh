#!/usr/bin/env bash

#gcloud get secrets ....

dbus_password_file=$(gcloud secrets versions access "latest" --secret="mc-dbus-api-htpasswd")
echo -n "$dbus_password_file" >/etc/dbus-api/auth
cat <<'EOF' >>/etc/dbus-api/environment
DBUS_API_SERVICE_NAME=minecraft.service
DBUS_API_AUTH_FILE=/etc/dbus-api/auth
;need to figure out internal ip
;maybe this? https://cloud.google.com/compute/docs/instances/view-ip-address#api
;DBUS_API_LISTEN_ADDRESS=:8080
DBUS_API_TLS_ENABLED=true
DBUS_API_TLS_CERT_PATH=/etc/dbus-api/cert.pem
DBUS_API_TLS_KEY_PATH=/etc/dbus-api/key.pem
EOF
