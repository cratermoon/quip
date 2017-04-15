build:
	go install

build-prod:
	env GOOS=linux GOARCH=amd64 go build -ldflags  "-w -s -X main.Version=`git log --pretty=format:'%H' -n 1` -X main.Build=`date +%FT%T%z`"
run:
	go install && quipserver

.PHONY: run test
