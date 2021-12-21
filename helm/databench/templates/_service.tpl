{{/*
Service Addresses For ApiServer
*/}}
{{- define "service.jaeger" -}}
{{ include "databench.fullname" . }}-jaeger:{{ .Values.ports.jaeger }}
{{- end -}}

{{- define "service.spacemanager" -}}
{{ include "databench.fullname" . }}-spacemanager:{{ .Values.ports.spaceManager }}
{{- end -}}

{{- define "service.udfmanager" -}}
{{ include "databench.fullname" . }}-udfmanager:{{ .Values.ports.udfManager }}
{{- end -}}

{{- define "service.scheduler" -}}
{{ include "databench.fullname" . }}-scheduler:{{ .Values.ports.scheduler }}
{{- end -}}

{{- define "service.flowmanager" -}}
{{ include "databench.fullname" . }}-flowmanager:{{ .Values.ports.flowManager }}
{{- end -}}

{{- define "service.sourcemanager" -}}
{{ include "databench.fullname" . }}-sourcemanager:{{ .Values.ports.sourceManager }}
{{- end -}}

{{- define "service.jobdeveloper" -}}
{{ include "databench.fullname" . }}-jobdeveloper:{{ .Values.ports.jobDeveloper }}
{{- end -}}

{{- define "service.jobwatcher" -}}
{{ include "databench.fullname" . }}-jobwatcher:{{ .Values.ports.jobWatcher }}
{{- end -}}

{{- define "service.jobmanager" -}}
{{ include "databench.fullname" . }}-jobmanager:{{ .Values.ports.jobManager }}
{{- end -}}

{{- define "service.zeppelin" -}}
{{ include "databench.fullname" . }}-zeppelin-server:{{ .Values.ports.zeppelin }}
{{- end -}}

{{- define "service.resourcemanager" -}}
{{ include "databench.fullname" . }}-resourcemanager:{{ .Values.ports.resourceManager }}
{{- end -}}

{{- define "service.enginemanager" -}}
{{ include "databench.fullname" . }}-enginemanager:{{ .Values.ports.engineManager }}
{{- end -}}

{{- define "service.account" -}}
{{ include "databench.fullname" . }}-account:{{ .Values.ports.account }}
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
- name: JOB_MANAGER_JOBDEVELOPER_SERVER_ADDRESS
  value: '{{ include "service.jobdeveloper" . }}'
- name: JOB_MANAGER_JOBWATCHER_SERVER_ADDRESS
  value: '{{ include "service.jobwatcher" . }}'
- name: JOB_MANAGER_ENGINEMANAGER_SERVER_ADDRESS
  value: '{{ include "service.enginemanager" . }}'
- name: JOB_MANAGER_ZEPPELIN_ADDRESS
  value: '{{ include "service.zeppelin" . }}'
{{- end -}}


{{- define "scheduler.link.services" -}}
- name: SCHEDULER_JOB_MANAGER_ADDRESS
  value: '{{ include "service.jobmanager" . }}'
- name: SCHEDULER_FLOW_MANAGER_ADDRESS
  value: '{{ include "service.flowmanager" . }}'
{{- end -}}

{{- define "spacemanager.link.services" -}}
- name: SCHEDULER_JOB_MANAGER_ADDRESS
  value: '{{ include "service.jobmanager" . }}'
- name: SCHEDULER_FLOW_MANAGER_ADDRESS
  value: '{{ include "service.flowmanager" . }}'
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

{{- define "service.hdfs" -}}
hdfs://{{ .Release.Name }}-hdfs-http:{{ .Values.ports.hdfs }}
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
