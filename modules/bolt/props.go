package bolt

import (
	"github.com/psanchezg/inbox-extract-data/interfaces"
)

var Props = map[string]interfaces.Prop{
	"from": {
		Text:   [...]string{"Periodo de %v a %v\n", ""},
		Unit:   "",
		Key:    "inicio",
		Column: [...]string{"Inicio", ""},
	},
	"to": {
		Text:   [...]string{"******* Plan activo **********\n", ""},
		Unit:   "",
		Key:    "fin",
		Column: [...]string{"Fin", ""},
	},
	"total": {
		Text:   [...]string{"Total pagado (incluyendo bono): %v %v\n", "Importe total: %v %v\n"},
		Unit:   "€",
		Key:    "total",
		Column: [...]string{"Total", ""},
	},
	"subtotal": {
		Text:   [...]string{"Importe total: %v %v\n", ""},
		Unit:   "€",
		Key:    "subtotal",
		Column: [...]string{"Subtotal", ""},
	},
	"purchased": {
		Text:   [...]string{"Plan de %v a %v\n", ""},
		Unit:   "",
		Key:    "purchased",
		Column: [...]string{"Comprado", ""},
	},
	"duration": {
		Text:   [...]string{"Duración del bono: %v %v\n", "Dias restantes del bono: %v\n"},
		Unit:   "días",
		Key:    "duracion",
		Column: [...]string{"Duración (días)", "Restante (días)"},
	},
	"duration2": {
		Text:   [...]string{"Duración del bono: %v %v\n", ""},
		Unit:   "minutos",
		Key:    "minutos",
		Column: [...]string{"Duración (mins)", ""},
	},
	"duration3": {
		Text:   [...]string{"Tiempo adicional usado (fuera bono): %v %v\n", ""},
		Unit:   "minutos",
		Key:    "minutos_dia",
		Column: [...]string{"Tiempo adicional usado (mins)", ""},
	},
	"usage": {
		Text:   [...]string{"", ""},
		Unit:   "",
		Key:    "uso",
		Column: [...]string{"Uso", ""},
	},
	"usage_firsttravel": {
		Text:   [...]string{"", ""},
		Unit:   "",
		Key:    "primer_viaje",
		Column: [...]string{"Primer viaje", ""},
	},
	"usage_time": {
		Text:   [...]string{"Tiempo total: %v %v\n", ""},
		Unit:   "minutos",
		Key:    "tiempo",
		Column: [...]string{"Tiempo usado (mins)", ""},
	},
	"usage_distance": {
		Text:   [...]string{"Distancia total: %v %v\n", ""},
		Unit:   "kms",
		Key:    "distancia",
		Column: [...]string{"Distancia usada (kms)", ""},
	},
	"usage_paid": {
		Text:   [...]string{"Importe total pagado: %v %v\n", "Pagado adicional al bono: %v %v\n"},
		Unit:   "€",
		Key:    "pagado",
		Column: [...]string{"Pagado total", "Pagado adicional"},
	},
	"usage_service": {
		Text:   [...]string{"Importe total del servicio: %v %v\n", "Coste total del servicio (sin bono): %v %v\n"},
		Unit:   "€",
		Key:    "servicio",
		Column: [...]string{"", "Coste servicio (sin bono)"},
	},
}
