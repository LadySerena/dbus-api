package main

import (
	"fmt"
	"github.com/godbus/dbus"
)

func main() {
	var list [][]interface{}

	conn, err := dbus.SystemBus()
	if err != nil {
		panic(err)
	}


	err = conn.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1").Call("org.freedesktop.systemd1.Manager.ListUnits", 0).Store(&list)
	if err != nil {
		panic(err)
	}
	for _, v := range list {
		fmt.Println(v)
	}
}