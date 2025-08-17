package preprocessor

import (
	"fmt"

	"os"
	"path"
	"strconv"
	"strings"
)

const (
	CODE = iota
	DATA
)

func Octopussify(filename string, options map[string]bool) (string, []error) {
	contents, err := loadFile(filename)
	if err != nil {
		return "", []error{err}
	}
	octopussified, errs := octopussify(contents, filename, options)
	if len(errs) > 0 {
		return "", errs
	}

	return reorder(octopussified), []error{}
}

func loadFile(filename string) (string, error) {
	if _, err := os.Stat(filename); err != nil {
		return "", fmt.Errorf("requested file '%s' not found", filename)
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %s", filename, err.Error())
	}
	return string(file), nil
}

func reorder(input string) string {
	lines := strings.Split(input, "\n")
	current := CODE
	code := []string{}
	data := []string{}
	for _, line := range lines {
		if match, param := find(line, ":segment"); match {
			if param == "code" {
				current = CODE
				continue
			}
			if param == "data" {
				current = DATA
				continue
			}
		}
		if current == CODE {
			code = append(code, line)
		}
		if current == DATA {
			data = append(data, line)
		}
	}
	return strings.Join(code, "\n") + strings.Join(data, "\n")
}

func octopussify(contents, filename string, options map[string]bool) (string, []error) {
	lines := strings.Split(contents, "\n")
	outputting := []bool{true}
	current := CODE
	output := []string{}
	errors := []error{}
	for lineNr, line := range lines {

		// Manage conditionals
		if match, params := find(line, ":if"); match {
			value, ok := options[params]
			outputting = append(outputting, ok && value)
			continue
		}
		if match, params := find(line, ":unless"); match {
			value, ok := options[params]
			outputting = append(outputting, !(ok && value))
			continue
		}
		if match, _ := find(line, ":else"); match {
			outputting[len(outputting)-1] = !outputting[len(outputting)-1]
			continue
		}
		if match, _ := find(line, ":end"); match {
			outputting = outputting[:len(outputting)-1]
			continue
		}

		// Do the conditionals say we should refrain from outputting?
		if len(outputting) > 0 && !outputting[len(outputting)-1] {
			continue
		}

		// Redefine options
		if match, params := find(line, ":const"); match {
			name, value := params[:strings.Index(params, " ")], params[strings.Index(params, " ")+1:]
			intValue, err := parseValue(value)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			options[name] = intValue != 0
		}

		// Dump options if requested
		if match, _ := find(line, ":dump-options"); match {
			fmt.Printf("Options in '%s' at line %d:\n", filename, lineNr+1)
			for name, value := range options {
				fmt.Printf("   %s: %v\n", name, value)
			}
			continue
		}

		// Manage includes
		if match, params := find(line, ":include"); match {
			start := strings.Index(params, `"`)
			end := strings.LastIndex(params, `"`)
			if start == -1 || end == -1 || start == end {
				errors = append(errors, fmt.Errorf("missing quotes in include statement in '%s' at line %d: '%s'", filename, lineNr+1, line))
				continue
			}
			subfilename := params[start+1 : end]
			if len(subfilename) < 1 {
				errors = append(errors, fmt.Errorf("invalid filename in include in '%s' at line %d: '%s'", filename, lineNr+1, line))
				continue
			}
			var file string
			var err error
			errs := []error{}
			switch path.Ext(subfilename) {
			case ".bin", ".ch8":
				file, err = loadBinaryFile(path.Join(path.Dir(filename), subfilename))
				if err != nil {
					errs = []error{err}
				}
			case ".png", ".bmp", ".jpg", ".jpeg", ".tif", ".tiff", ".gif":
				file, err = loadImageFile(path.Join(path.Dir(filename), subfilename), params[end+1:])
				if err != nil {
					errs = []error{err}
				}
			default:
				file, err = loadFile(path.Join(path.Dir(filename), subfilename))
				if err != nil {
					errs = []error{err}
				} else {
					file, errs = octopussify(file, path.Join(path.Dir(filename), subfilename), options)
				}
			}
			if err != nil {
				errors = append(errors, err)
				continue
			}
			errors = append(errors, errs...)
			output = append(output, file)
			if current == CODE {
				output = append(output, ":segment code")
			} else {
				output = append(output, ":segment data")
			}
			continue
		}

		// Keep track of segments for includes
		if match, param := find(line, ":segment"); match {
			switch param {
			case "code":
				current = CODE
			case "data":
				current = DATA
			default:
				errors = append(errors, fmt.Errorf("unknown segment in '%s' at line %d: %s", filename, lineNr+1, param))
			}
		}

		// Otherwise, just append this line to the total output
		output = append(output, line)
	}
	return ":segment code\n" + strings.Join(output, "\n"), errors
}

func parseValue(value string) (int, error) {
	var i int64
	var err error
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "0x") {
		i, err = strconv.ParseInt(value[2:], 16, 64)
	} else if strings.HasPrefix(value, "0b") {
		i, err = strconv.ParseInt(value[2:], 2, 64)
	} else {
		i, err = strconv.ParseInt(value, 10, 64)
	}
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func find(line string, pattern string) (bool, string) {
	index := strings.Index(line, pattern)
	commentIndex := strings.Index(line, "#")
	if index == -1 || (commentIndex > -1 && index > commentIndex) {
		return false, ""
	}
	params := line[index+len(pattern):]
	if commentIndex != -1 {
		params = line[index+len(pattern) : commentIndex]
	}
	return true, strings.Trim(params, " \t")
}

func loadBinaryFile(filename string) (string, error) {
	if _, err := os.Stat(filename); err != nil {
		return "", fmt.Errorf("requested file '%s' not found", filename)
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("error reading file '%s': %s", filename, err.Error())
	}
	return dataToOctoText(file), nil
}

func dataToOctoText(data []byte) string {
	output := ""
	stride := 16
	for i := 0; i < len(data); i += stride {
		output += " "
		for j := i; j < len(data) && j < i+stride; j++ {
			output += fmt.Sprintf(" 0x%02x", data[j])
		}
		output += "\n"
	}
	return output
}
