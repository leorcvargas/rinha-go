dev:
	go run ./cmd/rinha.go

build: clean deps
	CGO_ENABLED=0 go build -v -o ./bin/rinha ./cmd/rinha.go

deps:
	go mod tidy

clean:
	rm -rf ./bin

docker:
	docker build -t leorcvargas/rinha-go .

docker-push:
	docker buildx build --push --platform linux/arm/v7,linux/arm64/v8,linux/amd64 --tag leorcvargas/rinha-go .

docker-down:
	docker compose down -v --remove-orphans

docker-dev: docker-down
	docker compose -f docker-compose.dev.yml up --build

docker-local: docker-down
	docker compose -f docker-compose.local.yml up --build -d