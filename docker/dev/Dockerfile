FROM golang:1.20-alpine as builder

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .
COPY config /usr/bin/config

RUN go build -o /usr/bin/bot

FROM docker.io/alpine:3.16

RUN apk add --no-cache su-exec ca-certificates olm bash tzdata

COPY --from=builder /usr/bin/bot /usr/bin/bot
COPY --from=builder /usr/bin/config config 

CMD ["/usr/bin/bot"]