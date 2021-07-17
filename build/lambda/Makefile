build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/api api/api.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/schedule schedule/schedule.go

package:
	yarn package

# tidy these up - maybe use package.json scripts...
deploy: build
	liam exec -e development -r deployer -- serverless deploy

deploy-function: build
	liam exec -e development -r deployer -- serverless deploy function -f schedule

deploy-nobuild:
	liam exec -e development -r deployer -- serverless deploy