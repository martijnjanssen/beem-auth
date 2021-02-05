FROM golang:1.15-alpine AS builder

WORKDIR /beem-auth
COPY main.go go.mod go.sum /beem-auth/
COPY internal /beem-auth/internal

RUN GO111MODULE=on go get beem-auth
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o ./bin/beem-auth .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /beem-auth/bin/beem-auth .
CMD ["./beem-auth"]
