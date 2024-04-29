package interfaces

import "time"

type BoltReceipt struct {
	Total       float64   `json:"total"`
	Subtotal    float64   `json:"subtotal"`
	Desbloquear float64   `json:"desbloquear"`
	Descuento   float64   `json:"descuento"`
	Fecha       time.Time `json:"fecha"`
	Duracion    int32     `json:"duracion"`
	Snippet     string    `json:"snippet"`
}
