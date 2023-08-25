sel=. # selection for test fuzz bench
opt=  # options for fuzz
cnt=2 # repetitions for bench

target=./cmd/...

default: test

build:
	go build $(target)

install: test
	go install $(target)

test:
	go test -run=$(sel) ./...

fuzz:
	go test -fuzz=$(sel) $(opt) ./lex

cover:
	go test -coverprofile cov ./...
	go tool cover -html cov

bench:
	go test -bench=$(sel) -count=$(cnt) -benchmem . ./...

.PHONY: go install test fuzz cover bench
