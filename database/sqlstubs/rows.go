package sqlstubs

import (
	"reflect"

	"github.com/FantLab/go-kit/database/rowscanner"
)

type StubRows struct {
	Values  [][]interface{}
	Columns []rowscanner.Column
	Err     error
}

func (rows *StubRows) Error() error {
	return rows.Err
}

func (rows *StubRows) Scan(output interface{}) error {
	if rows.Err != nil {
		return rows.Err
	}
	return rowscanner.Scan(output, rows)
}

func (rows *StubRows) AltNameTag() string {
	return "db"
}

func (rows *StubRows) IterateUsing(fn rowscanner.RowFunc) error {
	if rows.Err != nil {
		return rows.Err
	}
	for _, values := range rows.Values {
		if err := fn(rows.Columns, values); err != nil {
			return err
		}
	}
	return nil
}

// *******************************************************

type StubColumn string

func (column StubColumn) Name() string {
	return string(column)
}

func (column StubColumn) Get(value reflect.Value) reflect.Value {
	return value
}
