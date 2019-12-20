FROM golang:1.13.0-stretch AS builder

WORKDIR /app

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

# Let's cache modules retrieval - those don't change so often
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENTRYPOINT go run main.go

EXPOSE 80
