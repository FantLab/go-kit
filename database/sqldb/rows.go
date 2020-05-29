package sqldb

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/FantLab/go-kit/database/rowscanner"
)

// *******************************************************

type sqlRows struct {
	data *sql.Rows
}

func (rows sqlRows) AltNameTag() string {
	return "db"
}

func (rows sqlRows) IterateUsing(fn rowscanner.RowFunc) error {
	return iterateOverRows(rows.data, false, fn)
}

// *******************************************************

type sqlColumn struct {
	name                string
	takeNonNullSubField bool
}

func (column *sqlColumn) Name() string {
	return column.name
}

func (column *sqlColumn) Get(value reflect.Value) reflect.Value {
	if column.takeNonNullSubField {
		return value.Elem().Field(0)
	}
	return value.Elem()
}

// *******************************************************

func iterateOverRows(rows *sql.Rows, allowNullTypes bool, fn rowscanner.RowFunc) error {
	columnTypes, err := rows.ColumnTypes()

	if err != nil {
		return err
	}

	values, columns := getColumnData(columnTypes, allowNullTypes)

	for rows.Next() {
		err = rows.Scan(values...)

		if err != nil {
			return err
		}

		err = fn(columns, values)

		if err != nil {
			return err
		}
	}

	return rows.Err()
}

func getColumnData(columnTypes []*sql.ColumnType, allowNullTypes bool) ([]interface{}, []rowscanner.Column) {
	size := len(columnTypes)

	values := make([]interface{}, size)
	columns := make([]rowscanner.Column, size)

	for i, columnType := range columnTypes {
		values[i] = reflect.New(columnType.ScanType()).Interface()

		isNullable := strings.HasPrefix(columnType.ScanType().Name(), "Null")

		columns[i] = &sqlColumn{
			name:                columnType.Name(),
			takeNonNullSubField: isNullable && !allowNullTypes,
		}
	}

	return values, columns
}
