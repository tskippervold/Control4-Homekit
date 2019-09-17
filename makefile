compile:
	@echo "Compiling for OSX and Linux..."
	GOOS=darwin GOARCH=amd64 go build -v -o bin/osx-hap-bridge ./main.go
	GOOS=linux GOARCH=386 go build -o bin/linux-hap-bridge ./main.go
