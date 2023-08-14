FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -v -o ./bin/rinha ./cmd/rinha.go

FROM alpine:3.14.10

ENV GIN_MODE release

EXPOSE 8080

COPY --from=builder /app/bin/rinha .

CMD ["./rinha"]
