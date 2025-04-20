#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <openssl/ssl.h>
#include <openssl/err.h>

void https_post(const char *hostname, const char *port, const char *endpoint, const char *payload, char *response_out, size_t response_out_len) {
    SSL_library_init();
    SSL_load_error_strings();
    const SSL_METHOD *method = TLS_client_method();
    SSL_CTX *ctx = SSL_CTX_new(method);
    if (!ctx) return;

    struct addrinfo hints = {0}, *res;
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_STREAM;
    getaddrinfo(hostname, port, &hints, &res);

    int sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol);
    connect(sock, res->ai_addr, res->ai_addrlen);

    SSL *ssl = SSL_new(ctx);
    SSL_set_fd(ssl, sock);
    if (SSL_connect(ssl) <= 0) {
        ERR_print_errors_fp(stderr);
        return;
    }

    char request[2048];
    snprintf(request, sizeof(request),
             "POST %s HTTP/1.1\r\n"
             "Host: %s\r\n"
             "Content-Type: application/json\r\n"
             "Content-Length: %zu\r\n"
             "Connection: close\r\n\r\n"
             "%s",
             endpoint, hostname, strlen(payload), payload);

    SSL_write(ssl, request, strlen(request));

    int bytes;
    if (response_out) {
        bytes = SSL_read(ssl, response_out, response_out_len - 1);
        if (bytes > 0) response_out[bytes] = 0;
    }

    SSL_free(ssl);
    close(sock);
    SSL_CTX_free(ctx);
    freeaddrinfo(res);
}

void https_get(const char *hostname, const char *port, const char *endpoint, char *response_out, size_t response_out_len) {
    SSL_library_init();
    SSL_load_error_strings();
    const SSL_METHOD *method = TLS_client_method();
    SSL_CTX *ctx = SSL_CTX_new(method);
    if (!ctx) return;

    struct addrinfo hints = {0}, *res;
    hints.ai_family = AF_INET;
    hints.ai_socktype = SOCK_STREAM;
    getaddrinfo(hostname, port, &hints, &res);

    int sock = socket(res->ai_family, res->ai_socktype, res->ai_protocol);
    connect(sock, res->ai_addr, res->ai_addrlen);

    SSL *ssl = SSL_new(ctx);
    SSL_set_fd(ssl, sock);
    if (SSL_connect(ssl) <= 0) {
        ERR_print_errors_fp(stderr);
        return;
    }

    char request[512];
    snprintf(request, sizeof(request),
             "GET %s HTTP/1.1\r\n"
             "Host: %s\r\n"
             "Connection: close\r\n\r\n",
             endpoint, hostname);

    SSL_write(ssl, request, strlen(request));

    int bytes;
    if (response_out) {
        bytes = SSL_read(ssl, response_out, response_out_len - 1);
        if (bytes > 0) response_out[bytes] = 0;
    }

    SSL_free(ssl);
    close(sock);
    SSL_CTX_free(ctx);
    freeaddrinfo(res);
}