{{/*
Service Addresses For ApiServer
*/}}
{{- define "service.jaeger" -}}
{{ include "dataomnis.fullname" . }}-jaeger:{{ .Values.ports.jaeger }}
{{- end -}}

{{- define "service.spacemanager" -}}
{{ include "dataomnis.fullname" . }}-spacemanager:{{ .Values.ports.spacemanager }}
{{- end -}}

{{- define "service.udfmanager" -}}
{{ include "dataomnis.fullname" . }}-udfmanager:{{ .Values.ports.udfmanager }}
{{- end -}}

{{- define "service.scheduler" -}}
{{ include "dataomnis.fullname" . }}-scheduler:{{ .Values.ports.scheduler }}
{{- end -}}

{{- define "service.flowmanager" -}}
{{ include "dataomnis.fullname" . }}-flowmanager:{{ .Values.ports.flowmanager }}
{{- end -}}

{{- define "service.sourcemanager" -}}
{{ include "dataomnis.fullname" . }}-sourcemanager:{{ .Values.ports.sourcemanager }}
{{- end -}}

{{- define "service.jobdeveloper" -}}
{{ include "dataomnis.fullname" . }}-jobdeveloper:{{ .Values.ports.jobdeveloper }}
{{- end -}}

{{- define "service.jobwatcher" -}}
{{ include "dataomnis.fullname" . }}-jobwatcher:{{ .Values.ports.jobwatcher }}
{{- end -}}

{{- define "service.jobmanager" -}}
{{ include "dataomnis.fullname" . }}-jobmanager:{{ .Values.ports.jobmanager }}
{{- end -}}

{{- define "service.zeppelin" -}}
{{ include "dataomnis.fullname" . }}-zeppelin-server:{{ .Values.ports.zeppelin }}
{{- end -}}

{{- define "service.resourcemanager" -}}
{{ include "dataomnis.fullname" . }}-resourcemanager:{{ .Values.ports.resourcemanager }}
{{- end -}}

{{- define "service.enginemanager" -}}
{{ include "dataomnis.fullname" . }}-enginemanager:{{ .Values.ports.enginemanager }}
{{- end -}}

{{- define "service.account" -}}
{{ include "dataomnis.fullname" . }}-account:{{ .Values.ports.account }}
{{- end -}}


{{- define "apiserver.link.services" -}}
- name: API_SERVER_TRACER_LOCAL_AGENT
  value: '{{- include "service.jaeger" . }}'
- name: API_SERVER_SPACE_MANAGER_ADDRESS
  value: '{{- include "service.spacemanager" . }}'
- name: API_SERVER_FLOW_MANAGER_ADDRESS
  value: '{{- include "service.flowmanager" . }}'
- name: API_SERVER_SCHEDULER_ADDRESS
  value: '{{- include "service.scheduler" . }}'
- name: API_SERVER_SOURCE_MANAGER_ADDRESS
  value: '{{ include "service.sourcemanager" . }}'
- name: API_SERVER_JOB_MANAGER_ADDRESS
  value: '{{ include "service.jobmanager" . }}'
- name: API_SERVER_UDF_MANAGER_ADDRESS
  value: '{{ include "service.udfmanager" . }}'
- name: API_SERVER_RESOURCE_MANAGER_ADDRESS
  value: '{{ include "service.resourcemanager" . }}'
- name: API_SERVER_ACCOUNT_SERVER_ADDRESS
  value: '{{ include "service.account" . }}'
- name: API_SERVER_ENGINE_MANAGER_ADDRESS
  value: '{{ include "service.enginemanager" . }}'
{{- end -}}


{{- define "jobdeveloper.link.services" -}}
- name: JOB_DEVELOPER_SOURCEMANAGER_SERVER_ADDRESS
  value: '{{ include "service.sourcemanager" . }}'
- name: JOB_DEVELOPER_UDFMANAGER_SERVER_ADDRESS
  value: '{{ include "service.udfmanager" . }}'
- name: JOB_DEVELOPER_RESOURCEMANAGER_SERVER_ADDRESS
  value: '{{ include "service.resourcemanager" . }}'
- name: JOB_DEVELOPER_ENGINEMANAGER_SERVER_ADDRESS
  value: '{{ include "service.enginemanager" . }}'
{{- end -}}


{{- define "jobmanager.link.services" -}}
- name: JOB_MANAGER_ZEPPELIN_ADDRESS
  value: '{{ include "service.zeppelin" . }}'
- name: JOB_MANAGER_RESOURCEMANAGER_SERVER_ADDRESS
  value: '{{ include "service.resourcemanager" . }}'
- name: JOB_MANAGER_UDFMANAGER_SERVER_ADDRESS
  value: '{{ include "service.udfmanager" . }}'
- name: JOB_MANAGER_ENGINEMANAGER_SERVER_ADDRESS
  value: '{{ include "service.enginemanager" . }}'
- name: JOB_MANAGER_SPACE_MANAGER_ADDRESS
  value: '{{- include "service.spacemanager" . }}'
{{- end -}}


{{- define "scheduler.link.services" -}}
- name: SCHEDULER_JOB_MANAGER_ADDRESS
  value: '{{ include "service.jobmanager" . }}'
- name: SCHEDULER_FLOW_MANAGER_ADDRESS
  value: '{{ include "service.flowmanager" . }}'
{{- end -}}

{{- define "spacemanager.link.services" -}}
- name: SPACE_MANAGER_JOB_MANAGER_ADDRESS
  value: '{{ include "service.jobmanager" . }}'
- name: SPACE_MANAGER_ENGINE_MANAGER_ADDRESS
  value: '{{ include "service.enginemanager" . }}'
- name: SPACE_MANAGER_SCHEDULER_ADDRESS
  value: '{{- include "service.scheduler" . }}'
{{- end -}}


{{/*
Mysql Settings
*/}}
{{- define "mysql.host" -}}
{{- if .Values.mysql.internal -}}
{{ .Release.Name }}-mysql
{{- else -}}
{{ .Values.mysql.externalHost }}
{{- end -}}
{{- end -}}

{{- define "mysql.hostPort" -}}
{{- if .Values.mysql.internal -}}
{{ .Release.Name }}-mysql:{{ .Values.ports.mysql }}
{{- else -}}
{{ .Values.mysql.externalHost }}:{{ .Values.ports.mysql }}
{{- end -}}
{{- end -}}

{{- define "mysql.waiting.cmd" -}}
until nc -z {{ include "mysql.host" . }} {{ .Values.ports.mysql }}; do echo "waiting for mysql.."; sleep 2; done;
{{- end -}}

{{/*
Etcd Settings
*/}}
{{- define "etcd.endpoints" -}}
{{ .Values.etcd.endpoint }}:{{- .Values.ports.etcd }}
{{- end -}}

{{- define "etcd.waiting.cmd" -}}
until nc -z {{ .Values.etcd.endpoint }} {{ .Values.ports.etcd }}; do echo "waiting for etcd.."; sleep 2; done;
{{- end -}}

{{- define "service.redis" -}}
{{ .Values.redis.address }}:{{ .Values.ports.redis }}
{{- end -}}
