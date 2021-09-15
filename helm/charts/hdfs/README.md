# Dataworkbench deployment 
deploy HDFS services on K8S with helm


## View k8s nodes 
 kubectl get node -o wide
NAME          STATUS   ROLES    AGE   VERSION   INTERNAL-IP   
master1       Ready    master   74d   v1.19.8   192.168.0.5  
worker-s001   Ready    worker   74d   v1.19.8   192.168.0.4  
worker-s002   Ready    worker   74d   v1.19.8   192.168.0.6   
worker-s003   Ready    worker   71d   v1.19.8   192.168.0.7   


## persistent for namenode journalnode zookeeper
### local pv 

Select worker-s001 and worker-s002 as the namenodes and create local pv path
```
# execute on worker-s001 and worker-s002

$ mkdir -p /mnt/hdfs/nn  
```

Select worker-s001 ,worker-s002 and worker-s003 as the journalnodes and zookeeper pod local pv path
```
# execute on worker-s001 ,worker-s002 and worker-s003

$ mkdir -p /mnt/hdfs/jn
$ mkdir -p /mnt/hdfs/zk
```

### storageClass

config __storageClass__ in __values.yaml__ .


## Label the datanode nodes that need to be deployed(required)
```
$ kubectl label node worker-s003 hdfs/nodetype=datanode
$ kubectl label node worker-s003 hdfs/nodetype=datanode
$ kubectl label node worker-s003 hdfs/nodetype=datanode
```


## deploy hdfs cluster
```
$ helm install hdfs . --set zk01Node=worker-s001,zk02Node=worker-s002,zk03Node=worker-s003,jn01Node=worker-s001,jn02Node=worker-s002,jn03Node=worker-s003,nn01Node=worker-s001,nn02Node=worker-s002
```
