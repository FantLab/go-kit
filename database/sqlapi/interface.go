package sqlapi

import "context"

type Result struct {
	LastInsertId int64
	RowsAffected int64
	Error        error
}

type Rows interface {
	Error() error
	Scan(output interface{}) error
}

type Reader interface {
	Read(ctx context.Context, q Query) Rows
}

type Writer interface {
	Write(ctx context.Context, q Query) Result
}

type ReaderWriter interface {
	Reader
	Writer
}

type Transactional interface {
	InTransaction(perform func(ReaderWriter) error) error
}

type DB interface {
	Transactional
	ReaderWriter
}
