package model

import "github.com/pkg/errors"

var (
	ErrModelNotFound    = errors.New("model not found")
	ErrModelSoftDeleted = errors.New("model has been soft-deleted")
)
