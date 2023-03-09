package models

import "time"

type Audio struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	Duration time.Duration `json:"duration"`
}
