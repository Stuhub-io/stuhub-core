FROM golang:1.22-alpine

ENV GOPATH /go
ENV CGO_ENABLED 0

RUN apk update && \
    apk add git && \
    apk add bash

WORKDIR /app

# Install go packages
COPY go.mod go.sum ./
RUN go mod download

# Install air for live reload
RUN go install github.com/air-verse/air@v1.52.2

COPY ./build/local/api/entrypoint.sh /entrypoint.sh
RUN sed -i 's/\r$//g' /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY ./build/local/api/start.sh /start.sh
RUN sed -i 's/\r$//g' /start.sh
RUN chmod +x /start.sh

COPY ./build/scripts/wait-for-it.sh /wait-for-it.sh
RUN sed -i 's/\r$//g' /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Source code
COPY . /app

ENTRYPOINT [ "/entrypoint.sh" ]

