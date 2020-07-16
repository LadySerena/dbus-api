package dbus

import (
	"time"

	"github.com/godbus/dbus"
)

const (
	destination                       = "org.freedesktop.systemd1"
	systemdPath       dbus.ObjectPath = "/org/freedesktop/systemd1"
	mode                              = "replace"
	managerInterface                  = "org.freedesktop.systemd1.Manager"
	getUnitMethod                     = managerInterface + ".GetUnit"
	startUnitMethod                   = managerInterface + ".StartUnit"
	restartUnitMethod                 = managerInterface + ".RestartUnit"
	stopUnitMethod                    = managerInterface + ".StopUnit"
	StartService      Operation       = "start"
	RestartService    Operation       = "restart"
	StopService       Operation       = "stop"
)

type Operation string

type Client struct {
	connection *dbus.Conn
}

type UnitResponse struct {
	ServiceName string `json:"service-name"`
	Active      string `json:"active"`
	SubStatus   string `json:"sub-status"`
}

type ServiceChangeRequest struct {
	Operation Operation `json:"operation"`
}

func (c *Client) GetUnit(unitName string) (*UnitResponse, error) {
	time.Sleep(time.Millisecond * 500)
	var path dbus.ObjectPath

	getObjectErr := c.connection.Object(destination, systemdPath).Call(getUnitMethod, 0, unitName).Store(&path)
	if getObjectErr != nil {
		return nil, getObjectErr
	}

	serviceObject := c.connection.Object(destination, path)
	Id, IdErr := serviceObject.GetProperty("org.freedesktop.systemd1.Unit.Id")
	if IdErr != nil {
		return nil, IdErr
	}
	active, activeErr := serviceObject.GetProperty("org.freedesktop.systemd1.Unit.ActiveState")
	if activeErr != nil {
		return nil, activeErr
	}
	subState, subStateErr := serviceObject.GetProperty("org.freedesktop.systemd1.Unit.SubState")
	if subStateErr != nil {
		return nil, subStateErr
	}

	response := &UnitResponse{
		ServiceName: Id.Value().(string),
		Active:      active.Value().(string),
		SubStatus:   subState.Value().(string),
	}

	return response, nil
}

func (c *Client) StartUnit(unitName string) error {
	time.Sleep(time.Millisecond * 500)
	var jobPath dbus.ObjectPath
	startErr := c.connection.Object(destination, systemdPath).Call(startUnitMethod, 0, unitName, mode).Store(&jobPath)
	if startErr != nil {
		return startErr
	}
	return nil
}

func (c *Client) RestartUnit(unitName string) error {

	time.Sleep(time.Millisecond * 500)
	var jobPath dbus.ObjectPath
	restartErr := c.connection.Object(destination, systemdPath).Call(restartUnitMethod, 0, unitName, mode).Store(&jobPath)
	if restartErr != nil {
		return restartErr
	}
	return nil
}

func (c *Client) StopUnit(unitName string) error {

	time.Sleep(time.Millisecond * 500)
	var jobPath dbus.ObjectPath
	stopErr := c.connection.Object(destination, systemdPath).Call(stopUnitMethod, 0, unitName, mode).Store(&jobPath)
	if stopErr != nil {
		return stopErr
	}
	return nil
}

func (c *Client) Close() error {
	return c.connection.Close()
}

func NewClient() (*Client, error) {
	conn, connectionErr := dbus.SystemBus()
	if connectionErr != nil {
		return nil, connectionErr
	}
	return &Client{connection: conn}, nil
}


