package dbus

import (
	"github.com/godbus/dbus"
)

type Client struct {
	connection *dbus.Conn
}

type UnitResponse struct {
	ServiceName string `json:"service-name"`
	Active      string `json:"active"`
	SubStatus   string `json:"sub-status"`
}

func (c *Client) GetUnit(unitName string) (*UnitResponse, error) {
	var path dbus.ObjectPath

	connection, connectionErr := dbus.SystemBus()
	if connectionErr != nil {
		return nil, connectionErr
	}

	getObjectErr := connection.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1").Call("org.freedesktop.systemd1.Manager.GetUnit", 0, unitName).Store(&path)
	if getObjectErr != nil {
		return nil, getObjectErr
	}

	serviceObject := connection.Object("org.freedesktop.systemd1", path)
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
