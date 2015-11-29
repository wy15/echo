#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <netdb.h>
#include <sodium.h>

void error(const char *msg)
{
    perror(msg);
    exit(0);
}

int main(int argc, char *argv[])
{
    int sockfd, portno, n;
    struct sockaddr_in serv_addr;
    struct hostent *server;
    char *mykey;
    unsigned char nonce[crypto_aead_chacha20poly1305_NPUBBYTES];
    unsigned char key[crypto_aead_chacha20poly1305_KEYBYTES];
    unsigned long long ciphertext_len;
    unsigned char ciphertext[20+crypto_aead_chacha20poly1305_ABYTES];
    /*char noncehex[crypto_aead_chacha20poly1305_NPUBBYTES*2+1];*/
    /*char ciphertexthex[(20+crypto_aead_chacha20poly1305_ABYTES)*2+1];*/
    char buffer[256];
    if (argc < 4) {
       fprintf(stderr,"usage %s hostname port key\n", argv[0]);
       exit(0);
    }
    portno = atoi(argv[2]);
    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0)
        error("ERROR opening socket");
    server = gethostbyname(argv[1]);
    if (server == NULL) {
        fprintf(stderr,"ERROR, no such host\n");
        exit(0);
    }
    bzero((char *) &serv_addr, sizeof(serv_addr));
    serv_addr.sin_family = AF_INET;
    bcopy((char *)server->h_addr,
         (char *)&serv_addr.sin_addr.s_addr,
         server->h_length);
    serv_addr.sin_port = htons(portno);
    if (connect(sockfd,(struct sockaddr *) &serv_addr,sizeof(serv_addr)) < 0)
        error("ERROR connecting");
    bzero(buffer,256);
    if (sodium_init() == -1) {
         error("sodium init error");
    }
    randombytes_buf(nonce, sizeof nonce);
    /*if(sodium_bin2hex(noncehex, crypto_aead_chacha20poly1305_NPUBBYTES*2+1,nonce,sizeof nonce))*/
        /*printf("%s\n",noncehex);*/

    bzero(key,crypto_aead_chacha20poly1305_KEYBYTES);
    mykey = argv[3];
    memcpy(key, mykey, strlen(mykey));
    if(crypto_aead_chacha20poly1305_encrypt(ciphertext,&ciphertext_len,"this is my c text",strlen("this is my c text"),NULL,0,NULL,nonce,key) == -1) {
         error("encrypt error");
    }
    /*if(sodium_bin2hex(ciphertexthex,(20+crypto_aead_chacha20poly1305_ABYTES)*2+1,ciphertext,ciphertext_len))*/
        /*printf("%s\n",ciphertexthex);*/
    memcpy(buffer, ciphertext, ciphertext_len);
    memcpy(buffer+ciphertext_len, nonce, crypto_aead_chacha20poly1305_NPUBBYTES);
    n = write(sockfd,buffer,strlen(buffer));
    if (n < 0)
         error("ERROR writing to socket");
    bzero(buffer,256);
    n = read(sockfd,buffer,255);
    if (n < 0)
         error("ERROR reading from socket");
    printf("%s\n",buffer);
    close(sockfd);
    return 0;
}
