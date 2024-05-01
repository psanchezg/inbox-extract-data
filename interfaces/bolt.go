package interfaces

import "time"

type BoltReceipt struct {
	Total       float64   `json:"total"`
	Subtotal    float64   `json:"subtotal"`
	Desbloquear float64   `json:"desbloquear"`
	Descuento   float64   `json:"descuento"`
	Fecha       time.Time `json:"fecha"`
	Duracion    int32     `json:"duracion"`
	Distancia   float64   `json:"distancia"`
	Snippet     string    `json:"snippet"`
}

type BoltPlan struct {
	Inicio     time.Time `json:"inicio"`
	Fin        time.Time `json:"fin"`
	Duracion   int64     `json:"duracion"`
	Minutos    int64     `json:"minutos"`
	MinutosDia int64     `json:"minutos_dia"`
	Total      float64   `json:"total"`
}
