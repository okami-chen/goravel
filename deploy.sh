#!/bin/bash
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o clash main.go
ssh -p 10022 root@36.7.120.174 'supervisorctl stop goravel'
dir=$(pwd)
scp -P 10022 -o "LogLevel=VERBOSE" ./clash root@36.7.120.174:/root/soft/clash/v2/linux
scp -P 10022 -o "LogLevel=VERBOSE" ./storage/clash/clash_v6.yaml root@36.7.120.174:/root/soft/clash/v2/storage/clash/clash_v6.yaml
scp -P 10022 -o "LogLevel=VERBOSE" ./storage/clash/clash_v7.yaml root@36.7.120.174:/root/soft/clash/v2/storage/clash/clash_v7.yaml
ssh -p 10022 root@36.7.120.174 'supervisorctl start goravel'


