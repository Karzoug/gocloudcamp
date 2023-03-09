.PHONY: build-server
build-server:
	cd cmd/server && go build -o server

.PHONY: linter
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test ./...

.PHONY: .install-linter
.install-linter:
	### INSTALL GOLANGCI-LINT ###
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin 