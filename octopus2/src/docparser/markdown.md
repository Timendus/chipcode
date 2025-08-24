# {{ .Title }}

{{if .Description -}}{{ .Description }}
{{end}}
{{- $path := .Filepath -}}

## Table of Contents

{{range .Sections}}{{if .Name}}- {{ .Name }}
{{end}}{{if .Consts}}  - Constants
{{range .Consts}}    - [{{.Name}}](#{{replaceAll .Name "." "" | toLower }})
{{end}}{{end}}{{if .Macros}}  - Macros
{{range .Macros}}    - [{{.Name}}](#{{replaceAll .Name "." "" | toLower }})
{{end}}{{end}}{{if .Routines}}  - Routines
{{range .Routines}}    - [{{.Name}}](#{{replaceAll .Name "." "" | toLower }}){{end}}{{end}}{{end}}

{{range .Sections -}}

{{- if .Name}}# {{ .Name }}

{{end -}}

{{- if .Consts -}}
## Constants

{{range .Consts -}}
### `{{ .Name }}`

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
Value: `{{ .Value }}`
{{end}}
{{end -}}

{{- if .Macros -}}
## Macros

{{range .Macros -}}
### `{{ .Name }}`

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
{{- if .Parameters}}
#### Parameters
{{range .Parameters}}- {{ . }}
{{end -}}{{- end}}
{{end}}{{end -}}

{{if .Routines -}}
### Routines

{{range .Routines -}}
### `{{ .Name }}`

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
{{- end}}
{{- end}}
{{end -}}
