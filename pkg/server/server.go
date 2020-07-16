package server

import (
	"dbus-api/pkg/dbus"
	"encoding/json"
	"fmt"
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
		_, _ = w.Write(generateErrorResponse(http.StatusInternalServerError, "could not get unit status", getUnitErr))
		return
	}

	responseBytes, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		_, _ = w.Write(generateErrorResponse(http.StatusInternalServerError, "could not marshal unit status", marshalErr))
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
		_, _ = w.Write(generateErrorResponse(http.StatusInternalServerError, "could not read request body", readErr))
		return
	}

	serviceRequest := dbus.ServiceChangeRequest{}

	unmarshalErr := json.Unmarshal(bodyBytes, &serviceRequest)
	if unmarshalErr != nil {
		_, _ = w.Write(generateErrorResponse(http.StatusInternalServerError, "could not unmarshal request", unmarshalErr))
		return
	}

	switch serviceRequest.Operation {
	case dbus.StartService:
		startErr := c.client.StartUnit(c.unitName)
		if startErr != nil {
			_, _ = w.Write(generateErrorResponse(http.StatusInternalServerError, "could not start unit: "+c.unitName, startErr))
			return
		}
		fallthrough
	default:
		response, getUnitErr := c.client.GetUnit(c.unitName)
		if getUnitErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(generateErrorResponse(http.StatusInternalServerError, "could not get unit status", getUnitErr))
			return
		}

		responseBytes, marshalErr := json.Marshal(response)
		if marshalErr != nil {
			_, _ = w.Write(generateErrorResponse(http.StatusInternalServerError, "could not marshal unit status", marshalErr))
			return
		}
		_, _ = w.Write(responseBytes)
		return
	}
}

func NewConfig(client *dbus.Client, unitName string) Config {
	return Config{
		client:   client,
		unitName: unitName,
		mux:      sync.Mutex{},
	}
}

func generateErrorResponse(code int32, text string, err error) []byte {
	return []byte(fmt.Sprintf(`{"status": %d, "reason":"%s due to: %s"}`, code, text, err.Error()))
}
