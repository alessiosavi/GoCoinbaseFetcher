#!/bin/bash
go build -o main
strip -s main
sudo clearcache.sh
./main
