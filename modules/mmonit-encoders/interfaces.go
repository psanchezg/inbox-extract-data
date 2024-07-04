package mmonitencoders

import "time"

type MmonitAlert struct {
	Channel string    `json:"channel"`
	Encoder string    `json:"encoder"`
	Action  string    `json:"action"`
	Date    time.Time `json:"date"`
	Snippet string    `json:"snippet"`
	Body    string    `json:"body"`
}

type MmonitUsage struct {
	Start    time.Time `json:"start"`
	Stop     time.Time `json:"stop"`
	Duration float64   `json:"duration"`
	Minutes  float64   `json:"minutes"`
	Channel  string    `json:"channel"`
	Client   string    `json:"client"`
	Encoder  string    `json:"encoder"`
}

type MmonitUsageStats struct {
	Total   float64 `json:"total"`
	Minutes int64   `json:"minutes"`
}
