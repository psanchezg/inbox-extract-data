package interfaces

type Prop struct {
	Text [2]string `json:"text"`
	Unit string    `json:"unit"`
	Key  string    `json:"key"`
}

type Any interface {
	map[string]interface{}
}
