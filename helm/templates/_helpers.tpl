{{/* vim: set filetype=mustache: */}}

{{- define "databench.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{- define "databench.image" -}}
{{- .Values.image.registry }}/{{ .Values.image.databench }}:{{ .Values.image.tag }}
{{- end -}}
{{- define "flyway.image" -}}
{{- .Values.image.registry }}/{{ .Values.image.flyway }}:{{ .Values.image.tag }}
{{- end -}}
{{- define "zeppelin.image" -}}
{{- .Values.image.registry }}/{{ .Values.image.zeppelin }}:{{ .Values.image.zeppelinTag }}
{{- end -}}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "databench.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "databench.labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
helm.sh/chart: {{ include "databench.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}
