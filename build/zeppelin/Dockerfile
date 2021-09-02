FROM dataworkbench/zeppelin-builder as builder

FROM apache/zeppelin:0.9.0

COPY --from=builder /opt/flink-* /zeppelin/flink/
COPY --from=builder /opt/hadoop-* /zeppelin/hadoop/

COPY startZeppelin.py /zeppelin/bin/
COPY interpreter.json /opt/zeppelin/conf/
ADD 1.12_lib /zeppelin/flink/1.12_lib/

COPY core-site.xml /zeppelin/hadoop/hadoop-2.7.5/etc/hadoop/

ENV HADOOP_HOME /zeppelin/hadoop/hadoop-2.7.5
ENV PATH $PATH:$HADOOP_HOME/bin

CMD ["bin/zeppelin.sh"]

USER root:root
ENTRYPOINT ["/zeppelin/bin/startZeppelin.py"]
