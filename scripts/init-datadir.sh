#!/bin/bash


###############################################################################################
# home dir for hdfs / mysql / etcd volume
VolumeHome=/data/dataomnis

###############################################################################################
AllNodes=(worker-s001 worker-s002 worker-s003)

# hdfsDatadir format: ${VolumeHome}/${HdfsReleaseName}/datanode ..
# hdfs role-node map
HdfsNodes=(worker-s001 worker-s002 worker-s003)

# MysqlDatadir format: ${VolumeHome}/${MysqlReleaseName}
# mysql role-node map
MysqlNodes=(worker-s001 worker-s002 worker-s003)

# EtcdDatadir format: ${VolumeHome}/${EtcdReleaseName}
# etcd role-node map
EtcdNodes=(worker-s001 worker-s002 worker-s003)

# RedisDatadir format: ${VolumeHome}/${RedisReleaseName}
# redis role-node map
RedisNodes=(worker-s001 worker-s002 worker-s003)


###############################################################################################
# Release name in volume
HdfsReleaseName=hdfs-cluster
MysqlReleaseName=mysql-cluster
EtcdReleaseName=etcd-cluster
RedisReleaseName=redis-cluster


###############################################################################################
# HelmRepodir: /root/.cache/helm/repository
# create it on all node
for node in ${AllNodes[@]}
do
  ssh root@${node} "mkdir -p /root/.cache/helm/repository"
done

# create hdfs dir
for node in ${HdfsNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${HdfsReleaseName}/{datanode,namenode,journalnode,zookeeper}"
done

# create mysql dir
for node in ${MysqlNodes[@]}
do
  ssh root@${node} "mkdir -p ${VolumeHome}/${MysqlReleaseName}/{data,log,mysql-bin}"
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
