package kde

import "errors"

var (
	// ErrUnknownKernel is returned when an unknown kernel name is provided.
	ErrUnknownKernel = errors.New("kde: unknown kernel")
	// ErrInvalidBandwidth is returned when an invalid bandwidth is provided.
	ErrInvalidBandwidth = errors.New("kde: invalid bandwidth")
	// ErrEmptyData is returned when empty data is provided.
	ErrEmptyData = errors.New("kde: empty data")
	// ErrDataMismatch is returned when data and weights have different lengths.
	ErrDataMismatch = errors.New("kde: data and weights length mismatch")
)
