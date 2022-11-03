package models

type Response struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
	Status  string      `json:"status"`
}
