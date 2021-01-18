export GO111MODULE=on

.PHONY: dependencies
dependencies:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: generate
generate:
# Generate messages from proto:
	protoc -I=./proto --go_out=./internal --go_opt=module=beem-auth ./proto/account-creation.proto
# Generate services
	protoc -I=./proto --go-grpc_out=./internal --go-grpc_opt=module=beem-auth ./proto/account-creation.proto
