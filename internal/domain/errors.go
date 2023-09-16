package domain

import "errors"

var (
	ErrKeyExpired   = errors.New("key expired")
	ErrKeyNotFound  = errors.New("key not found")
	ErrStorageEmpty = errors.New("storage is empty")
)
