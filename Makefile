sel=. # selection for test fuzz bench
opt=  # options for fuzz
cnt=2 # repetitions for bench

go:
	go vet  ./...
	go test ./...
	go install ./...

test:
	go vet ./...
	go test -run=$(sel) ./...

fuzz:
	go test -fuzz=$(sel) $(opt) ./lex

cover:
	go test -coverprofile cov ./...
	go tool cover -html cov

bench:
	go test -bench=$(sel) -count=$(cnt) -benchmem . ./...

.PHONY: go cover bench
