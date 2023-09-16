package api

import (
	"bytes"
	"encoding/json"
	"github.com/gynshu-one/in-memory-storage/internal/domain"
	"github.com/gynshu-one/in-memory-storage/internal/infra/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers_Set(t *testing.T) {
	type args struct {
		key        string
		value      string
		expiration int64
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
	}{
		{
			name: "Set returns 201 Created for a valid request",
			args: args{
				key:        "key1",
				value:      "value1",
				expiration: 0,
			},
			wantStatus: http.StatusCreated,
			wantBody:   KeyAddedSuccessfully,
		},
		{
			name: "Set returns 400 Bad Request for an invalid request",
			args: args{
				key:        "key2",
				value:      "",
				expiration: 0,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   ValueCanNotBeEmpty,
		},
		{
			name: "Set returns 500 Internal Server Error for a server error",
			args: args{
				key:        "key3",
				value:      "value3",
				expiration: -1,
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   InvalidDuration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(storage.NewInMemory())

			entity := domain.Entity{
				Key:        tt.args.key,
				Value:      tt.args.value,
				Expiration: tt.args.expiration,
			}
			body, err := json.Marshal(entity)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/set", bytes.NewReader(body))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Set)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}

func TestHandlers_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
	}{
		{
			name: "Delete returns 200 OK for a valid request",
			args: args{
				key: "key1",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Delete returns 400 Bad Request for an invalid request",
			args: args{
				key: "",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Delete returns 404 Not Found for a missing key",
			args: args{
				key: "key2",
			},
			wantStatus: http.StatusNoContent,
		},
	}
	stor := storage.NewInMemory()
	_ = stor.Set("key1", "value1", 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(stor)

			req, err := http.NewRequest(http.MethodDelete, "/delete?key="+tt.args.key, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Delete)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestHandlers_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantBody   string
		wantErr    error
	}{
		{
			name: "Get returns 200 OK for a valid request",
			args: args{
				key: "key1",
			},
			wantStatus: http.StatusOK,
			wantBody:   "value1",
		},
		{
			name: "Get returns 400 Bad Request for an invalid request",
			args: args{
				key: "",
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   KeyCanNotBeEmpty,
		},
		{
			name: "Get returns 404 Not Found for a missing key",
			args: args{
				key: "key2",
			},
			wantStatus: http.StatusNoContent,
			wantBody:   domain.ErrKeyNotFound.Error(),
		},
	}
	stor := storage.NewInMemory()
	_ = stor.Set("key1", "value1", 0)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandlers(stor)

			req, err := http.NewRequest(http.MethodGet, "/get?key="+tt.args.key, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.Get)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			assert.Equal(t, tt.wantBody, strings.TrimSpace(rr.Body.String()))
		})
	}
}
