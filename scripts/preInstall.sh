#!/usr/bin/env bash

getent group dbus-api >/dev/null || groupadd -r dbus-api
getent passwd dbus-api >/dev/null || \
    useradd -r -g dbus-api -d /dbus-api -s /sbin/nologin \
    -c "service account for dbus-api" dbus-api