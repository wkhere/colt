go:
	go fmt
	go test -cover
	go install

cover:
	go test -coverprofile cov
	go tool cover -html cov

.PHONY: go
