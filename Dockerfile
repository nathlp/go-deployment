FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache make gcc libc-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o main

RUN whoami
RUN id

EXPOSE 8080

CMD ["./main"]
