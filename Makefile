.PHONY: build-server
build-server: proto
	cd cmd/server && go build -o server

.PHONY: build-client
build-client: proto
	cd cmd/client && go build -o client

.PHONY: linter
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test ./...

.PHONY: proto
proto:
	protoc --proto_path=internal/grpcapi/protos --go_out=internal/grpcapi/ --go_opt=paths=source_relative --go-grpc_out=internal/grpcapi/ --go-grpc_opt=paths=source_relative internal/grpcapi/protos/service.proto

.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin 