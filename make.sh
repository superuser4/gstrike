#!/bin/bash
go build -o build/server ./pkg
./gstrike-beacon/make.sh
