build:
	@echo "Building binary...."
	GOARCH=arm64 GOOS=darwin go build -o ./builds/darwin/psgen main.go
	GOARCH=amd64 GOOS=linux go build -o ./builds/linux/psgen main.go

.PHONY: clean
clean:
	@echo "Removing binaries...."
	@rm -rf builds/

