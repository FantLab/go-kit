package sqlapi

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

const (
	BindVarChar = '?'
	TimeLayout  = "2006-01-02 15:04:05"
)

// *******************************************************

func flatQuery(text string, args []interface{}) (string, []interface{}) {
	newArgs, counts := flatArgs(args)
	newText := expandQuery(text, BindVarChar, counts)
	return newText, newArgs
}

func expandQuery(q string, bindVarChar rune, counts []int) string {
	end := len(counts) - 1
	cursor := 0

	var sb strings.Builder

	for _, char := range q {
		if char != bindVarChar {
			sb.WriteRune(char)
			continue
		}

		if cursor > end {
			sb.WriteRune(bindVarChar)

			continue
		}

		for j := 0; j < counts[cursor]-1; j++ {
			sb.WriteRune(bindVarChar)
			sb.WriteRune(',')
		}

		sb.WriteRune(bindVarChar)

		cursor += 1
	}

	return sb.String()
}

func flatArgs(args []interface{}) ([]interface{}, []int) {
	var flatSlice []interface{}

	counts := make([]int, len(args))

	for i, arg := range args {
		flatArg, count := deepFlat(arg)

		flatSlice = append(flatSlice, flatArg...)

		counts[i] = count
	}

	return flatSlice, counts
}

func deepFlat(input interface{}) ([]interface{}, int) {
	var flatSlice []interface{}
	var totalCount int

	queue := []interface{}{input}

	for len(queue) > 0 {
		item := queue[0]
		queue = queue[1:]

		value := reflect.ValueOf(item)

		if value.Kind() != reflect.Slice {
			flatSlice = append(flatSlice, item)
			totalCount += 1
			continue
		}

		for i := 0; i < value.Len(); i++ {
			queue = append(queue, value.Index(i).Interface())
		}
	}

	return flatSlice, totalCount
}

// *******************************************************

func formatQuery(q string, bindVarChar rune, args ...interface{}) string {
	end := len(args)
	cursor := 0

	var sb strings.Builder

	prevIsPrint := false
	shouldAppendSpace := false

	for _, char := range q {
		if unicode.IsPrint(char) && !unicode.IsSpace(char) {
			if shouldAppendSpace {
				sb.WriteRune(' ')

				shouldAppendSpace = false
			}

			if char == bindVarChar {
				if cursor < end {
					sb.WriteString(formatArg(args[cursor]))

					cursor += 1
				}
			} else {
				sb.WriteRune(char)
			}

			prevIsPrint = true
		} else {
			if prevIsPrint {
				shouldAppendSpace = true
			}

			prevIsPrint = false
		}
	}

	return sb.String()
}

func formatArg(arg interface{}) string {
	switch x := arg.(type) {
	case string:
		return fmt.Sprintf("'%s'", x)
	case time.Time:
		return fmt.Sprintf("'%s'", x.Format(TimeLayout))
	default:
		return fmt.Sprintf("%v", arg)
	}
}
