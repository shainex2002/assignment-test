FROM golang:1.10 as builder
WORKDIR /go/src/github.com/kousik93/simple-dump-server

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/simple-dump-server . && \
    go test ./... -cover

FROM alpine:3.6
RUN mkdir /opt
RUN apk add --no-cache ca-certificates
WORKDIR /opt/
COPY --from=builder /go/src/github.com/kousik93/simple-dump-server/bin/simple-dump-server .
ENTRYPOINT ["./simple-dump-server"]
