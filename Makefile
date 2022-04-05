.PHONY: build clean deploy
function_name = $(FUNCTION_NAME)
environment = $(DEPLOY_ENV)
build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/create cmd/create/create.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/update cmd/update/update.go
clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --stage=$(environment) --verbose
	
remove: clean build
	sls remove --stage=$(environment) --verbose

deployFunc: clean build
	sls deploy --stage=$(environment) -f ${function_name} --verbose 

format: 
	gofmt -w todos/create.go
	gofmt -w todos/update.go