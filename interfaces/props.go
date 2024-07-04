package interfaces

type Prop struct {
	Text   [2]string `json:"text"`
	Unit   string    `json:"unit"`
	Key    string    `json:"key"`
	Column [2]string `json:"column"`
}

type Any interface {
	map[string]interface{}
}
