FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

# ENV GOEXPERIMENT arenas

RUN go build -v -o ./bin/rinha ./cmd/rinha.go

FROM alpine:3.18.2

COPY --from=builder /app/bin/rinha /app/rinha

ENV GIN_MODE release

EXPOSE 8080

CMD ["/app/rinha"]
