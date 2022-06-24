SHELL := /bin/bash

# ==============================================================================
# Administration
key:
	go run app/tooling/sales-admin/main.go genkey

# ==============================================================================
# Generate mocks

mocks:
	mockgen -source=business/user/usecases/repository.go -destination=business/user/usecases/mock/mock.go -package=mock

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1
	staticcheck -checks=all ./...
test-verbose:
	go test ./... -count=1 -v
	staticcheck -checks=all ./...

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Running the system within the local computer

run:
	go run app/services/sales-api/main.go