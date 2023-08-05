dev:
	go run ./cmd/rinha.go

build: clean deps
	go build -o ./bin/rinha ./cmd/rinha.go

deps:
	go mod tidy

clean:
	rm -rf ./bin