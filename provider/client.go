package provider

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

type OpenhabClient struct{ baseUrl string }

func NewOpenhabClient() OpenhabClient {
	return OpenhabClient{baseUrl: os.Getenv("OPENHAB_HOST")}
}

func (c OpenhabClient) Ping(cregentials string) (bool, error) {
	_, status, err := c.sendRequest(cregentials, c.baseUrl+"/rest/uuid", "GET", nil)
	return status == http.StatusOK, err
}

func (c OpenhabClient) GetItem(cregentials string, itemName string) (item Item, status int, err error) {
	body, status, err := c.sendRequest(cregentials, c.baseUrl+"/rest/items/"+itemName, "GET", nil)
	json.Unmarshal(body, &item)
	return
}

func (c OpenhabClient) GetItems(cregentials string) (items []Item, status int, err error) {
	body, status, err := c.sendRequest(cregentials, c.baseUrl+"/rest/items?tags=yandex", "GET", nil)
	json.Unmarshal(body, &items)
	return
}

func (c OpenhabClient) GetRooms(cregentials string) (items []Item, status int, err error) {
	body, status, err := c.sendRequest(cregentials, c.baseUrl+"/rest/items?tags=Room", "GET", nil)
	json.Unmarshal(body, &items)
	return
}

func (c OpenhabClient) SetState(cregentials string, itemName string, value string) (status int, err error) {
	_, status, err = c.sendRequest(cregentials, c.baseUrl+"/rest/items/"+itemName, "POST", []byte(value))
	return
}

func (c OpenhabClient) sendRequest(cregentials string, url string, methodType string, reqBody []byte) (body []byte, status int, err error) {
	req, err := http.NewRequest(methodType, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Basic "+cregentials)
	if methodType == http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "text/plain")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	status = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

type Item struct {
	Name   string   `json:"name"`
	Label  string   `json:"label"`
	Link   string   `json:"link"`
	State  string   `json:"state"`
	Type   string   `json:"type"`
	Groups []string `json:"groupNames"`
	Tags   []string `json:"tags"`
}
