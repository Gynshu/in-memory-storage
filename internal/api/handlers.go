package api

import (
	"encoding/json"
	"github.com/gynshu-one/in-memory-storage/internal/domain"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Handlers struct {
	UseCase domain.Repository
}

// NewHandlers returns a new instance of Handlers.
func NewHandlers(useCase domain.Repository) *Handlers {
	return &Handlers{UseCase: useCase}
}

// Set adds a key-value pair to the in-memory storage.
// Body example:
//
//	{
//	  "key": "key1",
//	  "value": "value1",
//	  "expiration": 0
//	}
//
// expiration is optional and is in seconds
func (h *Handlers) Set(w http.ResponseWriter, r *http.Request) {
	var entity domain.Entity
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		http.Error(w, UnableToParseRequestBody, http.StatusBadRequest)
		return
	}

	if entity.Value == "" {
		http.Error(w, ValueCanNotBeEmpty, http.StatusBadRequest)
		return
	}

	if entity.Expiration < 0 {
		http.Error(w, InvalidDuration, http.StatusBadRequest)
		return
	}

	err = h.UseCase.Set(entity.Key, entity.Value, time.Duration(entity.Expiration)*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(KeyAddedSuccessfully))
	if err != nil {
		log.Error().Err(err).Msg(FailToWriteResponse)
		return
	}

}

// Delete deletes a key-value pair from the in-memory storage.
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, KeyCanNotBeEmpty, http.StatusBadRequest)
		return
	}

	err := h.UseCase.Delete(key)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(KeyDeletedSuccessfully))
	if err != nil {
		log.Error().Err(err).Msg(FailToWriteResponse)
		return
	}
}

// Get returns a value for a given key from the in-memory storage.
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key == "" {
		http.Error(w, KeyCanNotBeEmpty, http.StatusBadRequest)
		return
	}
	value, err := h.UseCase.Get(key)
	if err != nil {
		handleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(value))
	if err != nil {
		log.Error().Err(err).Msg(FailToWriteResponse)
		return
	}
}

// GetAll returns all keys from the in-memory storage.
func (h *Handlers) GetAll(w http.ResponseWriter, r *http.Request) {
	keys, err := h.UseCase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(keys)
	if err != nil {
		log.Error().Err(err).Msg(FailToWriteResponse)
		return
	}
}
