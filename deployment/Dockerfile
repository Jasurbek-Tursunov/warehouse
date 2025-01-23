FROM golang:1.23.4 as builder
LABEL authors="jasurbek"

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o build/app cmd/main.go

FROM debian:bookworm-slim

WORKDIR /root

COPY --from=builder app/build/app .

CMD ["./app"]