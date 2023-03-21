lint:
	go vet ./...

test:
	go test $(go list ./... | grep -v /testdata)
