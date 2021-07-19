# Dataworkbench deployment 
deploy dataworkbench services on K8S with helm

if add new dataworkbench service, create file named `SERVICE-deployment.yaml` under templates folder;

if add dependency service by dataworkbench service, eg: mysql / etcd service, create helm folder named the service name under charts;