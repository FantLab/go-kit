package signed

import (
	"crypto/ed25519"
	"encoding/base64"
	"testing"

	"github.com/FantLab/go-kit/assert"
)

func Test_Coder(t *testing.T) {
	coder, _ := Generate()

	t.Run("positive", func(t *testing.T) {
		x := coder.Encode([]byte("success"))

		y, err := coder.Decode(x)

		assert.True(t, err == nil)
		assert.True(t, string(y) == "success")
	})

	t.Run("negative_1", func(t *testing.T) {
		y, err := coder.Decode([]byte(""))

		assert.True(t, err == ErrInput)
		assert.True(t, y == nil)
	})

	t.Run("negative_2", func(t *testing.T) {
		y, err := coder.Decode([]byte(base64.RawURLEncoding.EncodeToString([]byte("test")) + ".sign"))

		assert.True(t, err == ErrSign)
		assert.True(t, y == nil)
	})

	t.Run("negative_3", func(t *testing.T) {
		y, err := coder.Decode([]byte(base64.RawURLEncoding.EncodeToString([]byte("test")) + ".***"))

		assert.True(t, err != nil)
		assert.True(t, y == nil)
	})

	t.Run("negative_4", func(t *testing.T) {
		s := "***"
		sign := ed25519.Sign(coder.privateKey, []byte(s))
		y, err := coder.Decode([]byte(s + "." + base64.RawURLEncoding.EncodeToString(sign)))

		assert.True(t, err != nil)
		assert.True(t, y == nil)
	})
}

func Benchmark_Encode(b *testing.B) {
	coder, _ := Generate()

	for i := 0; i < b.N; i++ {
		_ = coder.Encode([]byte("123456789123456789123456789"))
	}
}

func Benchmark_Decode(b *testing.B) {
	coder, _ := Generate()
	bytes := coder.Encode([]byte("123456789123456789123456789"))

	for i := 0; i < b.N; i++ {
		_, _ = coder.Decode(bytes)
	}
}
