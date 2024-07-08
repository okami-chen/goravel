#!/bin/bash
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o clash main.go
ssh -p 10022 root@36.7.120.174 'supervisorctl stop clash:clash_00'
scp -P 10022 -o "LogLevel=VERBOSE" ./clash root@36.7.120.174:/home/data/scripts/clash/linux
scp -P 10022 -o "LogLevel=VERBOSE" ./storage/clash/clash_v6.yaml root@36.7.120.174:/home/data/scripts/clash/storage/clash/clash_v6.yaml
scp -P 10022 -o "LogLevel=VERBOSE" ./storage/clash/clash_v7.yaml root@36.7.120.174:/home/data/scripts/clash/storage/clash/clash_v7.yaml
scp -P 10022 -o "LogLevel=VERBOSE" ./storage/clash/clash_v8.yaml root@36.7.120.174:/home/data/scripts/clash/storage/clash/clash_v8.yaml
ssh -p 10022 root@36.7.120.174 'supervisorctl start clash:clash_00'

youfeng@2024


