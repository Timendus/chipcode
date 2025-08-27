package docparser

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func findTargetRegisters(source []string) []string {
	tokens := sourceToTokens(source)
	touched := make(map[string]bool, 0) // Use a map so we get "unique" for free
	var prevToken, nextToken string
	for i, token := range tokens {
		prevToken = ""
		nextToken = ""
		if i > 0 {
			prevToken = tokens[i-1]
		}
		if i+1 < len(tokens) {
			nextToken = tokens[i+1]
		}

		switch token {
		case ":=":
			touched[prevToken] = true
		case "|=", "&=", "^=", "=-", "<<=", ">>=":
			touched[prevToken] = true
			touched["vF"] = true
		case "+=", "-=":
			touched[prevToken] = true
			if !isNumeric(nextToken) {
				// If the source of the operation is a register, mark vF as
				// destroyed. Since registers can be aliased, use "not a number"
				// as a stand-in.
				touched["vF"] = true
			}
		case "load", "loadflags":
			if nextToken == "" {
				continue // operation missing an operand, bail
			}

			// Regular load behaviour
			start := "v0"
			end := nextToken

			// Check if we have an XO-CHIP range
			if i+3 < len(tokens) && tokens[i+2] == "-" {
				start = nextToken
				end = tokens[i+3]
			}

			startReg, startIsReg := register(start)
			endReg, endIsReg := register(end)

			if startIsReg == nil && endIsReg == nil {
				// Add all registers in range individually
				if endReg < startReg {
					temp := endReg
					endReg = startReg
					startReg = temp
				}
				for reg := startReg; reg <= endReg; reg++ {
					touched[fmt.Sprintf("v%X", reg)] = true
				}
			} else {
				// Fallback for when using aliases or macro parameters
				touched[fmt.Sprintf("range %s - %s", start, end)] = true
			}

			touched["i"] = true
		case "save":
			touched["i"] = true
		}

	}

	targetRegisters := make([]string, 0)
	for str := range touched {
		if str != "" {
			targetRegisters = append(targetRegisters, str)
		}
	}

	slices.Sort(targetRegisters)
	return targetRegisters
}

func sourceToTokens(source []string) []string {
	tokens := make([]string, 0)
	for _, line := range source {
		// Make sure we ignore comments
		line = strings.Split(line, "#")[0]
		// Split on whitespace
		tokens = append(tokens, strings.Fields(line)...)
	}
	return tokens
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

func register(value string) (int, error) {
	value = strings.ToLower(value)
	re := regexp.MustCompile("v[0-9a-f]")
	if !re.Match([]byte(value)) {
		return -1, fmt.Errorf("input string is not a register")
	}
	val, _ := strconv.ParseInt(value[1:], 16, 64)
	return int(val), nil
}
