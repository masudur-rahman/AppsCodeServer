{{- define "labels" -}}
chart-name: {{ .Chart.Name | quote }}
release-name: {{ .Release.Name | quote }}
{{- range $key, $val := .Values.labels }}
{{ $key }}: {{ $val | quote }}
{{- end }}
{{- end }}
