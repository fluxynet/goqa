package repo

import "errors"

var (
	// ErrFailedToWriteCompletely we wrote some but not all :(
	ErrFailedToWriteCompletely = errors.New("failed to write completely")
)
