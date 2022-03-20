# Dataomnis Deployment Operators

redis cluster requires PV twice the number of {{ $.Values.masterSize }} by default.
For example, please create {{ $.Values.storageSpec.localPv.home }}/{{ $.Release.Name }}/01,{{ $.Values.storageSpec.localPv.home }}/{{ $.Release.Name }}/02 in advance
