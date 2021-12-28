FROM dataworkbench/builder:zeppelin as builder
FROM apache/zeppelin:0.9.0

ENV HADOOP_HOME /zeppelin/hadoop/hadoop-2.7.5
ENV PATH $PATH:$HADOOP_HOME/bin
ENV HADOOP_CLASSPATH /zeppelin/hadoop/hadoop-2.7.5/etc/hadoop:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/common/lib/*:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/common/*:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/hdfs:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/hdfs/lib/*:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/hdfs/*:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/yarn/lib/*:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/yarn/*:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/mapreduce/lib/*:/zeppelin/hadoop/hadoop-2.7.5/share/hadoop/mapreduce/*
COPY --from=builder /opt/flink-1.12.3 /zeppelin/flink/flink-1.12.3
COPY --from=builder /opt/hadoop-2.7.5 ${HADOOP_HOME}

COPY conf/* /opt/zeppelin/conf/
COPY zeppelin.sh /opt/zeppelin/bin/
COPY 100-interpreter-spec.yaml /opt/zeppelin/k8s/interpreter/
COPY 1.12_lib/* /zeppelin/flink/1.12_lib/
COPY depends/* /zeppelin/flink/depends/

USER root
RUN chmod 777 -R /zeppelin/flink/ && \
    chmod 777 -R /opt/zeppelin/bin
CMD ["bin/zeppelin.sh"]
