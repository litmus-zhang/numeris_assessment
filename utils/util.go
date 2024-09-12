package utils

type ResponseUtils[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}


