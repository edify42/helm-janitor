test-cover: clean
	go test ./... -coverprofile coverage.html

lint:
	go install github.com/mgechev/revive@latest
	revive ./...

vet:
	go vet ./...

clean:
	find . -name 'node_modules' -type d -prune -exec rm -rf '{}' +