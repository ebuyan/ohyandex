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

func (m Mapper) getType(item openhab.Item) (itemType string) {
	switch item.Type {
	case "Switch":
		itemType = "devices.types.switch"
	case "Light", "Dimmer", "Color":
		itemType = "devices.types.light"
	}
	return
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
		capabilities = append(capabilities, m.getSwitchCapability(item))
		return
	//todo
	case "Dimmer":
	case "Color":
	}
	return
}

func (m Mapper) getSwitchCapability(item openhab.Item) Capability {
	return Capability{
		Type:        "devices.capabilities.on_off",
		Retrievable: true,
		State: State{
			Instance: "on",
			Value:    item.State == "ON",
		},
	}
}
