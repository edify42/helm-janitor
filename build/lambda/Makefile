build: build-custom
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/api/bootstrap api/api.go
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/schedule/bootstrap schedule/schedule.go

build-custom:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o bin/custom/bootstrap custom/custom.go

zip: build
	mkdir -p bin/lambda
	zip -j bin/lambda/api.zip bin/api/bootstrap
	zip -j bin/lambda/schedule.zip bin/schedule/bootstrap
	zip -j bin/lambda/custom.zip bin/custom/bootstrap

package:
	yarn package

# tidy these up - maybe use package.json scripts...
deploy: zip
	liam exec -e development -r platform-poweruser -- serverless deploy

deploy-schedule: build
	liam exec -e development -r deployer -- serverless deploy function -f schedule

deploy-custom: build-custom
	liam exec -e development -r deployer -- serverless deploy function -f custom

deploy-nobuild:
	liam exec -e development -r deployer -- serverless deploy