build_path=./builds

build:
	@echo "Building binary...."
	GOARCH=amd64 GOOS=windows go build -o $(build_path)/psgen-win main.go
	GOARCH=arm64 GOOS=darwin go build -o $(build_path)/psgen-mac main.go
	GOARCH=amd64 GOOS=linux go build -o $(build_path)/psgen-linux main.go

.PHONY: clean
clean:
	@echo "Removing binaries...."
	@rm -rf $(build_path)

