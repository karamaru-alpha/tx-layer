include .env
export

.PHONY: run-context-pattern
run-context-pattern:
	(cd context-pattern && go run main.go)
