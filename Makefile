include .env
export

.PHONY: run-context-pattern
run-context-pattern:
	(cd context-pattern && go run main.go)

.PHONY: run-di-pattern
run-di-pattern:
	(cd di-pattern && go run main.go)
