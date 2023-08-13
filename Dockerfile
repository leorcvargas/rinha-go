FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . /app/

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -v -o /app/bin/rinha /app/cmd/rinha.go

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  ca-certificates && \
  rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/bin/rinha /app/bin/rinha

ENV GIN_MODE release

EXPOSE 8080

CMD ["/app/bin/rinha"]
