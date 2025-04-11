package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New(validator.WithRequiredStructEnabled())

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteValidationError(w http.ResponseWriter, status int, errs validator.ValidationErrors) {
	formattedErrs := make(map[string]string)

	for i := range errs {
		formattedErrs[strings.ToLower(errs[i].Field())] = writeValidationTagAsErrorString(errs[i])
	}

	// errString := "An error occurred while validating"
	// json, err := json.Marshal(formattedErrs)
	// if err == nil {
	// 	errString = string(json)
	// }

	WriteJSON(w, status, map[string]map[string]string{"error": formattedErrs})
}

func writeValidationTagAsErrorString(fieldErr validator.FieldError) string {
	switch fieldErr.Field() {
	case "Email":
		switch fieldErr.Tag() {
		case "email":
			return "Email is not valid."
		default:
			return "Problem while validating email."
		}
	case "Password":
		switch fieldErr.Tag() {
		case "min":
			// TODO below shouldn't be hardcoded
			return "Password must be at least 8 characters long."
		default:
			return "Problem while validating password."
		}
	default:
		return "Problem?"
	}
}
