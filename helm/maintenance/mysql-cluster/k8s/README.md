全量备份：
    先执行kubectl get pxc-backup -n dataomnis 命令查询pxc集群的名称
    如：
    NAME                                           CLUSTER           STORAGE   DESTINATION                                           STATUS      COMPLETED   AGE
    cron-dataomnis-mysql-fs-pvc-2021126000-372f8   dataomnis-mysql   fs-pvc    pvc/xb-cron-dataomnis-mysql-fs-pvc-2021126000-372f8   Succeeded   24h         32h
    cron-dataomnis-mysql-fs-pvc-2021127000-372f8   dataomnis-mysql                                                                                           8h
    那么集群的名称就是 dataomnis-mysql，然后修改backup.yaml里 pxcCluster 的值，再执行kubectl apply -f backup.yaml，就可以全量备份数据了
    
全量恢复：
   修改restore.yaml里 pxcCluster 跟 backupName， pxcCluster 是pxc集群的名称，backupName是备份的名称
   再执行kubectl apply -f restore.yaml，就可以全量恢复数据了，恢复数据期间会删除现在pod，等数据恢复好会重新启动3个pxc
   
   