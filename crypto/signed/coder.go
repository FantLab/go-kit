package signed

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
)

var (
	ErrInput = errors.New("signed: invalid input")
	ErrSign  = errors.New("signed: signature check failed")
)

var signatureSize = base64.RawURLEncoding.EncodedLen(ed25519.SignatureSize)

type Key []byte

func (key Key) String() string {
	return base64.StdEncoding.EncodeToString(key)
}

type Coder struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func (c *Coder) PublicKey() Key {
	return Key(c.publicKey)
}

func (c *Coder) PrivateKey() Key {
	return Key(c.privateKey)
}

func (c *Coder) Encode(input []byte) []byte {
	size := base64.RawURLEncoding.EncodedLen(len(input))
	output := make([]byte, size+1+signatureSize)
	base64.RawURLEncoding.Encode(output, input)
	output[size] = '.'
	sign := ed25519.Sign(c.privateKey, output[:size])
	base64.RawURLEncoding.Encode(output[size+1:], sign)
	return output
}

func (c *Coder) Decode(input []byte) ([]byte, error) {
	dotIndex := bytes.IndexByte(input, '.')
	if dotIndex < 2 {
		return nil, ErrInput
	}
	sign, err := base64.RawURLEncoding.DecodeString(string(input[dotIndex+1:]))
	if err != nil {
		return nil, err
	}
	output64 := input[:dotIndex]
	if !ed25519.Verify(c.publicKey, output64, sign) {
		return nil, ErrSign
	}
	output, err := base64.RawURLEncoding.DecodeString(string(output64))
	if err != nil {
		return nil, err
	}
	return output, nil
}
