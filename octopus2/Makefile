build: linux windows macos

linux: src/*.go
	@GOOS=linux GOARCH=amd64 go build -o ./dist/linux/octopus src/*.go

windows: src/*.go
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" go build -o ./dist/windows/octopus.exe src/*.go

macos: src/*.go
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC="zig cc -target x86_64-macos -nostdlib -nodefaultlibs" go build -o ./dist/macos-intel/octopus src/*.go
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="zig cc -target aarch64-macos -nodefaultlibs" go build -o ./dist/macos-apple-silicon/octopus src/*.go

test:
	@go test ./...

regression-test: build
	@cd tests && go run regression-tests.go
