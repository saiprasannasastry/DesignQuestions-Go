GO_PACKAGES		:= $(shell go list ./... | grep -v vendor)

.PHONY: fmt
fmt:
	@go fmt $(GO_PACKAGES)
