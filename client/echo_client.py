#!/usr/bin/env python3

import pysodium
import socket

def udpServe(key=b'this is my key value!',addr=('localhost',8080)) :
    s = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
    nonce = pysodium.randombytes(8)
    ciphertext = pysodium.crypto_aead_chacha20poly1305_encrypt(b'this is my key value!','',nonce,key)
    plaintext = pysodium.crypto_aead_chacha20poly1305_decrypt(ciphertext,'',nonce,key)
    print(plaintext)
    print(ciphertext)
    print(len(ciphertext))
    s.sendto(ciphertext,addr)
    s.close()

udpServe()


