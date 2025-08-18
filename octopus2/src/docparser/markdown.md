# {{ .Title }}

{{if .Description -}}{{ .Description }}
{{end}}
{{- $path := .Filepath -}}

## Table of Contents

{{range .Sections}}{{if .Name}}- {{ .Name }}
{{end}}{{if .Consts}}  - Constants
{{range .Consts}}    - [{{.Name}}](#{{.Name}})
{{end}}{{end}}{{if .Macros}}  - Macros
{{range .Macros}}    - [{{.Name}}](#{{.Name}})
{{end}}{{end}}{{if .Routines}}  - Routines
{{range .Routines}}    - [{{.Name}}](#{{.Name}}){{end}}{{end}}{{end}}

{{range .Sections}}{{if .Name}}# {{ .Name }}

{{end -}}{{if .Consts -}}
## Constants
{{end -}}

{{range .Consts}}
### `{{ .Name }}`

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
Value: `{{ .Value }}`
{{end}}

{{if .Macros -}}
## Macros
{{end -}}

{{range .Macros}}
### `{{ .Name }}`

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
{{- if .Parameters}}
#### Parameters
{{range .Parameters}}- {{ . }}
{{end -}}{{- end -}}
{{- end}}

{{if .Routines -}}
### Routines
{{end -}}

{{range .Routines}}
### `{{ .Name }}`

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
{{- end}}
{{- end}}
