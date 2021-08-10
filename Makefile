test-cover: clean
	go test ./... -coverprofile coverage.html

clean:
	find . -name 'node_modules' -type d -prune -exec rm -rf '{}' +