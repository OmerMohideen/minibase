benchmark:
	@go test -bench . -benchmem

test:
	@go test -v ./...
