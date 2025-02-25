FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o mindfulBot .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/mindfulBot .

COPY .env .

CMD ["./mindfulBot"]