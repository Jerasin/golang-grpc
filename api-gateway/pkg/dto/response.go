package dto

type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}
