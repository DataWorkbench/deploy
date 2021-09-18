{{/*
Service Addresses For ApiServer
*/}}
{{- define "service.jaeger" -}}
{{ include "databench.fullname" . }}-jaeger:{{ .Values.ports.jaeger }}
{{- end -}}

{{- define "service.udfmanager" -}}
{{ include "databench.fullname" . }}-udfmanager:{{ .Values.ports.udfmanager }}
{{- end -}}

{{- define "service.flowmanager" -}}
{{ include "databench.fullname" . }}-flowmanager:{{ .Values.ports.flowmanager }}
{{- end -}}

{{- define "service.sourcemanager" -}}
{{ include "databench.fullname" . }}-sourcemanager:{{ .Values.ports.sourcemanager }}
{{- end -}}

{{- define "service.jobdeveloper" -}}
{{ include "databench.fullname" . }}-jobdeveloper:{{ .Values.ports.jobdeveloper }}
{{- end -}}

{{- define "service.jobwatcher" -}}
{{ include "databench.fullname" . }}-jobwatcher:{{ .Values.ports.jobwatcher }}
{{- end -}}

{{- define "service.jobmanager" -}}
{{ include "databench.fullname" . }}-jobmanager:{{ .Values.ports.jobmanager }}
{{- end -}}

{{- define "service.zeppelinscale" -}}
{{ include "databench.fullname" . }}-zeppelinscale:{{ .Values.ports.zeppelinscale }}
{{- end -}}

{{- define "service.zeppelin" -}}
{{ include "databench.fullname" . }}-zeppelin:{{ .Values.ports.zeppelin }}
{{- end -}}

{{- define "service.hdfs" -}}
hdfs://{{- .Values.hdfs.service }}:{{ .Values.ports.hdfs }}
{{- end -}}

{{- define "apiserver.link.services" -}}
- name: API_SERVER_TRACER_LOCAL_AGENT
  value: '{{ include "service.jaeger" . }}'
- name: API_SERVER_SPACE_MANAGER_ADDRESS
  value: "{{ include "databench.fullname" . }}-spacemanager:{{ .Values.ports.spacemanager }}"
- name: API_SERVER_FLOW_MANAGER_ADDRESS
  value: '{{- include "service.flowmanager" . }}'
- name: API_SERVER_SCHEDULER_ADDRESS
  value: "{{ include "databench.fullname" . }}-scheduler:{{ .Values.ports.scheduler }}"
- name: API_SERVER_SOURCE_MANAGER_ADDRESS
  value: '{{ include "service.sourcemanager" . }}'
- name: API_SERVER_JOB_MANAGER_ADDRESS
  value: "{{ .Release.Name }}-jobmanager:{{ .Values.ports.jobmanager }}"
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

{{- define "mysql.url" -}}
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
{{ .Release.Name }}-client:{{- .Values.ports.etcd }}
{{- end -}}

{{- define "etcd.waiting.cmd" -}}
until nc -z {{ .Release.Name }}-etcd {{ .Values.ports.etcd }}; do echo "waiting for etcd.."; sleep 2; done;
{{- end -}}
