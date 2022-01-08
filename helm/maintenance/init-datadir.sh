#!/bin/bash


###############################################################################################
# home dir for hdfs / mysql / etcd volume
VolumeHome=/data/databench


###############################################################################################
# Release name in volume
HdfsReleaseName=hdfs-cluster
MysqlReleaseName=mysql-cluster
EtcdReleaseName=etcd-cluster


###############################################################################################
# hdfsDatadir format: ${VolumeHome}/${HdfsReleaseName}/datanode ..
# hdfs role-node map
HdfsDatanodeNodes=(worker-s001 worker-s002 worker-s003)
HdfsNamenodeNodes=(worker-s001 worker-s002)
HdfsJournalnodeNodes=(worker-s001 worker-s002 worker-s003)
HdfsZookeeperNodes=(worker-s001 worker-s002 worker-s003)

# MysqlDatadir format: ${VolumeHome}/${MysqlReleaseName}
# mysql role-node map
MysqlNodes=(worker-s001 worker-s002 worker-s003)

# EtcdDatadir format: ${VolumeHome}/${EtcdReleaseName}
# etcd role-node map
EtcdNodes=(worker-s001 worker-s002 worker-s003)


###############################################################################################
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
