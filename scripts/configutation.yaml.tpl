
# the home path of local pv to storage data, recommend：the path that mount a disk;
localPvHome: /data

# if config pullSecrets, create docker registry secret by kubectl:
# kubectl create secret docker-registry my-docker-registry-secret
#                                       --docker-server=<your-registry-server>
#                                       --docker-username=<your-name>
#                                       --docker-password=<your-pword>
#                                       --docker-email=<your-email>
# then:
# pullSecrets:
#   - my-docker-registry-secret
image:
  registry:
  pullSecrets:
  pullPolicy:

dataomnis:
  version: dev
  domain: dataomnis.192.168.27.90.nip.io
  port:

  mysql:
    maxIdleConn: 16
    maxOpenConn: 128
    connMaxLifetime: 10m
    logLevel: 4  # 1 => Silent, 2 => Error, 3 => Warn, 4 => Info
    slowThreshold: 2s

  global:
    replicas: 1
    strategy: RollingUpdate
    logLevel: 1  # 1=>"debug", 2=>"info", 3=>"warn", 4=>"error", 5=>"fatal"
    grpcLog:
      level: 2  #  1 => info, 2 => waring, 3 => error, 4 => fatal
      verbosity: 99
    metrics:
      enabled: true
      urlPath: "/metrics"

  webservice:
    enabled: false

  apiglobal:
    enabled: false
    regions:
      - hosts: http://api.dataomnis.192.168.27.90.nip.io
        enUsName: testing
        zhCnName: 测试区
    envs:

  iaas:

mysqlCluster:
  pxc:
    resources:
      requests:
        cpu:
        memory:
      limits:
        cpu:
        memory:

    persistent:
      size: 50Gi
      localPv:
        nodes:
          - worker-s001
          - worker-s002
          - worker-s003

hdfsCluster:
  nodes:
    - worker-s001
    - worker-s002
    - worker-s003

  namenode:
    persistent:
      size: 10Gi
      localPv:
        nodes:

  datanode:
    persistent:
      size: 100Gi
      localPv:
        nodes:

  journalnode:
    persistent:
      size: 20Gi
      localPv:
        nodes:

  zookeeper:
    persistent:
      size: 20Gi
      localPv:
        nodes:

redisCluster:
  redis:
    resources:
      requests:
        cpu:
        memory:
      limits:
        cpu:
        memory:

    persistent:
      size: 50Gi
      localPv:
        nodes:
          - worker-s001
          - worker-s002
          - worker-s003

etcdCluster:
  persistent:
    size: 50Gi
    localPv:
      nodes:
        - worker-s001
        - worker-s002
        - worker-s003