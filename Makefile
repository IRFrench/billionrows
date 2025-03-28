DEFAULT: help

help: ## Show commands of the makefile (and any included files)
	@awk 'BEGIN {FS = ":.*?## "}; /^[0-9a-zA-Z_.-]+:.*?## .*/ {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

probe: ## Run the probe test to get metadata about the file
	CGOENABLED=0 go run cmd/probing/main.go

run: ## Run the attempt
	CGOENABLED=0 go run cmd/challenge/main.go

test: ## Run the tests!
	CGOENABLED=0 go test ./attempt1/... -v -cover

benchmark: ## Run the benchmarks!
	CGOENABLED=0 go test -bench=. ./attempt1 -benchmem