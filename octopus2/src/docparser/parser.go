package docparser

import (
	"path"
	"slices"
	"strings"
)

type Doc struct {
	Title       string
	Description *string
	Filename    string
	Filepath    string
	Sections    []Section
}

type Section struct {
	Name     string
	Consts   []Constant
	Macros   []Macro
	Routines []Routine
	start    int
	end      int
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

func findData(source string, filename string) *Doc {
	lines := strings.Split(source, "\n")
	docBlocks := findDocBlocks(lines)
	title, description := findFileProperties(docBlocks)
	sections := findSections(lines, docBlocks)

	return &Doc{
		Title:       strOrDefault(title, path.Base(filename)),
		Description: description,
		Filename:    path.Base(filename),
		Filepath:    filename,
		Sections:    sections,
	}
}

func findDocBlocks(source []string) []DocBlock {
	result := make([]DocBlock, 0)
	inBlock := false
	curBlock := DocBlock{}
	for i, line := range source {
		if strings.HasPrefix(line, "#") {
			if !inBlock {
				curBlock.Line = i + 1
				curBlock.IsPrimary = strings.HasPrefix(line, "###")
				inBlock = true
			}
			line = strings.TrimSpace(line)
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
	// We want the first comment in the file, that is not a "primary" comment
	// (`###`) and is in the first five lines.
	if len(docBlocks) == 0 {
		return nil, nil
	}
	block := docBlocks[0]
	if block.IsPrimary || block.Line > 5 {
		return nil, nil
	}

	// We have a valid comment, now see if we can parse it into a title and a
	// description
	if block.Lines == 1 {
		return &block.Content, nil // Single line comment is considered a title
	}
	lines := strings.Split(block.Content, "\n")
	if len(lines) > 2 && strings.TrimSpace(lines[1]) == "" {
		title := lines[0]
		description := strings.Join(lines[2:], "\n")
		if strings.TrimSpace(description) == "" {
			return &title, nil
		} else {
			return &title, toDescription(description) // Block split by an empty line on line 2 is considered a title and a description
		}
	}
	return nil, toDescription(block.Content) // Whole block is considered a description
}

func findSections(source []string, docBlocks []DocBlock) []Section {
	sections := []Section{{
		Name:  "",
		start: 0,
	}}
	for _, block := range docBlocks {
		if !block.IsPrimary {
			continue
		}
		if !(block.Lines == 1) {
			continue
		}
		name := block.Content
		suffixes := 0
		name = strings.TrimSpace(name)
		for strings.HasSuffix(name, "#") {
			name = strings.TrimSuffix(name, "#")
			suffixes += 1
		}
		name = strings.TrimSpace(name)
		if len(name) >= 1 && suffixes >= 3 {
			sections[len(sections)-1].end = block.Line
			sections = append(sections, Section{
				start: block.Line,
				Name:  name,
			})
		}
	}
	sections[len(sections)-1].end = len(source)
	for i, s := range sections {
		sections[i].Consts = findConstants(source, s.start, s.end, docBlocks)
		sections[i].Macros = findMacros(source, s.start, s.end, docBlocks)
		sections[i].Routines = findRoutines(source, s.start, s.end, docBlocks)
	}
	return slices.DeleteFunc(sections, func(s Section) bool {
		return len(s.Consts) == 0 && len(s.Macros) == 0 && len(s.Routines) == 0
	})
}

func findConstants(source []string, start, end int, docBlocks []DocBlock) []Constant {
	result := make([]Constant, 0)
	for i := start; i < end; i++ {
		line := source[i]
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
					Description: toDescription(block.Content),
				})
			}
		}
	}
	return result
}

func findMacros(source []string, start, end int, docBlocks []DocBlock) []Macro {
	result := make([]Macro, 0)
	for i := start; i < end; i++ {
		line := source[i]
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
					Description: toDescription(block.Content),
					Parameters:  parameters,
				})
			}
		}
	}
	return result
}

func findRoutines(source []string, start, end int, docBlocks []DocBlock) []Routine {
	result := make([]Routine, 0)
	for i := start; i < end; i++ {
		line := source[i]
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
					Description: toDescription(block.Content),
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

func toDescription(descr string) *string {
	for strings.HasPrefix(descr, "\n") {
		descr = strings.TrimPrefix(descr, "\n")
	}
	for strings.HasSuffix(descr, "\n") {
		descr = strings.TrimSuffix(descr, "\n")
	}

	// Parse code blocks in the description
	result := ""
	inCodeBlock := false
	inIndentedBlock := false
	for _, line := range strings.Split(descr, "\n") {
		if strings.TrimSpace(line) == "```" {
			inCodeBlock = !inCodeBlock
		}
		if !inCodeBlock {
			if inIndentedBlock {
				if !strings.HasPrefix(line, " ") {
					result += "```\n"
					inIndentedBlock = false
				}
			} else {
				if strings.HasPrefix(line, " ") {
					result += "```\n"
					inIndentedBlock = true
				}
			}
		}
		result += line + "\n"
	}
	if inIndentedBlock {
		result += "```\n"
	}

	return &result
}
