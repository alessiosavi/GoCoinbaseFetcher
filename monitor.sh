#!/bin/bash
p=$(ps -ef | grep "./main" | grep -v grep | head -n1 | awk '{print $2}')

while true;
do 
  echo $(date)
  sudo cat /proc/$p/smaps  | grep -i pss |  awk '{Total+=$2} END {print Total/1024" MB"}'
  sudo cat /proc/$p/smaps  | grep -i rss |  awk '{Total+=$2} END {print Total/1024" MB"}'
  echo "=========="
  sleep 5
done
