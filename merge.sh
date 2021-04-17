#!/bin/bash
sudo clearcache.sh
./main -merge true
rm data.zip

sudo clearcache.sh
zip -r9 data.zip data
python3 -m http.server 8081
