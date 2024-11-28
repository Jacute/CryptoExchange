FROM golang:1.22.2 AS builder

WORKDIR /cryptoexchange

# build app
COPY cmd ./cmd
COPY internal ./internal
COPY go.* ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/app ./cmd/CryptoExchange/main.go

FROM alpine:3.20.3

WORKDIR /app

COPY --from=builder /build/app app
COPY config/docker.yaml /app/config.yaml

CMD [ "/app/app", "--config", "/app/config.yaml" ]