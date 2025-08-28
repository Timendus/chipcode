package docparser

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func findInputsAndOutputs(source []string) ([]string, []string) {
	scanner := newScanner(source)

	touched := make(map[string]bool, 0) // Use a map so we get "unique" for free
	read := make(map[string]bool, 0)

	for !scanner.done() {
		switch scanner.next() {
		case ":=":
			switch scanner.peek() {
			case "random", "delay", "key":
				// Ignore
			case "hex":
				if !touched[scanner.peekN(2)] {
					read[scanner.peekN(2)] = true
				}
			default:
				if !isNumeric(scanner.peek()) && !touched[scanner.peek()] {
					read[scanner.peek()] = true
				}
			}
			touched[scanner.previous()] = true
			scanner.next()
		case "|=", "&=", "^=", "=-", "<<=", ">>=":
			if !touched[scanner.peek()] {
				read[scanner.peek()] = true
			}
			touched[scanner.previous()] = true
			touched["vF"] = true
			scanner.next()
		case "+=", "-=":
			touched[scanner.previous()] = true
			op := scanner.next()
			if !isNumeric(op) {
				// If the source of the operation is a register, mark vF as
				// destroyed. Since registers can be aliased, use "not a number"
				// as a stand-in.
				touched["vF"] = true
				if !touched[op] {
					// If the register hasn't been written to yet, it's an input
					read[op] = true
				}
			}
		case "load", "loadflags":
			if scanner.peek() == "" {
				continue // operation missing an operand, bail
			}

			// Regular load behaviour
			start := "v0"
			end := scanner.peek()

			// Check if we have an XO-CHIP range
			if scanner.peekN(2) == "-" && scanner.peekN(3) != "" {
				start = scanner.peek()
				end = scanner.peekN(3)
			}

			startIndex, startOk := register(start)
			endIndex, endOk := register(end)

			if startOk && endOk {
				// Add all registers in range individually
				if endIndex < startIndex {
					temp := endIndex
					endIndex = startIndex
					startIndex = temp
				}
				for reg := startIndex; reg <= endIndex; reg++ {
					touched[fmt.Sprintf("v%X", reg)] = true
				}
			} else {
				// Fallback for when using aliases or macro parameters
				touched[fmt.Sprintf("range %s - %s", start, end)] = true
			}

			touched["i"] = true
		case "save":
			if scanner.peek() == "" {
				continue // operation missing an operand, bail
			}

			// Regular save behaviour
			start := "v0"
			end := scanner.peek()

			// Check if we have an XO-CHIP range
			if scanner.peekN(2) == "-" && scanner.peekN(3) != "" {
				start = scanner.peek()
				end = scanner.peekN(3)
			}

			startIndex, startOk := register(start)
			endIndex, endOk := register(end)

			if startOk && endOk {
				// Add all registers in range individually
				if endIndex < startIndex {
					temp := endIndex
					endIndex = startIndex
					startIndex = temp
				}
				for reg := startIndex; reg <= endIndex; reg++ {
					regName := fmt.Sprintf("v%X", reg)
					if !touched[regName] {
						read[regName] = true
					}
				}
			} else {
				// Fallback for when using aliases or macro parameters
				read[fmt.Sprintf("range %s - %s", start, end)] = true
			}

			if !touched["i"] {
				read["i"] = true
			}
			touched["i"] = true
		case "bcd":
			op := scanner.next()
			if !touched[op] {
				read[op] = true
			}
		case "sprite":
			op1 := scanner.next()
			if !touched[op1] {
				read[op1] = true
			}
			op2 := scanner.next()
			if !touched[op2] {
				read[op2] = true
			}
			scanner.next()
		case "jump0":
			if !touched["v0"] {
				read["v0"] = true
			}
			scanner.next()
		case "if", "while":
			op1 := scanner.next()
			if !isNumeric(op1) && !touched[op1] {
				read[op1] = true
			}
			switch scanner.next() {
			case "==", "!=":
				op2 := scanner.next()
				if !isNumeric(op2) && !touched[op2] {
					read[op2] = true
				}
			case "key", "-key":
				// Ignore
			}
		}

	}

	sourceRegisters := make([]string, 0)
	for str := range read {
		if str != "" {
			sourceRegisters = append(sourceRegisters, str)
		}
	}

	targetRegisters := make([]string, 0)
	for str := range touched {
		if str != "" {
			targetRegisters = append(targetRegisters, str)
		}
	}

	slices.Sort(sourceRegisters)
	slices.Sort(targetRegisters)

	return sourceRegisters, targetRegisters
}

func isNumeric(value string) bool {
	var err error
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "0x") {
		_, err = strconv.ParseInt(value[2:], 16, 64)
	} else if strings.HasPrefix(value, "0b") {
		_, err = strconv.ParseInt(value[2:], 2, 64)
	} else {
		_, err = strconv.ParseInt(value, 10, 64)
	}
	if err != nil {
		return false
	}
	return true
}

func register(value string) (int, bool) {
	value = strings.ToLower(value)
	re := regexp.MustCompile("v[0-9a-f]")
	if !re.Match([]byte(value)) {
		return -1, false
	}
	val, _ := strconv.ParseInt(value[1:], 16, 64)
	return int(val), true
}
