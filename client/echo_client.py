#!/usr/bin/env python3

import pysodium
import socket
import reprlib

def udpServe(key=b'this is my key value!',addr=('localhost',8080)) :
    s = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
    nonce = pysodium.randombytes(8)
    ciphertext = pysodium.crypto_aead_chacha20poly1305_encrypt(b'this is my key value!',None,nonce,key)
    plaintext = pysodium.crypto_aead_chacha20poly1305_decrypt(ciphertext,None,nonce,key)
    print(plaintext)
    print(ciphertext)
    print(reprlib.repr(ciphertext))
    print(nonce)
    print(ciphertext+nonce)
    s.sendto(ciphertext+nonce,addr)
    s.close()

udpServe()


