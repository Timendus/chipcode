package docparser

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed markdown.md
var markdown string

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

func ParseToMarkdown(source string, filename string) (string, error) {
	return Parse(source, filename, markdown)
}
