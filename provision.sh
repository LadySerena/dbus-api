#!/usr/bin/env bash

# install go shit
wget https://golang.org/dl/go1.15.5.linux-amd64.tar.gz -P /tmp
tar -C /usr/local -xzf /tmp/go1.15.5.linux-amd64.tar.gz
echo "PATH=$PATH:/usr/local/go/bin" >>/etc/profile
# polkit nonsense (gotta keep it old school cuz buster ships with policykit 105.x instead of polkit :/ )
apt-get update
apt-get install -y software-properties-common prometheus-node-exporter openjdk-11-jdk #has policy kit in it and some dummy service to manage


