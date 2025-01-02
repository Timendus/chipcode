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
			outputFile := testsDirectory + file.Name() + "/output.8o"
			expectedFile := testsDirectory + file.Name() + "/expected.8o"

			if !runCommand(octopus+" "+inputFile+" "+outputFile) ||
				!runCommand("git diff --no-index --color "+expectedFile+" "+outputFile) {
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
