package docparser

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"
)

type Doc struct {
	Title       string
	Description *string
	Filename    string
	Filepath    string
	Consts      []Constant
	Macros      []Macro
	Routines    []Routine
}

type Constant struct {
	Line        int
	Name        string
	Value       string
	Description *string
}

type Macro struct {
	Line        int
	Name        string
	Parameters  []string
	Description *string
}

type Routine struct {
	Line        int
	Name        string
	Description *string
}

type DocBlock struct {
	Line      int
	Lines     int
	Content   string
	IsPrimary bool
}

func Parse(source string, filename string, tmpl string) (string, error) {
	data := findData(source, filename)
	tpl, err := template.New("docparser").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("could not parse template: %v", err)
	}
	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, data)
	if err != nil {
		return "", fmt.Errorf("could not apply template: %v", err)
	}
	return buf.String(), nil
}

func findData(source string, filename string) *Doc {
	lines := strings.Split(source, "\n")
	docBlocks := findDocBlocks(lines)
	title, description := findFileProperties(docBlocks)

	return &Doc{
		Title:       strOrDefault(title, path.Base(filename)),
		Description: description,
		Filename:    path.Base(filename),
		Filepath:    filename,
		Consts:      findConstants(lines, docBlocks),
		Macros:      findMacros(lines, docBlocks),
		Routines:    findRoutines(lines, docBlocks),
	}
}

// Find all blocks of comments in the source file, with their line numbers
func findDocBlocks(source []string) []DocBlock {
	result := make([]DocBlock, 0)
	inBlock := false
	curBlock := DocBlock{}
	for i, line := range source {
		line = strings.TrimSpace(line)
		isComment := strings.HasPrefix(line, "#")
		if isComment {
			if !inBlock {
				curBlock.Line = i + 1
				curBlock.IsPrimary = strings.HasPrefix(line, "###")
				inBlock = true
			}
			for strings.HasPrefix(line, "#") {
				line = strings.TrimPrefix(line, "#")
			}
			curBlock.Content += strings.TrimPrefix(line, " ") + "\n"
			curBlock.Lines += 1
		} else {
			if inBlock {
				inBlock = false
				result = append(result, curBlock)
				curBlock = DocBlock{}
			}
		}
	}
	return result
}

func findFileProperties(docBlocks []DocBlock) (*string, *string) {
	var title *string
	var description *string
	titleFound := false
	descriptionFound := false
	for _, block := range docBlocks {
		if block.IsPrimary {
			return title, description
		}
		lines := strings.Split(block.Content, "\n")
		if !titleFound && (len(lines) == 1 || len(lines) > 1 && strings.TrimSpace(lines[1]) == "") {
			value := lines[0]
			title = &value
			titleFound = true
			if len(lines) > 2 && strings.TrimSpace(lines[1]) == "" {
				value := strings.Join(lines[2:], "\n")
				description = &value
				return title, description
			}
			continue
		}
		if !descriptionFound {
			description = &block.Content
			return title, description
		}
	}
	return title, description
}

func findConstants(source []string, docBlocks []DocBlock) []Constant {
	result := make([]Constant, 0)
	for i, line := range source {
		line = strings.ToLower(strings.TrimSpace(line))
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		if parts[0] == ":const" {
			block := findRelatedDocBlock(source, i+1, docBlocks)
			if block != nil && block.IsPrimary {
				result = append(result, Constant{
					Line:        i + 1,
					Name:        parts[1],
					Value:       parts[2],
					Description: &block.Content,
				})
			}
		}
	}
	return result
}

func findMacros(source []string, docBlocks []DocBlock) []Macro {
	result := make([]Macro, 0)
	for i, line := range source {
		line = strings.ToLower(strings.TrimSpace(line))
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		if parts[0] == ":macro" {
			block := findRelatedDocBlock(source, i+1, docBlocks)
			if block != nil && block.IsPrimary {
				parameters := make([]string, 0)
				for _, part := range parts[2:] {
					if part != "{" {
						parameters = append(parameters, part)
					}
				}
				result = append(result, Macro{
					Line:        i + 1,
					Name:        parts[1],
					Description: &block.Content,
					Parameters:  parameters,
				})
			}
		}
	}
	return result
}

func findRoutines(source []string, docBlocks []DocBlock) []Routine {
	result := make([]Routine, 0)
	for i, line := range source {
		line = strings.ToLower(strings.TrimSpace(line))
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		if parts[0] == ":" {
			block := findRelatedDocBlock(source, i+1, docBlocks)
			if block != nil && block.IsPrimary {
				result = append(result, Routine{
					Line:        i + 1,
					Name:        parts[1],
					Description: &block.Content,
				})
			}
		}
	}
	return result
}

func findRelatedDocBlock(source []string, line int, docBlocks []DocBlock) *DocBlock {
	var prevBlock DocBlock
	for _, block := range docBlocks {
		if block.Line > line {
			// We passed the requested line, is the last one relevant?
			endPos := prevBlock.Line + prevBlock.Lines
			if endPos == line || (endPos+1 == line && strings.TrimSpace(source[endPos]) == "") {
				return &prevBlock
			}
			return nil
		}
		prevBlock = block
	}
	return nil
}

func strOrDefault(str *string, def string) string {
	if str == nil {
		return def
	}
	return *str
}
