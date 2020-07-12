package main

import (
	"dbus-api/pkg/dbus"
	"fmt"
)

func main() {
	client, clientCreateErr := dbus.NewClient()
	if clientCreateErr != nil {
		fmt.Printf("could not create the dbus client due to: %s\n", clientCreateErr.Error())
		return
	}
	defer client.Close()
	unitDetails, getErr := client.GetUnit("docker.service")
	if getErr != nil {
		fmt.Printf("could not get unit %s due to: %s\n", "docker.service", getErr.Error())
	}
	fmt.Printf("%#v\n", unitDetails)
}
