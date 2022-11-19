package device

import (
	"encoding/json"
	"net/http"
	"strings"

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
	var payload RequestPayload
	var items []openhab.Item
	_ = json.NewDecoder(r.Body).Decode(&payload)
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
	var request Request
	if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
		return
	}
	for _, item := range request.Payload.Devices {
		val := ""
		// Костыль для штор
		if strings.Contains(item.Id, "Roller") {
			if item.Capabilities[0].State.Value == true {
				val = "UP"
			} else {
				val = "DOWN"
			}
		} else {
			if item.Capabilities[0].State.Value == true {
				val = "ON"
			} else {
				val = "OFF"
			}
		}
		logger.Info(r, "Item "+item.Id+" set "+val)
		var status int
		if status, err = p.SetState(credentials, item.Id, val); err != nil {
			return
		}
		item.Capabilities[0].State.ActionResult = NewActionResult(status, err)
	}
	p.sendResponse(w, credentials, request.Payload.Devices)
	return
}

func (p Provider) sendResponse(w http.ResponseWriter, credentials string, devices []Device) {
	response := NewResponse(credentials, devices)
	js, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(js)
}

type Request struct {
	Payload RequestPayload `json:"payload"`
}

type RequestPayload struct {
	Devices []Device `json:"devices"`
}
