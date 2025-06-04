package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	testsDirectory := "./"
	octopus := "../dist/linux/octopus"

	dir, err := os.ReadDir(testsDirectory)
	if err != nil {
		panic(err)
	}

	allGood := true
	for _, file := range dir {
		if file.IsDir() {
			inputFile := testsDirectory + file.Name() + "/src/index.8o"

			/* Test Octopussification */

			outputFile := testsDirectory + file.Name() + "/output.8o"
			expectedFile := testsDirectory + file.Name() + "/expected.8o"

			if !runCommand(octopus + " " + inputFile + " " + outputFile) {
				fmt.Println("Could not Octopussify " + inputFile)
				allGood = false
				continue
			}

			if !runCommand("git diff --no-index --color " + expectedFile + " " + outputFile) {
				fmt.Println("No match for output " + expectedFile + " and " + outputFile)
				allGood = false
				continue
			}

			/* Test Octopussification + assembly */

			outputBinary := testsDirectory + file.Name() + "/output.ch8"
			expectedBinary := testsDirectory + file.Name() + "/expected.ch8"
			outputHex := testsDirectory + file.Name() + "/output.hex"
			expectedHex := testsDirectory + file.Name() + "/expected.hex"

			if !runCommand(octopus + " " + inputFile + " " + outputBinary) {
				fmt.Println("Could not assemble " + inputFile)
				allGood = false
				continue
			}

			hexDump(expectedBinary, expectedHex)
			hexDump(outputBinary, outputHex)
			if !runCommand("git diff --no-index --color " + expectedHex + " " + outputHex) {
				fmt.Println("No match for binaries " + expectedBinary + " and " + outputBinary)
				allGood = false
				continue
			}
		}
	}

	if allGood {
		fmt.Println("\033[92;1m ✔️ All good!\033[0m")
		os.Exit(0)
	} else {
		fmt.Println("\033[91;1m ❌ One or more projects did not result in the expected output\033[0m")
		os.Exit(1)
	}
}

func runCommand(command string) bool {
	parts := strings.Split(command, " ")
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
		return false
	}
	return true
}

func hexDump(inputFile string, outputFile string) error {
	cmd := exec.Command("hexdump", "-v", "-C", inputFile)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	os.WriteFile(outputFile, output, 0644)
	return nil
}
