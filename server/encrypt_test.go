package server_test

import "testing"

func TestEncryptAndDecryptData(t *testing.T) {
	key := make([]byte, 32)
	copy(key, []byte("this is my key!"))
	plaintext := "yay for me"
	addata := make([]byte, 10)
	copy(addata, []byte("ok"))

	//Test key length.Key must be 32 bytes.
	if _, err := server.EncryptData(key[1:], []byte(plaintext), addata); err == nil {
		t.Errorf("EncryptData error: %s", "key is not 32 length, this should be error, but not!")
	}

	//Test EncryptData main function.
	ciphertext, err := server.EncryptData(key, []byte(plaintext), addata)
	if err != nil {
		t.Errorf("EncryptData error: %v", err)
	} else {
		t.Logf("ciphertext is %x", ciphertext)
	}

	//Check decrypt
	p, err := server.DecryptData(key, ciphertext, addata)
	//c, err := chacha20poly1305.New(key)

	//p, err := c.Open(nil, ciphertext[len(ciphertext)-c.NonceSize():], ciphertext[:len(ciphertext)-c.NonceSize()], addata)
	if err != nil {
		t.Errorf("Decrypt data fail: %v", err)
	} else {
		t.Logf("plantext is %s", p)
	}

}
