package storage

import (
	"github.com/gynshu-one/in-memory-storage/internal/domain"
	"github.com/stretchr/testify/assert"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewInMemory(t *testing.T) {
	tests := []struct {
		name string
		want *storage
	}{
		{
			name: "NewInMemory returns a new instance of storage",
			want: &storage{
				mu:      &sync.RWMutex{},
				storage: make(map[string]domain.Entity),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewInMemory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInMemory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_Delete(t *testing.T) {
	type fields struct {
		mu      *sync.RWMutex
		storage map[string]domain.Entity
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Delete a key that exists in the storage",
			fields: fields{
				mu: &sync.RWMutex{},
				storage: map[string]domain.Entity{
					"key1": {
						Value:      "value1",
						Expiration: time.Now().Add(time.Minute).UnixNano(),
					},
				},
			},
			args: args{
				key: "key1",
			},
			wantErr: false,
		},
		{
			name: "Delete a key that does not exist in the storage",
			fields: fields{
				mu:      &sync.RWMutex{},
				storage: make(map[string]domain.Entity),
			},
			args: args{
				key: "key1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			if err := s.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("storage.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_storage_Get(t *testing.T) {
	type fields struct {
		mu      *sync.RWMutex
		storage map[string]domain.Entity
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
	}{
		{
			name: "Get a key that exists in the storage",
			fields: fields{
				mu: &sync.RWMutex{},
				storage: map[string]domain.Entity{
					"key1": {
						Value:      "value1",
						Expiration: time.Now().Add(time.Minute).UnixNano(),
					},
				},
			},
			args: args{
				key: "key1",
			},
			want:    "value1",
			wantErr: nil,
		},
		{
			name: "Get a key that does not exist in the storage",
			fields: fields{
				mu:      &sync.RWMutex{},
				storage: make(map[string]domain.Entity),
			},
			args: args{
				key: "key1",
			},
			want:    "",
			wantErr: domain.ErrKeyNotFound,
		},
		{
			name: "Get a key that has expired",
			fields: fields{
				mu: &sync.RWMutex{},
				storage: map[string]domain.Entity{
					"key1": {
						Value:      "value1",
						Expiration: time.Now().Add(-time.Minute).UnixNano(),
					},
				},
			},
			args: args{
				key: "key1",
			},
			want:    "",
			wantErr: domain.ErrKeyExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			got, err := s.Get(tt.args.key)
			if !assert.Equal(t, tt.wantErr, err) {
				t.Errorf("storage.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("storage.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_GetAll(t *testing.T) {
	type fields struct {
		mu      *sync.RWMutex
		storage map[string]domain.Entity
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		{
			name: "Get all key-value pairs from the storage",
			fields: fields{
				mu: &sync.RWMutex{},
				storage: map[string]domain.Entity{
					"key1": {
						Value:      "value1",
						Expiration: time.Now().Add(time.Minute).UnixNano(),
					},
					"key2": {
						Value:      "value2",
						Expiration: time.Now().Add(time.Minute).UnixNano(),
					},
				},
			},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			wantErr: false,
		},
		{
			name: "Get all key-value pairs from an empty storage",
			fields: fields{
				mu:      &sync.RWMutex{},
				storage: make(map[string]domain.Entity),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			got, err := s.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("storage.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("storage.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_Set(t *testing.T) {
	type fields struct {
		mu      *sync.RWMutex
		storage map[string]domain.Entity
	}
	type args struct {
		key       string
		value     string
		expiresAt time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Set a new key-value pair in the storage",
			fields: fields{
				mu:      &sync.RWMutex{},
				storage: make(map[string]domain.Entity),
			},
			args: args{
				key:       "key1",
				value:     "value1",
				expiresAt: time.Minute,
			},
			wantErr: false,
		},
		{
			name: "Set a key-value pair with an existing key in the storage",
			fields: fields{
				mu: &sync.RWMutex{},
				storage: map[string]domain.Entity{
					"key1": {
						Value:      "value1",
						Expiration: time.Now().Add(time.Minute).UnixNano(),
					},
				},
			},
			args: args{
				key:       "key1",
				value:     "value2",
				expiresAt: time.Minute,
			},
			wantErr: false,
		},
		{
			name: "Set a key-value pair with an expired TTL in the storage",
			fields: fields{
				mu: &sync.RWMutex{},
				storage: map[string]domain.Entity{
					"key1": {
						Value:      "value1",
						Expiration: time.Now().Add(-time.Minute).UnixNano(),
					},
				},
			},
			args: args{
				key:       "key1",
				value:     "value2",
				expiresAt: time.Minute,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				mu:      tt.fields.mu,
				storage: tt.fields.storage,
			}
			if err := s.Set(tt.args.key, tt.args.value, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("storage.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
