package sqlstubs

import (
	"context"
	"testing"

	"github.com/FantLab/go-kit/database/sqlapi"

	"github.com/FantLab/go-kit/assert"
)

func Test_Stubs(t *testing.T) {
	t.Run("positive_1", func(t *testing.T) {
		var output []int

		db := StubDB{
			ReadTable: make(map[string]interface{}),
		}

		db.ReadTable["x"] = interface{}([]int{1, 2, 3, 4})

		err := db.Read(context.Background(), sqlapi.NewQuery("x"), &output)

		assert.DeepEqual(t, output, []int{1, 2, 3, 4})
		assert.True(t, err == nil)
	})

	t.Run("positive_2", func(t *testing.T) {
		var output int

		db := StubDB{
			ReadTable: make(map[string]interface{}),
		}

		db.ReadTable["x"] = interface{}(10)

		err := db.Read(context.Background(), sqlapi.NewQuery("x"), &output)

		assert.True(t, output == 10)
		assert.True(t, err == nil)
	})

	t.Run("positive_3", func(t *testing.T) {
		type xx struct {
			x int
			s string
		}

		var output xx

		db := StubDB{
			ReadTable: make(map[string]interface{}),
		}

		db.ReadTable["x"] = interface{}(xx{
			x: 10,
			s: "test",
		})

		err := db.Read(context.Background(), sqlapi.NewQuery("x"), &output)

		assert.DeepEqual(t, output, xx{
			x: 10,
			s: "test",
		})
		assert.True(t, err == nil)
	})
}
