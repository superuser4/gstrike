#!/bin/bash
go build -o build/server ./pkg
gcc -Wall -Wextra -Wpedantic -o build/agent payloads/* -lssl -lcrypto 
cd static && npm run build && cd ..
