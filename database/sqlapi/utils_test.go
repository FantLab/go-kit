package sqlapi

import (
	"testing"
	"time"

	"github.com/FantLab/go-kit/assert"
)

func Test_expandQuery(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		x := expandQuery("? (?) ?", '?', []int{1, 2, 3})

		assert.True(t, x == "? (?,?) ?,?,?")
	})

	t.Run("zero", func(t *testing.T) {
		x := expandQuery("", '?', []int{})

		assert.True(t, x == "")
	})

	t.Run("negative", func(t *testing.T) {
		x := expandQuery("? (?) ?", '?', []int{1, 2})

		assert.True(t, x == "? (?,?) ?")
	})
}

func Test_deepFlat(t *testing.T) {
	t.Run("single1", func(t *testing.T) {
		x, count := deepFlat(16)

		assert.True(t, count == 1)
		assert.DeepEqual(t, x, []interface{}{16})
	})

	t.Run("single2", func(t *testing.T) {
		x, count := deepFlat("a")

		assert.True(t, count == 1)
		assert.DeepEqual(t, x, []interface{}{"a"})
	})

	t.Run("slice", func(t *testing.T) {
		x, count := deepFlat([]string{"x", "y"})

		assert.True(t, count == 2)
		assert.DeepEqual(t, x, []interface{}{"x", "y"})
	})

	t.Run("2Dslice1", func(t *testing.T) {
		x, count := deepFlat([][]int{{1, 2, 3}, {4, 5, 6}})

		assert.True(t, count == 6)
		assert.DeepEqual(t, x, []interface{}{1, 2, 3, 4, 5, 6})
	})

	t.Run("2Dslice1", func(t *testing.T) {
		x, count := deepFlat([][]interface{}{{1, 2, 3}, {"x", "y", "z"}})

		assert.True(t, count == 6)
		assert.DeepEqual(t, x, []interface{}{1, 2, 3, "x", "y", "z"})
	})

	t.Run("4Dslice", func(t *testing.T) {
		x, count := deepFlat([][][][]int{{{{1, 2}, {3, 4}}, {{5, 6}, {7, 8}}}, {{{9, 10}, {11, 12}}, {{13, 14}, {15, 16}}}})

		assert.True(t, count == 16)
		assert.DeepEqual(t, x, []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	})
}

func Test_formatQuery(t *testing.T) {
	t.Run("numbers", func(t *testing.T) {
		x := formatQuery("???", '?', 1, 2, 3)

		assert.True(t, x == "123")
	})

	t.Run("string", func(t *testing.T) {
		x := formatQuery("? ? ?", '?', "x", "y", "z")

		assert.True(t, x == "'x' 'y' 'z'")
	})

	t.Run("time", func(t *testing.T) {
		x := formatQuery("?", '?', time.Date(2010, 10, 11, 15, 20, 33, 0, time.UTC))

		assert.True(t, x == "'2010-10-11 15:20:33'")
	})

	t.Run("multi_spaces", func(t *testing.T) {
		x := formatQuery("   ?   ?   ?   ", '?', 1, 2, "x")

		assert.True(t, x == "1 2 'x'")
	})

	t.Run("complex", func(t *testing.T) {
		x := formatQuery("id = ? and id in (?,?,?,?,?,?)", '?', 1, 2, 3, 4, 5, 6, 7)

		assert.True(t, x == "id = 1 and id in (2,3,4,5,6,7)")
	})
}

func Test_flatQuery(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		text, args := flatQuery("? (?)", []interface{}{"s", []int{1, 2, 3}})

		assert.True(t, text == "? (?,?,?)")
		assert.DeepEqual(t, args, []interface{}{"s", 1, 2, 3})
	})
}
