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

docker-push:
	docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag leorcvargas/rinha-go .

docker-down:
	docker compose -f docker-compose.dev.yml down -v --remove-orphans

docker-up: docker-down
	docker compose -f docker-compose.dev.yml up --build