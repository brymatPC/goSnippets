#! /bin/bash
cd src/$1
env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.compileDate=`date -u +%Y%m%d.%H%M%S` -X main.gitHash=`git rev-parse --verify HEAD` -X main.gitBranch=`git branch | grep \* | cut -d ' ' -f2` -X main.buildNumber=`git rev-list --count HEAD`"
mv $1.exe ../../winbin

