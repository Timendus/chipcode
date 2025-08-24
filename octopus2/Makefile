build: linux windows macos stdlib-docs
test: go-tests regression-test stdlib-tests

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

copy-docs:
	@cp README.md ${TARGET}
	@cp LICENSE ${TARGET}
	@cp src/assembler/LICENSE "${TARGET}/LICENSE for c-octo"

go-tests:
	@go test ./...

regression-test: linux
	@cd tests && go run regression-tests.go


# Run tests and build documentation for std library

LIB_DIR := library/lib
LIB_DOCS := library/docs
LIB_TESTS := library/tests
LIB_RESULTS := library/.test
OCTOPUS := dist/linux/octopus

LIBS := $(wildcard $(LIB_DIR)/*.8o)
DOCS := $(patsubst $(LIB_DIR)/%.8o,$(LIB_DOCS)/%.md,$(LIBS))
TESTS := $(patsubst $(LIB_TESTS)/%.8o,$(LIB_RESULTS)/%.ok,$(wildcard $(LIB_TESTS)/*.8o))

$(LIB_DOCS)/%.md: $(LIB_DIR)/%.8o
	$(OCTOPUS) -i $< -o $@

stdlib-docs: linux $(DOCS)

$(LIB_RESULTS):
	@mkdir -p "$@"

# One stamp per input file; only created if it passes
$(LIB_RESULTS)/%.ok: $(LIB_TESTS)/%.8o | $(LIB_RESULTS)
	@out=`"$(OCTOPUS)" -i "$<" -run "mode: xochip, 1000000, load: 0x200: 1" 2> /dev/null`; \
	if [ "$$out" == "0200: 42" ]; then \
	  printf 'PASS %s -> %s\n' "$<" "$$out"; \
	  printf '%s\n' "$$out" > "$@"; \
	else \
	  printf 'FAIL %s -> %s (expected "0200: 42")\n' "$<" "$$out" >&2; \
	  rm -f "$@"; \
	  exit 1; \
	fi

stdlib-tests: linux $(TESTS)
	@echo "✅ All $(words $(TESTS)) tests passed."
