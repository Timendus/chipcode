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
{{range .Routines}}    - [{{.Name}}](#{{replaceAll .Name "." "" | toLower }})
{{end}}{{end}}{{end}}

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
### `{{ .Name }}`{{range .Parameters}} `{{ . }}`{{end}}

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
{{- if .Inputs}}
#### Inputs
{{range .Inputs}}- {{index . 0}}{{ if index . 1}} - {{index . 1}}{{end}}
{{end -}}{{- end}}
{{- if .Outputs}}
#### Outputs
{{range .Outputs}}- {{index . 0}}{{ if index . 1}} - {{index . 1}}{{end}}
{{end -}}{{- end}}
{{- if .Destroys}}
#### Destroys
{{range .Destroys}}- {{ . }}
{{end -}}{{- end}}
{{- if .DependsOn}}
#### Depends on
{{range .DependsOn}}- {{ . }}
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
