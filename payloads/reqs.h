#ifndef REQS_H
#define REQS_H
#include "stdio.h" 

void https_post(const char *hostname, const char *port, const char *endpoint, const char *payload, char *response_out, size_t response_out_len);
void https_get(const char *hostname, const char *port, const char *endpoint, char *response_out, size_t response_out_len);

#endif