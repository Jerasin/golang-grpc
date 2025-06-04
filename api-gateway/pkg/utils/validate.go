package utils

import (
	"api-gateway/pkg/constant"
	"strings"

	"github.com/go-playground/validator"
	"github.com/goforj/godump"
)

func ReplaceWithMap(s string, replacements map[string]string) string {
	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}
	return s
}

func Validate(req any) {
	var validate = validator.New()
	if err := validate.Struct(req); err != nil {
		replacements := map[string]string{
			":":   " ",
			".":   " ",
			",":   " ",
			"Key": " ",
		}
		newError := ReplaceWithMap(err.Error(), replacements)
		godump.Dump(newError)
		PanicException(constant.ValidateError, &newError)
	}

}
