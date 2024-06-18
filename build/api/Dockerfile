FROM golang:1.22 AS builder

ENV GOPATH /go
ENV CGO_ENABLED 0
ENV GO111MODULE on
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src

COPY ./go.* ./

RUN go mod download

COPY . .

RUN go build -o main .

# FROM gcr.io/distroless/static-debian11 //TODO: fix deployed image

COPY --from=builder /go/src/main /main

EXPOSE 8080

CMD [ "main" ]
