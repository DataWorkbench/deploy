{{/*
Service Addresses For ApiServer
*/}}
{{- define "service.udfmanager" -}}
{{ .Release.Name }}-udfmanager:{{ .Values.ports.udfmanager }}
{{- end -}}
{{- define "service.flowmanager" -}}
{{ .Release.Name }}-flowmanager:{{ .Values.ports.flowmanager }}
{{- end -}}
{{- define "service.sourcemanager" -}}
{{ .Release.Name }}-sourcemanager:{{ .Values.ports.sourcemanager }}
{{- end -}}
{{- define "service.jobdeveloper" -}}
{{ .Release.Name }}-jobdeveloper:{{ .Values.ports.jobdeveloper }}
{{- end -}}
{{- define "service.jobwatcher" -}}
{{ .Release.Name }}-jobwatcher:{{ .Values.ports.jobwatcher }}
{{- end -}}
{{- define "service.jobmanager" -}}
{{ .Release.Name }}-jobmanager:{{ .Values.ports.jobmanager }}
{{- end -}}
{{- define "service.zeppelinscale" -}}
{{ .Release.Name }}-zeppelinscale:{{ .Values.ports.zeppelinscale }}
{{- end -}}

{{/*
- name: API_SERVER_SCHEDULER_ADDRESS
  value: "{{ .Release.Name }}-scheduler:{{ .Values.ports.scheduler }}"
- name: API_SERVER_SOURCE_MANAGER_ADDRESS
  value: {{ include "service.sourcemanager" . | quote }}
*/}}
{{- define "apiserver.link.services" -}}
- name: API_SERVER_TRACER_LOCAL_AGENT
  value: "{{ .Release.Name }}-jaeger:{{ .Values.ports.jaeger }}"
- name: API_SERVER_SPACE_MANAGER_ADDRESS
  value: "{{ .Release.Name }}-spacemanager:{{ .Values.ports.spacemanager }}"
- name: API_SERVER_FLOW_MANAGER_ADDRESS
  value: "{{ .Release.Name }}-flowmanager:{{ .Values.ports.flowmanager }}"
{{- end -}}


{{/*
Mysql Settings
*/}}
{{- define "mysql.host" -}}
{{ .Release.Name }}-mysql
{{- end -}}

{{- define "mysql.port" -}}
{{ .Values.ports.mysql }}
{{- end -}}

{{- define "mysql.hostPort" -}}
{{ .Release.Name }}-mysql:{{- .Values.ports.mysql }}
{{- end -}}

{{- define "mysql.root.password" -}}
{{- .Values.mysql.password }}
{{- end -}}

{{- define "mysql.waiting.cmd" -}}
until nc -z {{ .Release.Name }}-mysql {{ .Values.ports.mysql }}; do echo "waiting for mysql.."; sleep 2; done;
{{- end -}}

{{/*
Etcd Settings
*/}}
{{- define "etcd.endpoints" -}}
{{ .Release.Name }}-etcd:{{- .Values.ports.mysql }}
{{- end -}}
{{- define "etcd.waiting.cmd" -}}
until nc -z {{ .Release.Name }}-etcd {{ .Values.ports.etcd }}; do echo "waiting for etcd.."; sleep 2; done;
{{- end -}}
