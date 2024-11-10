ifeq ($(OS), Windows_NT)
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
else
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
endif

.DEFAULT_GOAL := run

.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

.PHONY: install-tools
install-tools: download
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

.PHONY: generate
generate:
	@wire simple-reconciliation-service/internal/inject
	@go generate ./...

.PHONY: go-lint
go-lint:
	@golangci-lint linters
	@golangci-lint run ./...

.PHONY: staticcheck
staticcheck:
	@staticcheck ./...

.PHONY: go-lint-fix-struct-staticcheck
go-lint-fix-struct-staticcheck:
	@go mod tidy
	@golangci-lint run ./... --fix
	@staticcheck ./...
	@fieldalignment -fix ./...

.PHONY: test
test:
	@go test -count=1 -p=8 -parallel=8 -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: run
run:
	@env $$(cat ".env" | grep -Ev '^#' | xargs) go run main.go

.PHONY: go-env
go-env:
	@go env

.PHONY: release-skip-publish
release-skip-publish: download install-tools generate
	@goreleaser release --skip-publish --snapshot --clean
