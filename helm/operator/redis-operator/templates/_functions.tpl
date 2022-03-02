{{/* Build the Spotahome standard labels */}}
{{- define "common-labels" -}}
app.kubernetes.io/name: {{ .Chart.Name | quote }}
{{- end }}

{{- define "helm-labels" -}}
{{ include "common-labels" . }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{- end }}

{{/* Build wide-used variables the application */}}
{{ define "name" -}}
{{- if contains .Chart.Name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{ define "operator.image" -}}
{{ printf "%s/%s:%s" .Values.image.registry .Values.image.operator }}
{{- end }}
