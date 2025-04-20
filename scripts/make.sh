#!/bin/bash
go build -o build/server ./pkg
gcc -Wall -Wextra -Wpedantic -static payloads/agent.c -o build/agent -lcurl
