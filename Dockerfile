FROM golang:1.21 as builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

ENV GOEXPERIMENT arenas

RUN CGO_ENABLED=0 go build -v -o ./bin/rinha ./cmd/rinha.go

FROM alpine:3.14.10

RUN apk add dumb-init

EXPOSE 8080

COPY --from=builder /app/bin/rinha .

ENV GOGC 12300

ENV GOMAXPROCS 3

CMD /rinha
