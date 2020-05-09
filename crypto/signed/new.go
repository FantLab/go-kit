package signed

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io/ioutil"
)

var (
	ErrPublicKeySize  = errors.New("signed: invalid public key size")
	ErrPrivateKeySize = errors.New("signed: invalid private key size")
)

func NewFileCoder64(publicKeyFile, privateKeyFile string) (*Coder, error) {
	publicKey64, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return nil, err
	}
	privateKey64, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	return NewCoder64(string(publicKey64), string(privateKey64))
}

func NewCoder64(publicKey64, privateKey64 string) (*Coder, error) {
	publicKey, err := base64.StdEncoding.DecodeString(publicKey64)
	if err != nil {
		return nil, err
	}
	privateKey, err := base64.StdEncoding.DecodeString(privateKey64)
	if err != nil {
		return nil, err
	}
	return NewCoder(publicKey, privateKey)
}

func NewCoder(publicKey, privateKey []byte) (*Coder, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return nil, ErrPublicKeySize
	}
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, ErrPrivateKeySize
	}
	return &Coder{publicKey: publicKey, privateKey: privateKey}, nil
}

func Generate() (*Coder, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Coder{publicKey: publicKey, privateKey: privateKey}, nil
}
