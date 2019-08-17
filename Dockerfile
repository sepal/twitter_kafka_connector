# Build stage
FROM golang:latest AS build-env
LABEL maintainer="Sebastian Gilits <sep.gil@gmail.com>"
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o twitter_kafka_connector

# Final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/twitter_kafka_connector /app/