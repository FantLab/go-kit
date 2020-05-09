package sqlbuilder

import (
	"testing"

	"github.com/FantLab/go-kit/assert"
)

type x struct {
	Id     int    `db:"id"`
	Name   string `db:"name"`
	Count  int    `db:"count"`
	Active bool
}

func Test_Insert(t *testing.T) {
	t.Run("positive_1", func(t *testing.T) {
		entries := []interface{}{
			x{Id: 1, Name: "a", Count: 100, Active: true},
			x{Id: 2, Name: "b", Count: 200, Active: false},
		}

		query := insertInto("table", "db", true, entries)

		assert.DeepEqual(t, query.Args(), []interface{}{1, "a", 100, true, 2, "b", 200, false})
		assert.True(t, query.Text() == "INSERT INTO table(id,name,count,Active) VALUES (?,?,?,?),(?,?,?,?)")
	})

	t.Run("positive_2", func(t *testing.T) {
		entries := []interface{}{
			x{Id: 1, Name: "a", Count: 100, Active: true},
			x{Id: 2, Name: "b", Count: 200, Active: false},
		}

		query := insertInto("table", "db", false, entries)

		assert.DeepEqual(t, query.Args(), []interface{}{1, "a", 100, true, 2, "b", 200, false})
		assert.True(t, query.Text() == "INSERT INTO table(id,name,count,Active) VALUES (?,?,?,?),(?,?,?,?)")
	})

	t.Run("negative_1", func(t *testing.T) {
		query := insertInto("table", "db", true, nil)

		assert.True(t, query == nil)
	})

	t.Run("negative_2", func(t *testing.T) {
		query := insertInto("table", "db", true, []interface{}{struct{}{}})

		assert.True(t, query == nil)
	})

	t.Run("negative_3", func(t *testing.T) {
		query := insertInto("table", "db", true, []interface{}{1, 2})

		assert.True(t, query == nil)
	})

	t.Run("negative_4", func(t *testing.T) {
		query := insertInto("table", "db", true, []interface{}{1, ""})

		assert.True(t, query == nil)
	})
}

func Benchmark_insertIntoSafe(b *testing.B) {
	n := 1000

	entries := make([]interface{}, 0, n)

	for i := 0; i < n; i++ {
		entries = append(entries, x{Id: 1, Name: "a", Count: 1000, Active: true})
	}

	for i := 0; i < b.N; i++ {
		q := insertInto("table", "db", true, entries)
		_ = q
	}
}

func Benchmark_insertIntoUnsafe(b *testing.B) {
	n := 1000

	entries := make([]interface{}, 0, n)

	for i := 0; i < n; i++ {
		entries = append(entries, x{Id: 1, Name: "a", Count: 1000, Active: true})
	}

	for i := 0; i < b.N; i++ {
		q := insertInto("table", "db", false, entries)
		_ = q
	}
}
