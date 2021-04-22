.PHONY: lint
lint:
ifeq (, $(shell which golangci-lint))
	$(error "No golangci-lint in $(PATH). Install it from https://github.com/golangci/golangci-lint")
endif
	golangci-lint run

.PHONY: update
update:
	go get ./... && \
	go mod tidy

.PHONY: test
test:
	go test -coverpkg=$(go list ./... | grep -v mocks | tr '\n' ',') -cover -coverprofile coverage.out ./... && \
    go tool cover -func coverage.out

.PHONY: swagger
swagger:
ifeq (, $(shell which swagger))
	$(error "No swagger in $(PATH). Install it from https://goswagger.io/install.html#homebrewlinuxbrew")
endif
	swagger generate spec -m -o ./swagger.json