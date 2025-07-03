package model

import "time"

type CreatedAtInfo struct {
	Date         time.Time `json:"date"`
	TimezoneType int       `json:"timezone_type"`
	Timezone     string    `json:"timezone"`
}
