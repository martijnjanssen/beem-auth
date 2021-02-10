export GO111MODULE=on

.PHONY: dependencies
dependencies:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go@v1.25.0
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go get -u google.golang.org/grpc@v1.35.0
	go get -u github.com/fullstorydev/grpcui/cmd/grpcui/...
	go install github.com/fullstorydev/grpcui/cmd/grpcui

.PHONY: generate
generate:
# Generate messages from proto:
	protoc -I=./proto --go_out=./internal --go_opt=module=beem-auth ./proto/account-creation.proto
# Generate services
	protoc -I=./proto --go-grpc_out=./internal --go-grpc_opt=module=beem-auth ./proto/account-creation.proto

.PHONY: grpcui
grpcui:
	grpcui -plaintext localhost:5051

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

.PHONY: coverreport
coverreport:
	go tool cover -html=coverage.out

.PHONY: clean
clean:
	rm -f coverage.out
