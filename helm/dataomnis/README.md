# Dataomnis deployment 
deploy dataomnis services on K8S with helm

if add new dataomnis service, create file named `SERVICE-deployment.yaml` under templates folder;

if add dependency service by dataomnis service, eg: mysql / etcd service, create helm folder named the service name under charts;