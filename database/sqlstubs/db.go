package sqlstubs

import (
	"context"
	"database/sql"
	"errors"

	"github.com/FantLab/go-kit/database/sqlapi"
)

var ErrSome = errors.New("sqlstubs: some error")

type (
	StubQueryTable map[string]*StubRows
	StubExecTable  map[string]sqlapi.Result
	StubDB         struct {
		QueryTable StubQueryTable
		ExecTable  StubExecTable
	}
)

func (db *StubDB) InTransaction(perform func(sqlapi.ReaderWriter) error) error {
	return perform(db)
}

func (db *StubDB) Write(ctx context.Context, q sqlapi.Query) sqlapi.Result {
	return db.ExecTable[q.String()]
}

func (db *StubDB) Read(ctx context.Context, q sqlapi.Query) sqlapi.Rows {
	if rows := db.QueryTable[q.String()]; rows != nil {
		return rows
	}
	return sqlapi.NoRows{
		Err: sql.ErrNoRows,
	}
}
