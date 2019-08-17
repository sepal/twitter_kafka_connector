# Build stage
FROM golang:alpine AS build-env
LABEL maintainer="Sebastian Gilits <sep.gil@gmail.com>"
WORKDIR /src
RUN apk update && apk add --no-cache git gcc build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o twitter_kafka_connector

# Final stage
FROM golang:alpine
WORKDIR /app
COPY --from=build-env /src/twitter_kafka_connector /app/
CMD ["/app/twitter_kafka_connector"]