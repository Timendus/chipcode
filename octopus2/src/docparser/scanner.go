package docparser

import "strings"

type scanner struct {
	tokens  *[]string
	pointer int
}

func newScanner(source []string) scanner {
	tokens := make([]string, 0)
	for _, line := range source {
		// Make sure we ignore comments
		line = strings.Split(line, "#")[0]
		// Split on whitespace
		tokens = append(tokens, strings.Fields(line)...)
	}
	return scanner{
		tokens:  &tokens,
		pointer: 0,
	}
}

func (s *scanner) peek() string {
	return s.peekN(1)
}

func (s *scanner) peekN(dist int) string {
	if s.pointer+dist-1 >= len(*s.tokens) {
		return ""
	}
	return (*s.tokens)[s.pointer+dist-1]
}

func (s *scanner) next() string {
	if s.done() {
		return ""
	}
	token := (*s.tokens)[s.pointer]
	s.pointer += 1
	return token
}

func (s *scanner) previous() string {
	if s.pointer <= 1 {
		return ""
	}
	return (*s.tokens)[s.pointer-2]
}

func (s *scanner) done() bool {
	return s.pointer+1 >= len(*s.tokens)
}
