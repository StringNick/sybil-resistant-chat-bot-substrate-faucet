FROM golang:1.20-alpine as tester

RUN apk add --no-cache git ca-certificates build-base su-exec olm-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .
COPY config /usr/bin/config

ENTRYPOINT ["go", "test", "-v", "./..."]