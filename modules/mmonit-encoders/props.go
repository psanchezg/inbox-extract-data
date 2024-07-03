package mmonitencoders

import (
	"github.com/psanchezg/inbox-extract-data/interfaces"
)

var Props = map[string]interfaces.Prop{
	"from": {
		Text: [...]string{"Uso de %v a %v\n", ""},
		Unit: "",
		Key:  "start",
	},
	"to": {
		Text: [...]string{"******* Plan activo **********\n", ""},
		Unit: "",
		Key:  "stop",
	},
	"duration": {
		Text: [...]string{"Duración del evento: %v %v\n", ""},
		Unit: "minutos",
		Key:  "minutes",
	},
	"channel": {
		Text: [...]string{"Canal: %v %v\n", ""},
		Unit: "",
		Key:  "channel",
	},
	"total": {
		Text: [...]string{"", ""},
		Unit: "segundos",
		Key:  "duration",
	},
}
