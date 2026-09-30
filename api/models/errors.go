package models

import "errors"

// Sentinel errors so the HTTP layer can translate a failure into the right
// status code instead of answering 400 for everything.
var (
	// ErrNotFound covers unknown books, unknown parts and missing chapters.
	ErrNotFound = errors.New("not found")

	// ErrBadRequest covers malformed requests and unsupported operations.
	ErrBadRequest = errors.New("bad request")
)
