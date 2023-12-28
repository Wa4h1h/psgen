build_path=./builds
windows_flags= GOARCH=amd64 GOOS=windows
darwin_flags=GOARCH=arm64 GOOS=darwin
linux_flags=GOARCH=amd64 GOOS=linux

build:
	@echo "Building binary...."
	$(windows_flags) go build -o $(build_path)/psgen_win main.go
	$(darwin_flags) go build -o $(build_path)/psgen_darwin main.go
	$(linux_flags) go build -o $(build_path)/psgen_linux main.go

.PHONY: clean
clean:
	@echo "Removing binaries...."
	@rm -rf $(build_path)

