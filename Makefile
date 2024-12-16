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

.PHONY: govulncheck
govulncheck:
	@govulncheck ./...

.PHONY: go-lint-fix-struct-staticcheck-govulncheck
go-lint-fix-struct-staticcheck-govulncheck: generate
	@go mod tidy
	@golangci-lint run ./... --fix
	@staticcheck ./...
	@fieldalignment -fix ./...
	@govulncheck -show verbose ./...

.PHONY: test
test:
	@go test -count=1 -p=8 -parallel=8 -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

.PHONY: run
run:
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) go run main.go

base_args="--showlog=true --listbank=bca,bni,mandiri,bri,danamon --from=$$(date -j -v -7d '+%Y-%m-%d') --to=$$(date -j '+%Y-%m-%d')"
process_args="process ${base_args} -g=true"
sample_args="sample ${base_args} --percentagematch=100 --amountdata=10 -g=true"

.PHONY: echo-sample-args
echo-sample-args:
	@echo $(sample_args)
	@echo $(sample_args) | pbcopy

.PHONY: run-sample
run-sample:
	@echo $(sample_args)
	@go build .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./simple-reconciliation-service  $$(echo $(sample_args))
	@#env $$(cat "params/.env" | grep -Ev '^#' | xargs) go run main.go  $$(echo $(sample_args))

.PHONY: echo-process-args
echo-process-args:
	@echo $(process_args)
	@echo $(process_args) | pbcopy

.PHONY: run-process
run-process:
	@echo "go run main.go process $(process_args)"
	@go build .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./simple-reconciliation-service $$(echo $(process_args))
#	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) go run main.go $$(echo $(process_args))

.PHONY: go-version
go-version:
	@go version

.PHONY: go-env
go-env:
	@go env

.PHONY: release-skip-publish
release-skip-publish: download install-tools generate
	@goreleaser release --skip-publish --snapshot --clean
