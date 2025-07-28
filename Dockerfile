FROM --platform=$BUILDPLATFORM golang:alpine AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /app ./cmd/main.go


FROM alpine:latest

RUN apk add --no-cache yt-dlp

COPY --from=builder /app /

ENTRYPOINT ["/app"]
