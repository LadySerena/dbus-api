package server

import (
	"dbus-api/pkg/common"
	"dbus-api/pkg/dbus"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
)

const (
	jsonContent          = "application/json"
	contentTypeHeaderKey = "Content-Type"
)

type APIError struct {
	Status int32  `json:"status"`
	Reason string `json:"reason"`
}

type Config struct {
	client   *dbus.Client
	unitName string
	mux      sync.Mutex
}

func (c *Config) GetService(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(contentTypeHeaderKey, jsonContent)
	defer r.Body.Close()
	c.mux.Lock()
	defer c.mux.Unlock()
	response, getUnitErr := c.client.GetUnit(c.unitName)
	if getUnitErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not get unit status", getUnitErr))
		return
	}

	responseBytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not marshal unit status", marshalErr))
		return
	}
	_, _ = w.Write(responseBytes)
}

func (c *Config) PostService(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(contentTypeHeaderKey, jsonContent)
	defer r.Body.Close()
	bodyBytes, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not read request body", readErr))
		return
	}

	serviceRequest := dbus.ServiceChangeRequest{}

	unmarshalErr := json.Unmarshal(bodyBytes, &serviceRequest)
	if unmarshalErr != nil {
		_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not unmarshal request", unmarshalErr))
		return
	}

	switch serviceRequest.Operation {
	case dbus.StartService:
		startErr := c.client.StartUnit(c.unitName)
		if startErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not start unit: "+c.unitName, startErr))
			return
		}
	case dbus.StopService:
		stopErr := c.client.StopUnit(c.unitName)
		if stopErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not stop unit: "+c.unitName, stopErr))
			return
		}
	case dbus.RestartService:
		restartErr := c.client.RestartUnit(c.unitName)
		if restartErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not restart unit: "+c.unitName, restartErr))
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(common.GenerateProblemResponse(http.StatusBadRequest, "invalid operation, please use "+string(dbus.StartService)+", "+string(dbus.RestartService)+", or "+string(dbus.StopService)))
		return
	}
	response, getUnitErr := c.client.GetUnit(c.unitName)
	if getUnitErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not get unit status", getUnitErr))
		return
	}

	responseBytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		_, _ = w.Write(common.GenerateErrorResponse(http.StatusInternalServerError, "could not marshal unit status", marshalErr))
		return
	}
	_, _ = w.Write(responseBytes)
}

func NewConfig(client *dbus.Client, unitName string) Config {
	return Config{
		client:   client,
		unitName: unitName,
		mux:      sync.Mutex{},
	}
}
