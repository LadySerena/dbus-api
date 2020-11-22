#!/usr/bin/env bash

# install go shit
wget https://golang.org/dl/go1.15.5.linux-amd64.tar.gz -P /tmp
tar -C /usr/local -xzf /tmp/go1.15.5.linux-amd64.tar.gz
echo "PATH=$PATH:/usr/local/go/bin" >>/etc/profile
# polkit nonsense (gotta keep it old school cuz buster ships with policykit 105.x instead of polkit :/ )
apt-get update
apt-get install -y software-properties-common prometheus-node-exporter #has policy kit in it and some dummy service to manage


#mkdir -p /etc/polkit-1/localauthority/50-local.d/
cp /vagrant/resources/etc/polkit-1/localauthority/com.serenacodes.dbus.api.pkla /etc/polkit-1/localauthority/50-local.d/
