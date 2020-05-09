package sqlbuilder

import (
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/FantLab/go-kit/database/sqlapi"
	"github.com/FantLab/go-kit/env"
)

func InsertInto(tableName string, records ...interface{}) sqlapi.Query {
	return *insertInto(tableName, "db", env.IsDebug(), records)
}

func insertInto(tableName, tagName string, checkTypes bool, records []interface{}) *sqlapi.Query {
	var recordType reflect.Type

	for _, record := range records {
		if recordType == nil {
			recordType = reflect.TypeOf(record)
			if !checkTypes {
				break
			}
		} else if reflect.TypeOf(record) != recordType {
			return nil
		}
	}

	if recordType == nil || recordType.Kind() != reflect.Struct {
		return nil
	}

	fieldNames := extractFieldNames(recordType, tagName)

	n, m := len(fieldNames), len(records)

	if n == 0 {
		return nil
	}

	args := make([]interface{}, 0, m*n)

	for _, record := range records {
		value := reflect.ValueOf(record)

		for j := 0; j < n; j++ {
			args = append(args, value.Field(j).Interface())
		}
	}

	text := makeInsertQueryText(tableName, fieldNames, sqlapi.BindVarChar, m)

	query := sqlapi.NewQuery(text).WithArgs(args...)

	return &query
}

func extractFieldNames(typ reflect.Type, tagName string) (fieldNames []string) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		name := field.Tag.Get(tagName)
		if len(name) == 0 {
			name = field.Name
		}

		fieldNames = append(fieldNames, name)
	}
	return
}

const (
	insertText = "INSERT INTO "
	valuesText = " VALUES "
)

func calculateInsertQuerySize(tableName string, fieldNames []string, bindVarChar rune, count int) (size int) {
	size += len(insertText)
	size += len(tableName)
	size++
	for i, fieldName := range fieldNames {
		if i > 0 {
			size++
		}
		size += len(fieldName)
	}
	size++
	size += len(valuesText)
	size += count*(len(fieldNames)*(utf8.RuneLen(bindVarChar)+1)+2) - 1
	return
}

func makeInsertQueryText(tableName string, fieldNames []string, bindVarChar rune, count int) string {
	var sb strings.Builder
	sb.Grow(calculateInsertQuerySize(tableName, fieldNames, bindVarChar, count))
	sb.WriteString(insertText)
	sb.WriteString(tableName)
	sb.WriteRune('(')
	for i, fieldName := range fieldNames {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(fieldName)
	}
	sb.WriteRune(')')
	sb.WriteString(valuesText)
	for i := 0; i < count; i++ {
		if i > 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('(')
		for j := 0; j < len(fieldNames); j++ {
			if j > 0 {
				sb.WriteRune(',')
			}
			sb.WriteRune(bindVarChar)
		}
		sb.WriteRune(')')
	}
	return sb.String()
}
