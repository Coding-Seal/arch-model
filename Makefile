.PHONY: lint 
lint:
	golangci-lint run --fix --disable="godox"
	wsl --fix ./...

.PHONY: clean 
clean:
	rm .journal/* && rm .logs/*

.PHONY: run 
run:
	go run cmd/server/*