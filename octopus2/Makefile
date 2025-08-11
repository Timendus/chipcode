build: linux windows macos

linux: src/*.go
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o ./dist/linux/octopus src/*.go
	@TARGET=./dist/linux make copy-docs

windows: src/*.go
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows-gnu" go build -o ./dist/windows/octopus.exe src/*.go
	@TARGET=./dist/windows make copy-docs

macos: src/*.go
	@CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o ./dist/macos-apple-silicon/octopus src/*.go
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o ./dist/macos-intel/octopus src/*.go
	@TARGET=./dist/macos-apple-silicon make copy-docs
	@TARGET=./dist/macos-intel make copy-docs

test:
	@go test ./...

regression-test: linux
	@cd tests && go run regression-tests.go

copy-docs:
	cp README.md ${TARGET}
	cp LICENSE ${TARGET}
	cp src/assembler/LICENSE "${TARGET}/LICENSE for c-octo"

