# Dataomnis Deployment Operators
deploy dataomnis-related services on K8S with operator

First, install operator using:
```
  $ helm install release-name operator/helm-operator
```

After, deploy dataomnis-related services using:
```
  $ helm install release-name operator/helm
```
