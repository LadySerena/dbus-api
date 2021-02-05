#!/usr/bin/env bash
auth_file_path="/etc/dbus-api/auth"
tls_key_path="/etc/dbus-api/key.pem"
tls_cert_path="/etc/dbus-api/cert.pem"
dbus_auth_file_secret_name=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/dbus-secret-name -H "Metadata-Flavor: Google")
api_managed_service=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/service-name -H "Metadata-Flavor: Google")
internal_ip=$(curl http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip -H "Metadata-Flavor: Google")
tls_secret_name=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/tls-secret-name -H "Metadata-Flavor: Google")
tls_cert=$(curl http://metadata.google.internal/computeMetadata/v1/instance/attributes/tls-cert -H "Metadata-Flavor: Google")
dbus_password_file=$(gcloud secrets versions access "latest" --secret="$dbus_auth_file_secret_name")
tls_key=$(gcloud secrets versions access "latest" --secret="$tls_secret_name")
echo -n "$tls_key" >$tls_key_path
echo -n "$tls_cert" >$tls_cert_path
echo -n "$dbus_password_file" >$auth_file_path

cat <<EOF >/etc/dbus-api/environment
DBUS_API_SERVICE_NAME=$api_managed_service
DBUS_API_AUTH_FILE=$auth_file_path
DBUS_API_LISTEN_ADDRESS=$internal_ip:8080
DBUS_API_TLS_ENABLED=true
DBUS_API_TLS_CERT_PATH=$tls_cert_path
DBUS_API_TLS_KEY_PATH=$tls_key_path
EOF
