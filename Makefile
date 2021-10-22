test-cover: clean
	go test ./... -coverprofile coverage.html

deps:
	go get -u golang.org/x/lint/golint

lint: deps
	goling ./...

vet:
	go vet ./...

clean:
	find . -name 'node_modules' -type d -prune -exec rm -rf '{}' +