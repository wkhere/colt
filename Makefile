sel=. # selection for test fuzz bench
opt=  # options for fuzz


go:
	go fmt  ./...
	go test ./...
	go install ./...

test:
	go fmt ./...
	go test -run=$(sel) ./...

fuzz:
	go test -fuzz=$(sel) $(opt) ./lex

cover:
	go test -coverprofile cov ./...
	go tool cover -html cov

bench:
	go test -bench=$(sel) -benchmem . ./...

.PHONY: go cover bench
