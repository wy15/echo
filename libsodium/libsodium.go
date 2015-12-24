/*********************************************************************************
*      File Name     :     libsodium.go
*      Created By    :     xbg@maqi.me
*      Creation Date :     [2015-12-24 14:37]
*      Last Modified :     [AUTO_UPDATE_BEFORE_SAVE]
*      Description   :
*      Copyright     :     2015 xbg@maqi.me
*      License       :     Licensed under the Apache License, Version 2.0
**********************************************************************************/
package libsodium

//#cgo pkg-config: libsodium
//#include <sodium.h>
import "C"
import (
	"crypto/rand"
	"errors"
	"fmt"
)

func init() {

	if i := C.sodium_init(); i == -1 {
		panic(fmt.Sprintf("Sodium initialization failed, error code is %d", i))
	}

}

func EncryptData(key, plaintext []byte) ([]byte, error) {
	if len(key) != C.crypto_aead_chacha20poly1305_KEYBYTES {
		return nil, errors.New("key must be 32 bytes!")
	}

	nonce := make([]byte, C.crypto_aead_chacha20poly1305_NPUBBYTES)
	_, err := rand.Read(nonce)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("generate rand error : %v", err))
	}

	ciphertext := make([]byte, len(plaintext)+C.crypto_aead_chacha20poly1305_ABYTES)
	var clen C.ulonglong
	r, err := C.crypto_aead_chacha20poly1305_encrypt((*C.uchar)(&ciphertext[0]), &clen,
		(*C.uchar)(&plaintext[0]), (C.ulonglong)(len(plaintext)), nil, 0, nil, (*C.uchar)(&nonce[0]),
		(*C.uchar)(&key[0]))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("call encrypt function error : %v", err))
	}

	if r == -1 {
		return nil, errors.New("encrypt data error")
	}

	ciphertextaddnonce := make([]byte, clen+C.crypto_aead_chacha20poly1305_NPUBBYTES)
	copy(ciphertextaddnonce, ciphertext[:clen])
	copy(ciphertextaddnonce[clen:], nonce)
	return ciphertextaddnonce, nil
}

func DecryptData(key, ciphertextaddnonce []byte) ([]byte, error) {
	if len(key) != C.crypto_aead_chacha20poly1305_KEYBYTES {
		return nil, errors.New("key must be 32 bytes!")
	}

	clen := len(ciphertextaddnonce) - C.crypto_aead_chacha20poly1305_NPUBBYTES
	if clen-C.crypto_aead_chacha20poly1305_ABYTES <= 0 {
		return nil, errors.New("Illegal ciphertext")
	}
	nonce := ciphertextaddnonce[clen:]
	ciphertext := ciphertextaddnonce[:clen]
	plaintext := make([]byte, clen-C.crypto_aead_chacha20poly1305_ABYTES)
	var mlen C.ulonglong
	r, err := C.crypto_aead_chacha20poly1305_decrypt((*C.uchar)(&plaintext[0]), &mlen, nil,
		(*C.uchar)(&ciphertext[0]), (C.ulonglong)(clen), nil, 0, (*C.uchar)(&nonce[0]),
		(*C.uchar)(&key[0]))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("call decrypt function error : %v", err))
	}

	if r != 0 {
		return nil, errors.New("decrypt data error")
	}
	return plaintext[:mlen], nil
}
