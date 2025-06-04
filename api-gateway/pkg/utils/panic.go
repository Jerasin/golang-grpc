package utils

import (
	"api-gateway/pkg/constant"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/goforj/godump"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func PanicException(responseKey constant.ResponseStatus, message *string) {
	if message == nil {
		PanicException_(responseKey.GetResponseStatus(), responseKey.GetResponseMessage())
		return
	} else {
		PanicException_(responseKey.GetResponseStatus(), responseKey.GetResponseMessage())
	}

}

func PanicException_(key string, message string) {
	err := errors.New(message)
	err = fmt.Errorf("%s: %w", key, err)
	if err != nil {
		panic(err)
	}
}

func GrpcPanicException(err error, message *string) {
	st, ok := status.FromError(err)
	if ok {
		godump.Dump("Code", st.Code())
		godump.Dump("Message", st.Message())
		var newErr error
		if message == nil {
			newErr = fmt.Errorf("%s: %s", st.Code(), st.Message())
		} else {
			newErr = fmt.Errorf("%s: %s", st.Code(), *message)
			godump.Dump("Message", newErr.Error())
		}
		if newErr != nil {
			panic(newErr)
		}
	}
}

func PanicHandler(c *fiber.Ctx) error {
	if err := recover(); err != nil {
		log.Errorf("PanicHandler = %+v  \n", err)
		str := fmt.Sprint(err)
		strArr := strings.Split(str, ":")
		key := strArr[0]
		msg := strings.Trim(strArr[1], " ")

		log.Errorf("key = %+v  \n", key)
		switch key {
		case
			constant.BadRequest.GetResponseStatus():
			return c.Status(http.StatusBadRequest).JSON(BuildResponse_(http.StatusBadRequest, msg, Null()))
		case
			codes.AlreadyExists.String():
			return c.Status(http.StatusBadRequest).JSON(BuildResponse_(http.StatusBadRequest, msg, Null()))
		case
			constant.ValidateError.GetResponseStatus():
			return c.Status(http.StatusUnprocessableEntity).JSON(BuildResponse_(http.StatusUnprocessableEntity, msg, Null()))
		case
			constant.Unauthorized.GetResponseStatus():
			return c.Status(http.StatusUnauthorized).JSON(BuildResponse_(http.StatusUnauthorized, msg, Null()))
		case
			constant.DataIsExit.GetResponseStatus():
			return c.Status(http.StatusBadRequest).JSON(BuildResponse_(http.StatusBadRequest, msg, Null()))
		default:
			return c.Status(http.StatusInternalServerError).JSON(BuildResponse_(http.StatusInternalServerError, msg, Null()))
		}
	}

	return nil
}
