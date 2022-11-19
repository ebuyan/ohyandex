package device

import "github.com/ebuyan/ohyandex/pkg/openhab"

type Mapper struct{}

func NewMapper() Mapper {
	return Mapper{}
}

func (m Mapper) MapOpenhabItemToYandexDevice(items []openhab.Item, rooms []openhab.Item) (devices []Device, err error) {
	for _, item := range items {
		m.fixType(&item)
		device := Device{
			Id:           item.Name,
			Name:         item.Label,
			Room:         m.getRoom(item, rooms),
			Type:         m.getType(item),
			Capabilities: m.getCapabilities(item),
		}
		devices = append(devices, device)
	}
	return
}

func (m Mapper) getRoom(item openhab.Item, rooms []openhab.Item) (roomName string) {
	roomName = "Дом"
	if rooms == nil {
		return
	}
	for _, room := range rooms {
		for _, group := range item.Groups {
			if room.Name == group {
				return room.Label
			}
		}
	}
	return
}

func (m Mapper) getType(item openhab.Item) string {
	switch item.Type {
	case "Switch":
		return TypeSwitch
	case "Light", "Dimmer", "Color":
		return TypeLight
	case "Rollershutter":
		return TypeCurtain
	default:
		return ""
	}
}

func (m Mapper) fixType(item *openhab.Item) {
	for _, tag := range item.Tags {
		if tag == "Light" {
			item.Type = tag
			return
		}
	}
}

func (m Mapper) getCapabilities(item openhab.Item) (capabilities []Capability) {
	switch item.Type {
	case "Switch", "Light":
		return []Capability{{
			Type:        CapabilitiesOnOff,
			Retrievable: true,
			State: State{
				Instance: "on",
				Value:    item.State == "ON",
			},
		}}
	case "Rollershutter":
		return []Capability{{
			Type:        CapabilitiesOnOff,
			Retrievable: true,
			State: State{
				Instance: "on",
				Value:    item.State == "0",
			},
		}}
	case "Dimmer":
	case "Color":
	}
	return
}
