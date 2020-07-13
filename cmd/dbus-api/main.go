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
		return
	}
	fmt.Printf("%#v\n", unitDetails)

	startUnitErr := client.StartUnit("docker.service")
	if startUnitErr != nil {
		fmt.Printf("%s\n", startUnitErr.Error())
		return
	}
	stopUnitErr := client.StopUnit("docker.service")
	if stopUnitErr != nil {
		fmt.Printf("%s\n", stopUnitErr.Error())
		return
	}
	restartUnitErr := client.RestartUnit("docker.service")
	if restartUnitErr != nil {
		fmt.Printf("%s\n", restartUnitErr.Error())
		return
	}
	fmt.Printf("finished test\n")
}
