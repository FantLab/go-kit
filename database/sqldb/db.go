package sqldb

import (
	"context"
	"database/sql"

	"github.com/FantLab/go-kit/database/sqlapi"
)

func New(sql *sql.DB) sqlapi.DB {
	return &sqlDB{sql: sql}
}

type sqlDB struct {
	sql *sql.DB
}

func (db sqlDB) InTransaction(perform func(sqlapi.ReaderWriter) error) error {
	return inTransaction(db.sql, func(tx *sql.Tx) error {
		return perform(readerWriter{tx})
	})
}

func (db sqlDB) Write(ctx context.Context, q sqlapi.Query) sqlapi.Result {
	return readerWriter{db.sql}.Write(ctx, q)
}

func (db sqlDB) Read(ctx context.Context, q sqlapi.Query) sqlapi.Rows {
	return readerWriter{db.sql}.Read(ctx, q)
}

// *******************************************************

type sqlReaderWriter interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type readerWriter struct {
	sql sqlReaderWriter
}

func (rw readerWriter) Write(ctx context.Context, q sqlapi.Query) sqlapi.Result {
	res, err := rw.sql.ExecContext(ctx, q.Text(), q.Args()...)

	if err != nil {
		return sqlapi.Result{
			Error: err,
		}
	}

	var lastInsertId, rowsAffected int64

	lastInsertId, _ = res.LastInsertId()
	rowsAffected, _ = res.RowsAffected()

	return sqlapi.Result{
		LastInsertId: lastInsertId,
		RowsAffected: rowsAffected,
	}
}

func (rw readerWriter) Read(ctx context.Context, q sqlapi.Query) sqlapi.Rows {
	r, err := rw.sql.QueryContext(ctx, q.Text(), q.Args()...)

	return sqlRows{
		data:           r,
		err:            err,
		allowNullTypes: false,
	}
}

// *******************************************************

func inTransaction(db *sql.DB, fn func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()

			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = fn(tx)

	return
}
