{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "redis-cluster.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "redis-cluster.fullname" -}}
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

{{- define "sentinel.image" -}}
{{- if .Values.image.registry -}}
{{- .Values.image.registry }}/library/{{- .Values.image.sentinel }}:{{- .Chart.AppVersion }}
{{- else -}}
{{- .Values.image.sentinel }}:{{- .Chart.AppVersion }}
{{- end -}}
{{- end -}}

{{- define "sentinel.replica" -}}
{{- if .Values.redis.persistent.localPv.nodes }}
{{ len .Values.redis.persistent.localPv.nodes }}
{{- else }}
{{ .Values.sentinel.replicaCount }}
{{- end }}
{{- end }}

{{- define "redis.image" -}}
{{- if .Values.image.registry -}}
{{- .Values.image.registry }}/library/{{- .Values.image.redis }}:{{- .Chart.AppVersion }}
{{- else -}}
{{- .Values.image.redis }}:{{- .Chart.AppVersion }}
{{- end -}}
{{- end -}}

{{- define "redis.replica" -}}
{{- if .Values.redis.persistent.localPv.nodes }}
{{ len .Values.redis.persistent.localPv.nodes }}
{{- else }}
{{ .Values.redis.replicaCount }}
{{- end }}
{{- end }}
