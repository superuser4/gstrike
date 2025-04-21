#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <openssl/ssl.h>
#include <openssl/err.h>


void compute_hmac(const char *key, const char *data, char *out_hex, size_t out_len) {
    unsigned char *digest;
    unsigned int digest_len;

    digest = HMAC(EVP_sha256(), key, strlen(key), (unsigned char *)data, strlen(data), NULL, &digest_len);

    // Convert to hex
    for (unsigned int i = 0; i < digest_len && i * 2 + 1 < out_len; ++i) {
        sprintf(&out_hex[i * 2], "%02x", digest[i]);
    }
    out_hex[digest_len * 2] = '\0';
}

void https_post(const char *ip, const char *port, const char *endpoint, const char *payload, char *response_out, size_t response_out_len) {
    SSL_library_init();
    SSL_load_error_strings();
    OpenSSL_add_all_algorithms();
    const SSL_METHOD *method = TLS_client_method();
    SSL_CTX *ctx = SSL_CTX_new(method);
    if (!ctx) return;

    int sock = socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in server_addr = {0};
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(atoi(port));
    inet_pton(AF_INET, ip, &server_addr.sin_addr);

    if (connect(sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) != 0) {
        perror("Connection failed");
        SSL_CTX_free(ctx);
        return;
    }

    SSL *ssl = SSL_new(ctx);
    SSL_set_fd(ssl, sock);
    if (SSL_connect(ssl) <= 0) {
        ERR_print_errors_fp(stderr);
        SSL_free(ssl);
        close(sock);
        SSL_CTX_free(ctx);
        return;
    }

    const char sharedSecret[] = "378432999013382759857861340953603067";
    char hmac_hex[65];
    compute_hmac(sharedSecret, payload, hmac_hex, sizeof(hmac_hex));

    char request[4096];
    snprintf(request, sizeof(request),
        "POST %s HTTP/1.1\r\n"
        "Host: %s\r\n"
        "Content-Type: application/json\r\n"
        "Content-Length: %zu\r\n"
        "X-Agent-Signature: %s\r\n"
        "Connection: close\r\n\r\n"
        "%s",
        endpoint, ip, strlen(payload), hmac_hex, payload);

    
    SSL_write(ssl, request, strlen(request));

    if (response_out) {
        size_t total = 0;
        int bytes;
        while ((bytes = SSL_read(ssl, response_out + total, response_out_len - total - 1)) > 0) {
            total += bytes;
            if (total >= response_out_len - 1) break;
        }
        response_out[total] = '\0';
    }

    SSL_shutdown(ssl);
    SSL_free(ssl);
    close(sock);
    SSL_CTX_free(ctx);
}

void https_get(const char *ip, const char *port, const char *endpoint, char *response_out, size_t response_out_len) {
    SSL_library_init();
    SSL_load_error_strings();
    OpenSSL_add_all_algorithms();
    const SSL_METHOD *method = TLS_client_method();
    SSL_CTX *ctx = SSL_CTX_new(method);
    if (!ctx) return;

    int sock = socket(AF_INET, SOCK_STREAM, 0);
    struct sockaddr_in server_addr = {0};
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(atoi(port));
    inet_pton(AF_INET, ip, &server_addr.sin_addr);

    if (connect(sock, (struct sockaddr *)&server_addr, sizeof(server_addr)) != 0) {
        perror("Connection failed");
        SSL_CTX_free(ctx);
        return;
    }

    SSL *ssl = SSL_new(ctx);
    SSL_set_fd(ssl, sock);
    if (SSL_connect(ssl) <= 0) {
        ERR_print_errors_fp(stderr);
        SSL_free(ssl);
        close(sock);
        SSL_CTX_free(ctx);
        return;
    }

    char request[512];
    snprintf(request, sizeof(request),
        "GET %s HTTP/1.1\r\n"
        "Host: %s\r\n"
        "Connection: close\r\n\r\n",
        endpoint, ip);

    SSL_write(ssl, request, strlen(request));

    if (response_out) {
        size_t total = 0;
        int bytes;
        while ((bytes = SSL_read(ssl, response_out + total, response_out_len - total - 1)) > 0) {
            total += bytes;
            if (total >= response_out_len - 1) break;
        }
        response_out[total] = '\0';
    }

    SSL_shutdown(ssl);
    SSL_free(ssl);
    close(sock);
    SSL_CTX_free(ctx);
}
