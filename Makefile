go:
	go fmt
	go build
	go test -cover
	go install

.PHONY: go
