#!/usr/bin/env python3

#import pysodium
import socket

def udpServe(key='this is my key value!',addr=('localhost',8080)) :
    s = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
    s.sendto(b"test",addr)
    s.close()

udpServe()


