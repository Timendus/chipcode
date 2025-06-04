build:
	@GOOS=linux GOARCH=amd64 go build -o ./dist/linux/octopus src/*.go
	@GOOS=windows GOARCH=amd64 go build -o ./dist/windows/octopus.exe src/*.go
	@GOOS=darwin GOARCH=amd64 go build -o ./dist/macos-intel/octopus src/*.go
	@GOOS=darwin GOARCH=arm64 go build -o ./dist/macos-apple-silicon/octopus src/*.go

test:
	@go test ./...

regression-test: build
	@cd tests && go run regression-tests.go
