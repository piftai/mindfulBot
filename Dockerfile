FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o mindfulBot .

FROM alpine:latest

WORKDIR /app

# Устанавливаем tzdata для поддержки временных зон
RUN apk add --no-cache tzdata

# Устанавливаем московское время
ENV TZ=Europe/Moscow
RUN ln -fs /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone

COPY --from=builder /app/mindfulBot .

COPY .env .

CMD ["./mindfulBot"]