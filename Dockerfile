FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd

ENV  GOOS linux
ENV GOARCH amd64

RUN go build -v -o ./bin/rinha ./cmd/rinha.go

FROM alpine:3.18.2

COPY --from=builder /app/bin/rinha /app/rinha

ENV GIN_MODE release

EXPOSE 8080

CMD ["/app/rinha"]
