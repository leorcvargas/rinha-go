dev:
	go run ./cmd/rinha.go

build: clean deps
	go build -o ./bin/rinha ./cmd/rinha.go

deps:
	go mod tidy

clean:
	rm -rf ./bin

docker-down:
	docker-compose down -v --remove-orphans

docker-up: docker-down
	docker-compose up --build