package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/timendus/chipcode/octopus2/src/assembler"
	"github.com/timendus/chipcode/octopus2/src/preprocessor"
)

var (
	USE_COLOR      = flag.Bool("color", false, "Use ANSI codes for color output to the terminal")
	DONT_USE_COLOR = flag.Bool("no-color", false, "Do not use ANSI codes for color output to the terminal")
	INPUT_FILE     = flag.String("input", "", "The path of the input file")
	OUTPUT_FILE    = flag.String("output", "STDOUT", "The path of the output file")
	USING_COLOR    = runtime.GOOS != "windows" // Disable ANSI colors by default on Windows
)

func main() {
	startTime := time.Now()

	// Parse parameters
	flag.StringVar(INPUT_FILE, "i", *INPUT_FILE, "Alias for -input")
	flag.StringVar(OUTPUT_FILE, "o", *OUTPUT_FILE, "Alias for -output")
	flag.Usage = print_usage
	flag.Parse()
	USING_COLOR = (USING_COLOR || *USE_COLOR) && !*DONT_USE_COLOR

	if *INPUT_FILE == "" {
		fmt.Fprintln(os.Stderr, bad("Input file is a required parameter\n"))
		print_usage()
		os.Exit(1)
	}

	options := map[string]bool{}
	for _, option := range flag.Args() {
		options[option] = true
	}

	if strings.ToLower(path.Ext(*INPUT_FILE)) != ".8o" {
		fmt.Fprintf(os.Stderr, bad("Don't know how to handle '%s' file\n"), path.Ext(*INPUT_FILE))
		os.Exit(1)
	}

	switch strings.ToLower(path.Ext(*OUTPUT_FILE)) {
	case ".8o", "":
		fmt.Fprintf(os.Stderr, "Octopussifying '%s' --> '%s'\n", *INPUT_FILE, *OUTPUT_FILE)

		// Octopussify the input file
		octopussified, errs := preprocessor.Octopussify(*INPUT_FILE, options)
		if len(errs) > 0 {
			fmt.Fprintln(os.Stderr, bad("Could not complete octopussification due to the following errors:"))
			for _, error := range errs {
				fmt.Fprintln(os.Stderr, "   - ", error)
			}
			os.Exit(1)
		}

		// And output to the destination file
		if *OUTPUT_FILE == "STDOUT" {
			fmt.Print(octopussified)
		} else {
			err := os.WriteFile(*OUTPUT_FILE, []byte(octopussified), 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

	case ".ch8":
		fmt.Fprintf(os.Stderr, "Octopussifying and assembling '%s' --> '%s'\n", *INPUT_FILE, *OUTPUT_FILE)

		// Octopussify the input file
		octopussified, errs := preprocessor.Octopussify(*INPUT_FILE, options)
		if len(errs) > 0 {
			fmt.Fprintln(os.Stderr, bad("Could not complete octopussification due to the following errors:"))
			for _, error := range errs {
				fmt.Fprintln(os.Stderr, "   - ", error)
			}
			os.Exit(1)
		}

		// Assemble the octopussified code
		binary, err := assembler.Assemble(octopussified)
		if err != nil {
			fmt.Fprintln(os.Stderr, bad("Could not complete assembly due to the following error:"))
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// And output to the destination file
		if *OUTPUT_FILE == "STDOUT" {
			fmt.Print(binary)
		} else {
			err = os.WriteFile(*OUTPUT_FILE, binary, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

	default:
		fmt.Fprintf(os.Stderr, bad("Don't know how to convert '%s' to '%s'\n"), path.Ext(*INPUT_FILE), path.Ext(*OUTPUT_FILE))
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, good("Finished processing in %s\n"), time.Since(startTime))
}

func good(input string) string {
	if USING_COLOR {
		return "\033[92;1m" + input + "\033[0m"
	} else {
		return input
	}
}

func bad(input string) string {
	if USING_COLOR {
		return "\033[91;1m" + input + "\033[0m"
	} else {
		return input
	}
}

func print_usage() {
	fmt.Fprintln(os.Stderr, `Usage:
   octopus -i file.8o -o result.ch8 OPTION1 OPTION2

Input file should be an assembly language file in Octo syntax with the extension
".8o". Output file can have the extensions ".8o" or ".ch8" to output either the
pre-processed intermediate assembly language or the resulting binary. If you do
not specify an output file, it will dump the pre-processed assembly to standard
output.

The options you provide as additional parameters will be "true" for the
Octopussification.

For more information, see:
https://github.com/Timendus/chipcode/blob/main/octopus2/README.md

Valid parameters:`)
	flag.PrintDefaults()
}
