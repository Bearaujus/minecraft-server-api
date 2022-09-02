run:
	@go mod tidy
	@go mod vendor
	@go run cmd/*.go