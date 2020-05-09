package rowscanner

import "reflect"

type Column interface {
	Name() string
	Get(value reflect.Value) reflect.Value
}

type RowFunc func(columns []Column, values []interface{}) error

type Rows interface {
	AltNameTag() string
	IterateUsing(fn RowFunc) error
}
