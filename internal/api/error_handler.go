package api

import (
	"errors"
	"github.com/gynshu-one/in-memory-storage/internal/domain"
	"net/http"
)

const (
	UnableToParseRequestBody = "Unable to parse request body"
	ValueCanNotBeEmpty       = "Value can not be empty"
	KeyCanNotBeEmpty         = "Key can not be empty"
	KeyAddedSuccessfully     = "Key added successfully"
	KeyDeletedSuccessfully   = "Key deleted successfully"
	FailToWriteResponse      = "Failed to write response"
	InvalidDuration          = "Invalid duration"
)

func handleError(err error, w http.ResponseWriter) {
	switch {
	case errors.Is(err, domain.ErrKeyNotFound):
		http.Error(w, err.Error(), http.StatusNoContent)
		return
	case errors.Is(err, domain.ErrKeyExpired):
		http.Error(w, err.Error(), http.StatusGone)
		return
	case errors.Is(err, domain.ErrStorageEmpty):
		http.Error(w, err.Error(), http.StatusNoContent)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
