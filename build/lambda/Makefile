build: build-custom
	env GOOS=linux go build -ldflags="-s -w" -o bin/api api/api.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/schedule schedule/schedule.go

build-custom:
	env GOOS=linux go build -ldflags="-s -w" -o bin/custom custom/custom.go

package:
	yarn package

# tidy these up - maybe use package.json scripts...
deploy: build
	liam exec -e development -r deployer -- serverless deploy

deploy-schedule: build
	liam exec -e development -r deployer -- serverless deploy function -f schedule

deploy-custom: build-custom
	liam exec -e development -r deployer -- serverless deploy function -f custom

deploy-nobuild:
	liam exec -e development -r deployer -- serverless deploy