package brand

import (
	"github.com/pkg/errors"
)

var (
	ErrBrandNotFound    = errors.New("brand not found")
	ErrBrandSoftDeleted = errors.New("brand has been soft-deleted")
)
