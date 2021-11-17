FROM dataworkbench/builder:zeppelin as builder


FROM apache/zeppelin:0.9.0

ENV HADOOP_HOME /zeppelin/hadoop/hadoop-2.7.5
ENV PATH $PATH:$HADOOP_HOME/bin

COPY --from=builder /opt/flink-* /zeppelin/flink/flink-1.12.3/
COPY --from=builder /opt/hadoop-* /zeppelin/hadoop/

COPY startZeppelin.py /zeppelin/bin/
COPY interpreter.json /opt/zeppelin/conf/
COPY 1.12_lib /zeppelin/flink/1.12_lib/
COPY core-site.xml ${HADOOP_HOME}/etc/hadoop/

USER root:root
ENTRYPOINT ["/zeppelin/bin/startZeppelin.py"]

CMD ["bin/zeppelin.sh"]
