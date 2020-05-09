package rowscanner

import "errors"

var (
	ErrInvalidRowCount    = errors.New("rowscanner: invalid row count")
	ErrInvalidColumnCount = errors.New("rowscanner: invalid column count")
	ErrUnsupportedType    = errors.New("rowscanner: unsupported type")
	ErrIsNil              = errors.New("rowscanner: output value must not be nil")
	ErrNotAPtr            = errors.New("rowscanner: output value must be a pointer")
)
