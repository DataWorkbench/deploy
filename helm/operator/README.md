# Dataworkbench Deployment Operators
deploy dataworkbench-related services on K8S with operator

First, install operator using:
```
  $ helm install release-name operator/helm-operator
```

After, deploy dataworkbench-related services using:
```
  $ helm install release-name operator/helm
```
