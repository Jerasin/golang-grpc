package utils

import (
	"github.com/go-playground/validator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Validate(req any) error {
	var validate = validator.New()

	if err := validate.Struct(req); err != nil {
		// คืน error message จาก validator
		return status.Errorf(
			codes.InvalidArgument,
			"Validation error: %v",
			err.Error(),
		)
	}
	return nil
}
