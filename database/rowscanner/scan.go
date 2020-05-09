package rowscanner

import "reflect"

func isKnownType(k reflect.Kind) bool {
	switch k {
	case
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

func setValue(column Column, value interface{}, output reflect.Value) {
	input := column.Get(reflect.ValueOf(value))

	output.Set(input.Convert(output.Type()))
}

func setValuesToStruct(values []interface{}, columns []Column, output reflect.Value, idxMap map[string]int) {
	for i, value := range values {
		column := columns[i]

		j, ok := idxMap[column.Name()]

		if !ok {
			continue
		}

		setValue(column, value, output.Field(j))
	}
}

func makeFieldNameIndexMapFromStruct(t reflect.Type, altNameTag string) map[string]int {
	m := make(map[string]int)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		m[f.Name] = i

		if altName, ok := f.Tag.Lookup(altNameTag); ok {
			m[altName] = i
		}
	}

	return m
}

func scanMultipleValuesOfStructType(output reflect.Value, elemType reflect.Type, rows Rows) error {
	idxMap := makeFieldNameIndexMapFromStruct(elemType, rows.AltNameTag())

	err := rows.IterateUsing(func(columns []Column, values []interface{}) error {
		newElem := reflect.Indirect(reflect.New(elemType))

		setValuesToStruct(values, columns, newElem, idxMap)

		output.Set(reflect.Append(output, newElem))

		return nil
	})

	return err
}

func scanMultipleValuesOfKnownType(output reflect.Value, elemType reflect.Type, rows Rows) error {
	err := rows.IterateUsing(func(columns []Column, values []interface{}) error {
		if len(columns) != 1 || len(values) != 1 {
			return ErrInvalidColumnCount
		}

		newElem := reflect.Indirect(reflect.New(elemType))

		setValue(columns[0], values[0], newElem)

		output.Set(reflect.Append(output, newElem))

		return nil
	})

	return err
}

func scanMultipleValuesOfKnownMapType(output reflect.Value, keyType, elemType reflect.Type, rows Rows) error {
	err := rows.IterateUsing(func(columns []Column, values []interface{}) error {
		if len(columns) != 2 || len(values) != 2 {
			return ErrInvalidColumnCount
		}

		newKey := reflect.Indirect(reflect.New(keyType))
		newElem := reflect.Indirect(reflect.New(elemType))

		setValue(columns[0], values[0], newKey)
		setValue(columns[1], values[1], newElem)

		output.SetMapIndex(newKey, newElem)

		return nil
	})

	return err
}

func scanSingleValueOfStructType(output reflect.Value, rows Rows) error {
	idxMap := makeFieldNameIndexMapFromStruct(output.Type(), rows.AltNameTag())

	once := false

	err := rows.IterateUsing(func(columns []Column, values []interface{}) error {
		if once {
			return ErrInvalidRowCount
		}

		setValuesToStruct(values, columns, output, idxMap)

		once = true

		return nil
	})

	if err != nil {
		return err
	}

	if !once {
		return ErrInvalidRowCount
	}

	return nil
}

func scanSingleValueOfKnownType(output reflect.Value, rows Rows) error {
	once := false

	err := rows.IterateUsing(func(columns []Column, values []interface{}) error {
		if len(columns) != 1 || len(values) != 1 {
			return ErrInvalidColumnCount
		}

		if once {
			return ErrInvalidRowCount
		}

		setValue(columns[0], values[0], output)

		once = true

		return nil
	})

	if err != nil {
		return err
	}

	if !once {
		return ErrInvalidRowCount
	}

	return nil
}

func scanRowsIntoValue(output reflect.Value, rows Rows) error {
	switch k := output.Type().Kind(); {
	case k == reflect.Slice:
		elemType := output.Type().Elem()

		switch k = elemType.Kind(); {
		case k == reflect.Struct:
			return scanMultipleValuesOfStructType(output, elemType, rows)
		case isKnownType(k):
			return scanMultipleValuesOfKnownType(output, elemType, rows)
		default:
			return ErrUnsupportedType
		}
	case k == reflect.Map:
		keyType := output.Type().Key()
		elemType := output.Type().Elem()

		if isKnownType(keyType.Kind()) && isKnownType(elemType.Kind()) {
			return scanMultipleValuesOfKnownMapType(output, keyType, elemType, rows)
		} else {
			return ErrUnsupportedType
		}
	case k == reflect.Struct:
		return scanSingleValueOfStructType(output, rows)
	case isKnownType(k):
		return scanSingleValueOfKnownType(output, rows)
	default:
		return ErrUnsupportedType
	}
}

func Scan(output interface{}, rows Rows) error {
	value := reflect.ValueOf(output)

	if value.Kind() != reflect.Ptr {
		return ErrNotAPtr
	}

	if value.IsNil() {
		return ErrIsNil
	}

	value = reflect.Indirect(value)

	err := scanRowsIntoValue(value, rows)

	if err != nil {
		value.Set(reflect.Zero(value.Type()))
	}

	return err
}
