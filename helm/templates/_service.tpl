{{/*
Service Addresses For ApiServer
*/}}
{{- define "service.udfmanager" -}}
{{ .Release.Name }}-udfmanager:{{ .Values.ports.udfmanager }}
{{- end -}}
{{- define "service.sourcemanager" -}}
{{ .Release.Name }}-sourcemanager:{{ .Values.ports.sourcemanager }}
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

{{- define "mysql.root.password" -}}
{{- .Values.mysql.password }}
{{- end -}}

{{- define "mysql.waiting.cmd" -}}
until nc -z {{ .Release.Name }}-mysql {{ .Values.ports.mysql }}; do echo "waiting for mysql.."; sleep 2; done;
{{- end -}}

