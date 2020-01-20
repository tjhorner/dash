package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func contentTypeJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
}

func respondJSON(resp interface{}, w http.ResponseWriter, r *http.Request) error {
	contentTypeJSON(w, r)
	j := json.NewEncoder(w)
	err := j.Encode(resp)
	return err
}

type apiError struct {
	Error     string `json:"error"`
	Code      string `json:"code"`
	Reference string `json:"reference"`
}

func respondError(err error, code string, w http.ResponseWriter, r *http.Request) {
	contentTypeJSON(w, r)
	w.WriteHeader(400)
	refCode := uuid.New().String()

	fmt.Printf("API error (%s): %s\n", refCode, err)

	errString := "Unknown error"
	if err != nil {
		errString = err.Error()
	}

	humanErr := apiError{
		Error:     errString,
		Code:      code,
		Reference: refCode,
	}

	internalErr := respondJSON(humanErr, w, r)
	if internalErr != nil {
		w.Write([]byte("{\"error\":\"Internal Server Error\"}"))
	}
}
