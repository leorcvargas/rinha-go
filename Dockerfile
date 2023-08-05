FROM golang:1.20.7-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN go build -v -o ./bin/rinha ./cmd/rinha.go

FROM alpine:3.18.2

COPY --from=builder /app/bin/rinha /app/rinha

EXPOSE 8080

CMD ["./app/rinha"]
