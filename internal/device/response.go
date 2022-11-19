package device

import "github.com/google/uuid"

func NewResponse(userId string, devices []Device) Response {
	return Response{
		Id: uuid.NewString(),
		Payload: Payload{
			UserId:  userId,
			Devices: devices,
		},
	}
}

type Response struct {
	Id      string  `json:"request_id"`
	Payload Payload `json:"payload"`
}

type Payload struct {
	UserId  string   `json:"user_id"`
	Devices []Device `json:"devices"`
}

type Device struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Room         string       `json:"room"`
	Type         string       `json:"type"`
	Capabilities []Capability `json:"capabilities"`
}

type Capability struct {
	Type        string        `json:"type"`
	Retrievable bool          `json:"retrievable"`
	Parameters  []interface{} `json:"parameters"`
	State       State         `json:"state"`
}

type Parameter struct {
	Instance string `json:"instance"`
	Unit     string `json:"unit"`
	Range    Range  `json:"range"`
}

type Range struct {
	Min       int `json:"min"`
	Max       int `json:"max"`
	Precision int `json:"precision"`
}

type State struct {
	Instance     string       `json:"instance"`
	Value        interface{}  `json:"value"`
	ActionResult ActionResult `json:"action_result"`
}
