# {{ .Title }}

{{if .Description -}}{{ .Description }}
{{end}}
{{- $path := .Filepath -}}

{{if .Consts -}}
## Constants
{{end -}}

{{range .Consts}}
### {{ .Name }}

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
### {{ .Name }}

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
## Routines
{{end -}}

{{range .Routines}}
### {{ .Name }}

_{{ $path }}:{{ .Line }}_
{{if .Description}}
{{ .Description }}
{{- end}}
{{- end}}
