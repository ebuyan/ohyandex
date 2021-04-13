package provider

import (
	"encoding/json"
	"net/http"

	"ohyandex/logger"
)

type ServiceProvider struct {
	Mapper
	OpenhabClient
}

func NewServiceProvider() ServiceProvider {
	return ServiceProvider{NewMapper(), NewOpenhabClient()}
}

func (p ServiceProvider) Devices(w http.ResponseWriter, r *http.Request, creadentials string) (err error) {
	items, _, err := p.OpenhabClient.GetItems(creadentials)
	if err != nil {
		return
	}
	rooms, _, err := p.OpenhabClient.GetRooms(creadentials)
	if err != nil {
		return
	}
	devices, err := p.Mapper.MapOpenhabItemToYandexDevice(items, rooms)
	if err != nil {
		return
	}
	logger.Info(r, "Get devices query")
	p.sendResponse(w, creadentials, devices)
	return
}

func (p ServiceProvider) DevicesState(w http.ResponseWriter, r *http.Request, creadentials string) (err error) {
	var payload DevicePayload
	var items []Item
	json.NewDecoder(r.Body).Decode(&payload)
	for _, rItem := range payload.Devices {
		item, _, err := p.OpenhabClient.GetItem(creadentials, rItem.Id)
		if err != nil {
			logger.Error(r, err)
			continue
		}
		logger.Info(r, "Item "+rItem.Id+" get State value "+item.State)
		items = append(items, item)
	}
	if err != nil {
		return
	}
	devices, err := p.Mapper.MapOpenhabItemToYandexDevice(items, nil)
	if err != nil {
		return
	}
	p.sendResponse(w, creadentials, devices)
	return
}

func (p ServiceProvider) ControlDevices(w http.ResponseWriter, r *http.Request, creadentials string) (err error) {
	var request DeviceRequest
	json.NewDecoder(r.Body).Decode(&request)
	for _, rItem := range request.Payload.Devices {
		state := &rItem.Capabilities[0].State
		val := state.getValue()
		logger.Info(r, "Item "+rItem.Id+" set "+val)
		status, err := p.OpenhabClient.SetState(creadentials, rItem.Id, val)
		state.ActionResult = NewActionResult(status, err)
	}
	p.sendResponse(w, creadentials, request.Payload.Devices)
	return
}

func (p ServiceProvider) sendResponse(w http.ResponseWriter, creadentials string, devices []Device) {
	response := NewResponse(creadentials, devices)
	js, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type DeviceRequest struct {
	Payload DevicePayload `json:"payload"`
}

type DevicePayload struct {
	Devices []Device `json:"devices"`
}
