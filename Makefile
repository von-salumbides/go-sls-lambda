.PHONY: build clean deploy
function_name = $(FUNCTION_NAME)
build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/create cmd/create/create.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update cmd/update/update.go
clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
	
remove: clean build
	sls remove --verbose

deployFunc: clean build
	sls deploy -f ${function_name} --verbose 