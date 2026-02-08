package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

type Response struct {
	Status string
	Error  string
}

const (
	StatusOK    = "OK"
	StatusError = "errrr"
)

func GeneralError(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}

func ValidationErrors(errs validator.ValidationErrors) Response {
	var errmessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errmessages = append(errmessages, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errmessages = append(errmessages, fmt.Sprintf("field %s is invalid", err.Field()))
		}

	}
	return Response{
		Error:  strings.Join(errmessages, ", "),
		Status: StatusError,
	}
}
