package emulator

func Assert(assertion bool, message string) {
	if !assertion {
		Fail(message)
	}
}

func Fail(message string) {
	panic(message)
}
