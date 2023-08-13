FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o ./bin/rinha ./cmd/rinha.go

FROM alpine:latest

COPY --from=builder /app/bin/rinha /app/rinha

ENV GIN_MODE release

EXPOSE 8080

CMD ["/app/rinha"]
