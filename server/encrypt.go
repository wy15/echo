package server

import (
	"crypto/rand"
	"errors"

	"github.com/codahale/chacha20poly1305"
)

/*
 *use chacha20poly1305 implements the AEAD_CHACHA20_POLY1305 algorithm.
 *For full document, go to "https://godoc.org/github.com/codahale/chacha20poly1305"
 *key size must be 32
 */
func EncryptData(key, plaintext, addata []byte) (edata []byte, err error) {
	if len(key) != 32 {
		return nil, errors.New("key size must be 32 byte")
	}
	c, e := chacha20poly1305.New(key)
	if e != nil {
		return nil, e
	}

	nonce, err := generateNonce(c.NonceSize())
	if err != nil {
		return
	}

	//EncryptData
	ciphertext := c.Seal(nil, nonce, plaintext, addata)

	rdata := make([]byte, len(ciphertext)+c.NonceSize())
	//Append nonce to the ciphertext
	copy(rdata, ciphertext)
	copy(rdata[len(ciphertext):], nonce[:])

	return rdata, nil
}

//Decrypt data function
func DecryptData(key, ciphertext, addata []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("key size must be 32 bytes")
	}
	c, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	nonce := ciphertext[len(ciphertext)-c.NonceSize():]

	p, err := c.Open(nil, nonce, ciphertext[:len(ciphertext)-c.NonceSize()], addata)

	return p, err
}

/*
 * generate nonce for chacha20poly1305
 * nonce size is 8
 */
func generateNonce(nonceSize int) (nonce []byte, err error) {
	nonce = make([]byte, nonceSize)
	_, err = rand.Read(nonce)
	return
}
