package rowscanner

import (
	"reflect"
	"testing"

	"github.com/FantLab/go-kit/assert"
)

// *******************************************************

type _testColumn string

func (column _testColumn) Name() string {
	return string(column)
}

func (column _testColumn) Get(value reflect.Value) reflect.Value {
	return value
}

// *******************************************************

type _testRows struct {
	values  [][]interface{}
	columns []Column
}

func (rows *_testRows) AltNameTag() string {
	return "altname"
}

func (rows *_testRows) IterateUsing(fn RowFunc) error {
	for _, values := range rows.values {
		err := fn(rows.columns, values)

		if err != nil {
			return err
		}
	}

	return nil
}

// *******************************************************

func Test_Scan(t *testing.T) {
	t.Run("negative_output_nil", func(t *testing.T) {
		rows := &_testRows{}

		var x *uint8

		err := Scan(x, rows)

		assert.True(t, err == ErrIsNil)
	})

	t.Run("negative_output_not_a_ptr", func(t *testing.T) {
		rows := &_testRows{}

		var x uint8

		err := Scan(x, rows)

		assert.True(t, err == ErrNotAPtr)
	})

	t.Run("positive_single_value_1", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{1}},
			columns: []Column{
				_testColumn(""),
			},
		}

		var x uint8

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.True(t, x == 1)
	})

	t.Run("positive_single_value_2", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"hello"}},
			columns: []Column{
				_testColumn(""),
			},
		}

		var x string

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.True(t, x == "hello")
	})

	t.Run("negative_single_value_multi_columns", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"hello", "world"}},
			columns: []Column{
				_testColumn(""),
				_testColumn(""),
			},
		}

		var x string

		err := Scan(&x, rows)

		assert.True(t, err == ErrInvalidColumnCount)
		assert.True(t, x == "")
	})

	t.Run("negative_single_value_no_rows", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{},
			columns: []Column{
				_testColumn(""),
			},
		}

		var x string

		err := Scan(&x, rows)

		assert.True(t, err == ErrInvalidRowCount)
		assert.True(t, x == "")
	})

	t.Run("negative_single_value_multi_rows", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"hello"}, {"world"}},
			columns: []Column{
				_testColumn(""),
			},
		}

		var x string

		err := Scan(&x, rows)

		assert.True(t, err == ErrInvalidRowCount)
		assert.True(t, x == "")
	})

	t.Run("positive_single_struct_field_name", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"a", "b"}},
			columns: []Column{
				_testColumn("FirstName"),
				_testColumn("LastName"),
			},
		}

		var x struct {
			FirstName string
			LastName  string
		}

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.True(t, x.FirstName == "a")
		assert.True(t, x.LastName == "b")
	})

	t.Run("positive_single_struct_alt_name", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"a", "b"}},
			columns: []Column{
				_testColumn("first_name"),
				_testColumn("last_name"),
			},
		}

		var x struct {
			FirstName string `altname:"first_name"`
			LastName  string `altname:"last_name"`
		}

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.True(t, x.FirstName == "a")
		assert.True(t, x.LastName == "b")
	})

	t.Run("negative_single_struct_no_rows", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{},
			columns: []Column{
				_testColumn("first_name"),
				_testColumn("last_name"),
			},
		}

		var x struct {
			FirstName string `altname:"first_name"`
			LastName  string `altname:"last_name"`
		}

		err := Scan(&x, rows)

		assert.True(t, err == ErrInvalidRowCount)
		assert.True(t, x.FirstName == "")
		assert.True(t, x.LastName == "")
	})

	t.Run("negative_single_value_multi_rows", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"a", "b"}, {"c", "d"}},
			columns: []Column{
				_testColumn("first_name"),
				_testColumn("last_name"),
			},
		}

		var x struct {
			FirstName string `altname:"first_name"`
			LastName  string `altname:"last_name"`
		}

		err := Scan(&x, rows)

		assert.True(t, err == ErrInvalidRowCount)
		assert.True(t, x.FirstName == "")
		assert.True(t, x.LastName == "")
	})

	t.Run("positive_slice_known_type_1", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{1}, {2}, {3}},
			columns: []Column{
				_testColumn(""),
			},
		}

		var x []int

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.DeepEqual(t, x, []int{1, 2, 3})
	})

	t.Run("positive_slice_known_type_2", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"x"}, {"y"}, {"z"}},
			columns: []Column{
				_testColumn(""),
			},
		}

		var x []string

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.DeepEqual(t, x, []string{"x", "y", "z"})
	})

	t.Run("negative_slice_known_type", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{1}, {2}, {3}},
			columns: []Column{
				_testColumn(""),
				_testColumn(""),
			},
		}

		var x []int

		err := Scan(&x, rows)

		assert.True(t, err == ErrInvalidColumnCount)
		assert.True(t, x == nil)
	})

	t.Run("positive_slice_alt_name", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"a", "b"}, {"c", "d"}},
			columns: []Column{
				_testColumn("first_name"),
				_testColumn("last_name"),
			},
		}

		var x []struct {
			FirstName string `altname:"first_name"`
			LastName  string `altname:"last_name"`
		}

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.True(t, x[0].FirstName == "a")
		assert.True(t, x[0].LastName == "b")
		assert.True(t, x[1].FirstName == "c")
		assert.True(t, x[1].LastName == "d")
	})

	t.Run("positive_slice_mix_names", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{"a", "b"}, {"c", "d"}},
			columns: []Column{
				_testColumn("FirstName"),
				_testColumn("last_name"),
			},
		}

		var x []struct {
			FirstName string
			LastName  string `altname:"last_name"`
		}

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.True(t, x[0].FirstName == "a")
		assert.True(t, x[0].LastName == "b")
		assert.True(t, x[1].FirstName == "c")
		assert.True(t, x[1].LastName == "d")
	})

	t.Run("positive_complex_slice", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{
				{"a", "b", 1, 2, true, 2.8},
				{"c", "d", 3, 4, false, 3.2},
			},
			columns: []Column{
				_testColumn("first_name"),
				_testColumn("last_name"),
				_testColumn("id1"),
				_testColumn("id2"),
				_testColumn("is_closed"),
				_testColumn("coef"),
			},
		}

		type testData struct {
			FirstName string  `altname:"first_name"`
			LastName  string  `altname:"last_name"`
			Id1       int     `altname:"id1"`
			Id2       uint8   `altname:"id2"`
			IsClosed  bool    `altname:"is_closed"`
			Coef      float64 `altname:"coef"`
		}

		var x []testData

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.DeepEqual(t, x[0], testData{
			FirstName: "a",
			LastName:  "b",
			Id1:       1,
			Id2:       2,
			IsClosed:  true,
			Coef:      2.8,
		})
		assert.DeepEqual(t, x[1], testData{
			FirstName: "c",
			LastName:  "d",
			Id1:       3,
			Id2:       4,
			IsClosed:  false,
			Coef:      3.2,
		})
	})

	t.Run("positive_map_scan", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{1, "v1"}, {2, "v2"}},
			columns: []Column{
				_testColumn(""),
				_testColumn(""),
			},
		}

		x := make(map[uint]string)

		err := Scan(&x, rows)

		assert.True(t, err == nil)
		assert.DeepEqual(t, x, map[uint]string{1: "v1", 2: "v2"})
	})

	t.Run("negative_map_scan", func(t *testing.T) {
		rows := &_testRows{
			values: [][]interface{}{{1, "v1", true}, {2, "v2", true}},
			columns: []Column{
				_testColumn(""),
				_testColumn(""),
				_testColumn(""),
			},
		}

		x := make(map[uint]string)

		err := Scan(&x, rows)

		assert.True(t, err == ErrInvalidColumnCount)
		assert.True(t, len(x) == 0)
	})
}
