# syntax=docker/dockerfile:1
FROM golang:1.19.4-alpine
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
COPY src ./src
RUN go build -o /app/bananascript src/main.go
RUN rm -rf /build
WORKDIR /app
CMD ["./bananascript", "src.bs"]