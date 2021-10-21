{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "hdfs.fullname" -}}
{{- if contains $.Chart.Name $.Release.Name }}
{{- $.Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" $.Release.Name $.Chart.Name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}

{{- define "zookeeper.quorum" -}}
{{- if .Values.zookeeper.quorum }}
{{- .Values.zookeeper.quorum }}
{{- else }}
{{- printf "zk-0.zk-hs.%s.svc.cluster.local:2181,zk-1.zk-hs.%s.svc.cluster.local:2181,zk-2.zk-hs.%s.svc.cluster.local:2181" .Release.Namespace .Release.Namespace .Release.Namespace }}
{{- end }}
{{- end }}