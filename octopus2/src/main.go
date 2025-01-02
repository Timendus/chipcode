package main

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/timendus/chipcode/octopus2/src/preprocessor"
)

func main() {
	startTime := time.Now()

	// Parse parameters
	if len(os.Args) < 3 {
		fmt.Println("\033[91;1mInput and output file are required parameters\033[0m\nUsage:\n   octopus <input file> <ouput file> <option 1> <option 2> ...")
		os.Exit(1)
	}
	input := os.Args[1]
	output := os.Args[2]
	options := map[string]bool{}
	for _, option := range os.Args[3:] {
		options[option] = true
	}

	switch path.Ext(input) {
	case ".8o":
		switch path.Ext(output) {
		case ".8o":
			fmt.Printf("Octopussifying '%s' 🡆 '%s'\n", input, output)

			// Octopussify the input file
			octopussified, errs := preprocessor.Octopussify(input, options)
			if len(errs) > 0 {
				fmt.Println("\033[91;1mCould not complete octopussification due to the following errors:\033[0m")
				for _, error := range errs {
					fmt.Println("   - ", error)
				}
				os.Exit(1)
			}

			// And output to the destination file
			err := os.WriteFile(output, []byte(octopussified), 0644)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		case ".ch8":
			fmt.Printf("Assembling '%s' 🡆 '%s'\n", input, output)

			// Octopussify the input file
			octopussified, errs := preprocessor.Octopussify(input, options)
			if len(errs) > 0 {
				fmt.Println("\033[91;1mCould not complete preprocessing due to the following errors:\033[0m")
				for _, error := range errs {
					fmt.Println("   - ", error)
				}
				os.Exit(1)
			}

			// Assemble the intermediate format to a CHIP-8 binary
			// binary, errs := assembler.Assemble(octopussified)
			// if len(errs) > 0 {
			// 	fmt.Println("\033[91;1mCould not complete assembly due to the following errors:\033[0m")
			// 	for _, error := range errs {
			// 		fmt.Println("   - ", error)
			// 	}
			// 	os.Exit(1)
			// }

			// And output to the destination file
			// err := os.WriteFile(output, binary, 0644)
			err := os.WriteFile(output, []byte(octopussified), 0644)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		default:
			fmt.Printf("\033[91;1mDon't know how to convert '%s' to '%s'\033[0m\n", path.Ext(input), path.Ext(output))
			os.Exit(1)
		}
	default:
		fmt.Printf("\033[91;1mDon't know how to convert '%s' to '%s'\033[0m\n", path.Ext(input), path.Ext(output))
		os.Exit(1)
	}

	fmt.Printf("\033[92;1mFinished processing in %s\033[0m\n", time.Since(startTime))
}
