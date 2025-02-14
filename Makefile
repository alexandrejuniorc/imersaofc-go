### Build methods for the project
build/dev:
	@echo "Building the project..."
	@go build .

build/executable/windows:
	@echo "Building the project..."
	@GOOS=windows go build .