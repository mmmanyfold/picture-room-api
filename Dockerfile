## CACHE DEPENDENCIES IN BASE CONTAINER
FROM golang:1.15 as base

WORKDIR /go/src/github.com/mmmanyfold/picture-room-api

ENV GO111MODULE=on
ENV GOPATH=/go

COPY go.mod .
COPY go.sum .

RUN go mod download

## BUILD SERVER BINARY IN BUILDR CONTAINER
FROM base as builder

COPY . .

WORKDIR ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/api

## CREATE FINAL IMAGE
FROM alpine:3.9

RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN mkdir -p /home/app

## ADD NON ROOT USER
RUN addgroup -S app && adduser app -S -G app
RUN chown app /home/app

WORKDIR /home/app

COPY --from=builder /go/bin/api ./api
RUN chmod +x ./api

# copy config
COPY --from=builder /go/src/github.com/mmmanyfold/picture-room-api/config/prod.yaml ./config/prod.yaml

USER app

ARG PORT
ENV PORT ${PORT:-8080}

EXPOSE 8080

CMD ["./api"]
