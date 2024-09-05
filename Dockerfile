FROM golang:1.21-alpine

RUN apk update

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

EXPOSE 50501

RUN go build -o /sdcc_registry

ENTRYPOINT ["/sdcc_registry"]