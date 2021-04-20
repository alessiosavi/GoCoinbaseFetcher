#!/bin/bash
sudo clearcache.sh
go build -o main
strip -s main
./main -merge true
rm data.zip

sudo clearcache.sh
7z a data.zip data
python3 -m http.server 8081
