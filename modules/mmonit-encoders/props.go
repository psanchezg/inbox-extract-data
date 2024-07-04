package mmonitencoders

import (
	"github.com/psanchezg/inbox-extract-data/interfaces"
)

var Props = map[string]interfaces.Prop{
	"from": {
		Text:   [...]string{"Uso de %v a %v\n", ""},
		Unit:   "",
		Key:    "start",
		Column: [...]string{"Inicio", ""},
	},
	"to": {
		Text:   [...]string{"******* Plan activo **********\n", ""},
		Unit:   "",
		Key:    "stop",
		Column: [...]string{"Fin", ""},
	},
	"duration": {
		Text:   [...]string{"Duración del evento: %v %v\n", ""},
		Unit:   "minutos",
		Key:    "minutes",
		Column: [...]string{"Duración (mins)", ""},
	},
	"channel": {
		Text:   [...]string{"Canal: %v %v\n", ""},
		Unit:   "",
		Key:    "channel",
		Column: [...]string{"Canal", ""},
	},
	"client": {
		Text:   [...]string{"Cliente: %v %v\n", ""},
		Unit:   "",
		Key:    "client",
		Column: [...]string{"Sección", ""},
	},
	"total": {
		Text:   [...]string{"", ""},
		Unit:   "segundos",
		Key:    "duration",
		Column: [...]string{"", ""},
	},
}
