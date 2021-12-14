FROM dataworkbench/builder:zeppelin as builder
FROM apache/zeppelin:0.9.0

ENV HADOOP_HOME /zeppelin/hadoop/hadoop-2.7.5
ENV PATH $PATH:$HADOOP_HOME/bin

COPY --from=builder /opt/flink-1.12.3 /zeppelin/flink/flink-1.12.3
COPY --from=builder /opt/hadoop-2.7.5 ${HADOOP_HOME}

COPY interpreter.json /opt/zeppelin/conf/
COPY core-site.xml ${HADOOP_HOME}/etc/hadoop/
ADD 1.12_lib /zeppelin/flink/1.12_lib/
ADD depends /zeppelin/flink/depends/

CMD ["bin/zeppelin.sh"]
