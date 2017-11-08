go:
	go fmt
	go build
	go vet
	go test -cover
	go install

.PHONY: go
