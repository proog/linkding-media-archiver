FROM golang:alpine AS builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /ldvd ./


FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache yt-dlp

COPY --from=builder /ldvd /

ENTRYPOINT ["/ldvd"]
