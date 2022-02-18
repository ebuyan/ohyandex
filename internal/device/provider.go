package device

import (
	"encoding/json"
	"net/http"

	"github.com/ebuyan/ohyandex/pkg/logger"
	"github.com/ebuyan/ohyandex/pkg/openhab"
)

type Provider struct {
	Mapper
	openhab.Client
}

func NewProvider() Provider {
	return Provider{NewMapper(), openhab.NewClient()}
}

func (p Provider) Devices(w http.ResponseWriter, r *http.Request, credentials string) (err error) {
	items, _, err := p.GetAllItemsByTag(credentials, "yandex")
	if err != nil {
		return
	}
	rooms, _, err := p.GetRooms(credentials)
	if err != nil {
		return
	}
	devices, err := p.MapOpenhabItemToYandexDevice(items, rooms)
	if err != nil {
		return
	}
	logger.Info(r, "Get devices query")
	p.sendResponse(w, credentials, devices)
	return
}

func (p Provider) DevicesState(w http.ResponseWriter, r *http.Request, credentials string) (err error) {
	var payload DevicePayload
	var items []openhab.Item
	json.NewDecoder(r.Body).Decode(&payload)
	for _, rItem := range payload.Devices {
		item, _, err := p.GetItem(credentials, rItem.Id)
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
	p.sendResponse(w, credentials, devices)
	return
}

func (p Provider) ControlDevices(w http.ResponseWriter, r *http.Request, credentials string) (err error) {
	var request DeviceRequest
	json.NewDecoder(r.Body).Decode(&request)
	for _, rItem := range request.Payload.Devices {
		state := &rItem.Capabilities[0].State
		val := state.getValue()
		logger.Info(r, "Item "+rItem.Id+" set "+val)
		status, err := p.SetState(credentials, rItem.Id, val)
		state.ActionResult = NewActionResult(status, err)
	}
	p.sendResponse(w, credentials, request.Payload.Devices)
	return
}

func (p Provider) sendResponse(w http.ResponseWriter, credentials string, devices []Device) {
	response := NewResponse(credentials, devices)
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
