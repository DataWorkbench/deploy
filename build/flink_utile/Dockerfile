FROM flink:1.12.3-scala_2.11

#RUN cp -r /opt/flink/opt/* /opt/flink/lib/
RUN mkdir -p /opt/flink/plugins/s3-fs-presto/
RUN mkdir -p /opt/flink/plugins/s3-fs-hadoop/
RUN cp /opt/flink/opt/flink-s3-fs-presto-1.12.3.jar /opt/flink/plugins/s3-fs-presto/
RUN cp /opt/flink/opt/flink-s3-fs-hadoop-1.12.3.jar /opt/flink/plugins/s3-fs-hadoop/
COPY flink-shaded-hadoop-2-uber-2.8.3-10.0.jar /opt/flink/lib/
#RUN echo 's3.access-key: aa' >>  /opt/flink/conf/flink-conf.yaml
#RUN echo 's3.secret-key: bb' >>  /opt/flink/conf/flink-conf.yaml
#RUN echo 's3.endpoint: http://s3.gd2.qingstor.com' >>  /opt/flink/conf/flink-conf.yaml


RUN sed -i 's#deb.debian.org#mirrors.aliyun.com#g' /etc/apt/sources.list
RUN apt-get update && \
    apt-get install -y python3 python3-pip && \
    ln /usr/bin/python3 /usr/bin/python && \
    apt-get clean
RUN pip3 install --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple/ apache-flink==1.12.3

USER root:root
