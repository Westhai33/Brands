FROM golang:1.22.3-alpine

WORKDIR /app

RUN apk update && apk upgrade && apk add --update alpine-sdk bash git openssh make cmake && \
    rm -rf /var/cache/apk/*

EXPOSE 8080
