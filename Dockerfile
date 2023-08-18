FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

ENV GOEXPERIMENT arenas

RUN go build -v -o ./bin/rinha ./cmd/rinha.go

FROM golang:1.21

RUN apt-get update && apt-get install -y dumb-init

EXPOSE 8080

COPY --from=builder /app/bin/rinha .

ENV GOGC 1000

ENV GOMAXPROCS 8

CMD ["./rinha"]
