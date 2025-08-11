package assembler

// This package is a very thin Go wrapper for the assembler in
// [c-octo](https://github.com/JohnEarnest/c-octo), an implementation of Octo in
// C by the original author of Octo. In his words, it "depends only upon the C
// standard library", so we should be able to build this for anything, anywhere.
// However, it does require the standard math library, so we do have to tell CGO
// to link that.

// #cgo LDFLAGS: -lm
// #include "octo_compiler.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func Assemble(input string) ([]byte, error) {
	C_input := C.CString(input)
	defer C.free(unsafe.Pointer(C_input))

	program := C.octo_compile_str(C_input)
	defer C.free(unsafe.Pointer(program))

	if program.is_error != 0 {
		return nil, fmt.Errorf(
			"c-octo encountered an error on line %v, column %v: %s",
			program.error_line+1,
			program.error_pos+1,
			C.GoString(&program.error[0]),
		)
	}

	rom := C.GoBytes(unsafe.Pointer(&program.rom[0]), program.length)
	return rom[0x200:], nil
}
