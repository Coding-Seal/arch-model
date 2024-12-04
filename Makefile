.PHONY: lint 
lint:
	wsl --fix ./...
	golangci-lint run --fix

.PHONY: clean 
clean:
	rm .journal/* && rm .logs/*

.PHONY: run 
run:
	go run cmd/server/*