FROM golang:1.20.3-alpine AS build

WORKDIR /build
COPY ./ /build
RUN go build -o bananascript ./cmd/bananascript

FROM alpine
WORKDIR /app
COPY --from=build /build/bananascript .

CMD ["./bananascript"]