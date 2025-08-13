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
	"github.com/timendus/chipcode/octopus2/src/emulator"
	"github.com/timendus/chipcode/octopus2/src/preprocessor"
)

var (
	USE_COLOR      = flag.Bool("color", false, "Use ANSI codes for color output to the terminal")
	DONT_USE_COLOR = flag.Bool("no-color", false, "Do not use ANSI codes for color output to the terminal")
	INPUT_FILE     = flag.String("input", "", "The path of the input file")
	OUTPUT_FILE    = flag.String("output", "STDOUT", "The path of the output file")
	EMULATE        = flag.String("run", "disabled", "Run the given code or binary in the embedded emulator instead")
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

	// Figure out what to do given the input and output file and the parameters
	input_ext := strings.ToLower(path.Ext(*INPUT_FILE))
	output_ext := strings.ToLower(path.Ext(*OUTPUT_FILE))
	do_preprocessing := input_ext == ".8o"
	build_binary := input_ext == ".8o" && (output_ext == ".ch8" || *EMULATE != "disabled")
	read_binary := input_ext == ".ch8" && *EMULATE != "disabled"

	// Do we have something to do?
	if !(do_preprocessing || build_binary || read_binary) {
		fmt.Fprintf(os.Stderr, bad("Don't know what to do with '%s' file\n"), path.Ext(*INPUT_FILE))
		os.Exit(1)
	}

	// Do preprocessing if we have received a .8o file
	var preprocessed string
	if do_preprocessing {
		fmt.Fprintf(os.Stderr, "Octopussifying '%s'...\n", *INPUT_FILE)
		var errs []error
		preprocessed, errs = preprocessor.Octopussify(*INPUT_FILE, options)
		if len(errs) > 0 {
			fmt.Fprintln(os.Stderr, bad("Could not complete pre-processing due to the following errors:"))
			for _, error := range errs {
				fmt.Fprintln(os.Stderr, "   - ", error)
			}
			os.Exit(1)
		}
	}

	// Build a binary from it if we're transforming a .8o file into a .ch8 file or we're emulating the thing
	var binary []byte
	if build_binary {
		fmt.Fprintf(os.Stderr, "Assembling '%s'...\n", *INPUT_FILE)
		var err error
		binary, err = assembler.Assemble(preprocessed)
		if err != nil {
			fmt.Fprintln(os.Stderr, bad("Could not complete assembly due to the following error:"))
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if *OUTPUT_FILE != "STDOUT" {
		var output []byte
		if output_ext == ".ch8" {
			output = binary
		} else {
			output = []byte(preprocessed)
		}
		err := os.WriteFile(*OUTPUT_FILE, output, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if do_preprocessing || build_binary {
		fmt.Fprintf(os.Stderr, good("Finished processing in %s\n"), time.Since(startTime))
	}

	if read_binary {
		fmt.Fprintf(os.Stderr, "Reading ROM from file '%s'...\n", *INPUT_FILE)
		var err error
		binary, err = os.ReadFile(*INPUT_FILE)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if *EMULATE != "disabled" {
		fmt.Fprintf(os.Stderr, "Running emulation sequence...\n")
		err := emulator.Emulate(binary, *EMULATE)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	} else if *OUTPUT_FILE == "STDOUT" {
		fmt.Println(preprocessed)
	}

	fmt.Fprint(os.Stderr, good("Done\n"))
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
