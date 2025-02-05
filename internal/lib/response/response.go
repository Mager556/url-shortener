package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(err string) Response {
	return Response{
		Status: StatusError,
		Error:  err,
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errMsgs []string

	for _, value := range errors {
		switch value.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field", value.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not a valid url", value.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", value.Field()))
		}
	}

	return Error(strings.Join(errMsgs, ", "))
}
