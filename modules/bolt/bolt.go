package bolt

import (
	"github.com/psanchezg/inbox-extract-data/interfaces"
)

var Props = map[string]interfaces.Prop{
	"from": {
		Text: [...]string{"Periodo de %v a %v\n", ""},
		Unit: "",
		Key:  "inicio",
	},
	"to": {
		Text: [...]string{"******* Plan activo **********\n", ""},
		Unit: "",
		Key:  "fin",
	},
	"total": {
		Text: [...]string{"Total pagado (incluyendo bono): %v %v\n", "Importe total: %v %v\n"},
		Unit: "€",
		Key:  "total",
	},
	"subtotal": {
		Text: [...]string{"Importe total: %v %v\n", ""},
		Unit: "€",
		Key:  "subtotal",
	},
	"purchased": {
		Text: [...]string{"Plan de %v a %v\n", ""},
		Unit: "",
		Key:  "purchased",
	},
	"duration": {
		Text: [...]string{"Duración del bono: %v %v\n", ""},
		Unit: "días",
		Key:  "duracion",
	},
	"duration2": {
		Text: [...]string{"Duración del bono: %v %v\n", ""},
		Unit: "minutos",
		Key:  "minutos",
	},
	"duration3": {
		Text: [...]string{"Tiempo adicional usado (fuera bono): %v %v\n", ""},
		Unit: "minutos",
		Key:  "minutos_dia",
	},
	"usage": {
		Text: [...]string{"", ""},
		Unit: "",
		Key:  "uso",
	},
	"usage_firsttravel": {
		Text: [...]string{"", ""},
		Unit: "",
		Key:  "primer_viaje",
	},
	"usage_time": {
		Text: [...]string{"Tiempo total: %v %v\n", ""},
		Unit: "minutos",
		Key:  "tiempo",
	},
	"usage_distance": {
		Text: [...]string{"Distancia total: %v %v\n", ""},
		Unit: "kms",
		Key:  "distancia",
	},
	"usage_paid": {
		Text: [...]string{"Importe total pagado: %v %v\n", "Pagado adicional al bono: %v %v\n"},
		Unit: "€",
		Key:  "pagado",
	},
	"usage_service": {
		Text: [...]string{"Importe total del servicio: %v %v\n", "Coste total del servicio (sin bono): %v %v\n"},
		Unit: "€",
		Key:  "servicio",
	},
}
