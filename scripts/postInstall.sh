#!/usr/bin/env bash

systemctl enable dbus-api.service
mkdir -p /dbus-api
mkdir -p /etc/dbus-api
chown dbus-api:dbus-api /dbus-api
chown dbus-api:dbus-api /etc/dbus-api
chmod 0751 /etc/dbus-api
chmod 0751 /dbus-api
