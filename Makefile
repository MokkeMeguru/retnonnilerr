deps:
	go install github.com/kisielk/errcheck@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

lint: deps
	go vet ./...
	errcheck ./...
	staticcheck ./...

test:
	go test $(go list ./... | grep -v /testdata)
