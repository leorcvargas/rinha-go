FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

ENV GOEXPERIMENT arenas

RUN go build -v -o ./bin/rinha ./cmd/rinha.go

# FROM alpine:3.14.10

RUN apt-get update && apt-get install dumb-init

EXPOSE 8080

# COPY --from=builder /app/bin/rinha .

ENV GOGC 1000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/app/bin/rinha"]
