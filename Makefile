dev:
	go run ./cmd/rinha.go

build: clean deps
	go build -o ./bin/rinha ./cmd/rinha.go

deps:
	go mod tidy

clean:
	rm -rf ./bin

docker:
	docker build -t leorcvargas/rinha-go .

docker-down:
	docker-compose down -v --remove-orphans

docker-up: docker docker-down
	docker-compose up --build