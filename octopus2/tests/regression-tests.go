package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var testsDirectory = "./"

func main() {
	/* If the user has specified a specific test to run */

	if len(os.Args) > 1 {
		success := runTest(os.Args[1])
		if success {
			fmt.Println("\033[92;1m ✔️ All good!\033[0m")
			os.Exit(0)
		} else {
			fmt.Println("\033[91;1m ❌ Did not result in the expected output\033[0m")
			os.Exit(1)
		}
	}

	/* Otherwise, run all tests */

	dir, err := os.ReadDir(testsDirectory)
	if err != nil {
		panic(err)
	}

	allGood := true
	for _, file := range dir {
		if file.IsDir() {
			if !runTest(file.Name()) {
				allGood = false
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

func runTest(test string) bool {
	octopus := "../dist/linux/octopus"
	inputFile := testsDirectory + test + "/src/index.8o"
	result := true

	/* Test Octopussification */

	outputFile := testsDirectory + test + "/output.8o"
	expectedFile := testsDirectory + test + "/expected.8o"

	if !runCommand(octopus + " " + inputFile + " " + outputFile) {
		fmt.Println("Could not Octopussify " + inputFile)
		return false
	}

	if !runCommand("git diff --no-index --color " + expectedFile + " " + outputFile) {
		fmt.Println("No match for output " + expectedFile + " and " + outputFile)
		result = false
	}

	/* Test Octopussification + assembly */

	outputBinary := testsDirectory + test + "/output.ch8"
	expectedBinary := testsDirectory + test + "/expected.ch8"
	outputHex := testsDirectory + test + "/output.hex"
	expectedHex := testsDirectory + test + "/expected.hex"

	if !runCommand(octopus + " " + inputFile + " " + outputBinary) {
		fmt.Println("Could not assemble " + inputFile)
		return false
	}

	hexDump(expectedBinary, expectedHex)
	hexDump(outputBinary, outputHex)
	if !runCommand("git diff --no-index --color " + expectedHex + " " + outputHex) {
		fmt.Println("No match for binaries " + expectedBinary + " and " + outputBinary)
		result = false
	}

	return result
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
