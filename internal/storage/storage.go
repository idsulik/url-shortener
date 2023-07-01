package storage

import "errors"

var (
	ErrUrlNotFound = errors.New("url not found")
	ErrorUrlExists = errors.New("url already exists")
)
