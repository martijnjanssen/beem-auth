FROM golang:1.15-alpine AS builder

WORKDIR /beem-auth

# Do dependencies first to allow layer to be re-used
COPY go.mod go.sum /beem-auth/
RUN GO111MODULE=on go mod download

COPY main.go /beem-auth/main.go
COPY internal /beem-auth/internal

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -o ./bin/beem-auth .


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /beem-auth/bin/beem-auth .
CMD ["./beem-auth"]
