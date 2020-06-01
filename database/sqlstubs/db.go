package sqlstubs

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/FantLab/go-kit/database/sqlapi"
)

type StubDB struct {
	ReadTable  map[string]interface{}
	WriteTable map[string]sqlapi.Result
}

func (db *StubDB) InTransaction(perform func(sqlapi.ReaderWriter) error) error {
	return perform(db)
}

func (db *StubDB) Write(ctx context.Context, q *sqlapi.Query) sqlapi.Result {
	return db.WriteTable[q.String()]
}

func (db *StubDB) Read(ctx context.Context, q *sqlapi.Query, output interface{}) error {
	if stub := db.ReadTable[q.String()]; stub != nil {
		reflect.Indirect(reflect.ValueOf(output)).Set(reflect.ValueOf(stub))
		return nil
	}
	return sql.ErrNoRows
}
