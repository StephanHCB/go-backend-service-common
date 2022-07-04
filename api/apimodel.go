package api

import "time"

type ErrorDto struct {
	Details   *string    `json:"details,omitempty"`
	Message   *string    `json:"message,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`
}

type HealthComponent struct {
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}
