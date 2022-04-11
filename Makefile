.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/create-task create-task/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/list-tasks list-tasks/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/authorizer authorizer/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/tasks-data tasks-data.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
