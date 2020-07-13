# dbus-api

REST API to interface with systemd on servers. The `resources/etc/polkit-1/rules.d/10-dbus-api.rules` contains a sample
polkit rule that will allow the api to manage the service listed in the `action.lookup("unit") == <your service name here>`. 

TODO add usage information

TODO poll service status to add "rate limiting"