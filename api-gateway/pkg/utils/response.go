package utils

import (
	"api-gateway/pkg/constant"
	"api-gateway/pkg/dto"
)

func Null() any {
	return nil
}

func BuildResponse[T any](responseStatus constant.ResponseStatus, data T) dto.Response[T] {
	return BuildResponse_(responseStatus.GetResponseStatusCode(), responseStatus.GetResponseMessage(), data)
}

func BuildResponse_[T any](status int, message string, data T) dto.Response[T] {
	return dto.Response[T]{
		Status:  status,
		Message: message,
		Data:    data,
		Error:   "",
	}
}
