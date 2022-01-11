#!/bin/bash


###############################################################################################
# home dir for hdfs / mysql / etcd volume
VolumeHome=/data/dataomnis


###############################################################################################
# Release name in volume
HdfsReleaseName=hdfs-cluster
MysqlReleaseName=mysql-cluster
EtcdReleaseName=etcd-cluster
RedisReleaseName=redis-cluster


###############################################################################################
AllNodes=(worker-1 worker-2 worker-3)

# hdfsDatadir format: ${VolumeHome}/${HdfsReleaseName}/datanode ..
# hdfs role-node map
HdfsDatanodeNodes=(worker-1 worker-2 worker-3)
HdfsNamenodeNodes=(worker-1 worker-2)
HdfsJournalnodeNodes=(worker-1 worker-2 worker-3)
HdfsZookeeperNodes=(worker-1 worker-2 worker-3)

# MysqlDatadir format: ${VolumeHome}/${MysqlReleaseName}
# mysql role-node map
MysqlNodes=(worker-1 worker-2 worker-3)

# EtcdDatadir format: ${VolumeHome}/${EtcdReleaseName}
# etcd role-node map
EtcdNodes=(worker-1 worker-2 worker-3)

# RedisDatadir format: ${VolumeHome}/${RedisReleaseName}
# redis role-node map
RedisNodes=(worker-1 worker-2 worker-3)


###############################################################################################
# HelmRepodir: /root/.cache/helm/repository
# create it on all node
for node in ${AllNodes[@]}
do
  ssh root@${node} "mkdir -p /root/.cache/helm/repository"
done

# create hdfs dir
for node in ${HdfsDatanodeNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${HdfsReleaseName}/datanode"
done
for node in ${HdfsNamenodeNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${HdfsReleaseName}/namenode"
done
for node in ${HdfsJournalnodeNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${HdfsReleaseName}/journalnode"
done
for node in ${HdfsZookeeperNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${HdfsReleaseName}/zookeeper"
done

# create mysql dir
for node in ${MysqlNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${MysqlReleaseName}"
done

# create etcd dir
for node in ${EtcdNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${EtcdReleaseName}"
done

# create etcd dir
for node in ${RedisNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${RedisReleaseName}"
done
