FROM apache/zeppelin:0.9.0
COPY startZeppelin.py /zeppelin/bin/
COPY interpreter.json /zeppelin/conf/
ADD flink-1.12.3-bin-scala_2.11.tgz /zeppelin/flink/
ADD 1.12_lib /zeppelin/flink/1.12_lib/
RUN pip3 install -i https://mirrors.aliyun.com/pypi/simple/ apache-flink

ADD hadoop-2.7.5.tar.gz  /zeppelin/hadoop/
COPY core-site.xml /zeppelin/hadoop/hadoop-2.7.5/etc/hadoop/
ENV HADOOP_HOME /zeppelin/hadoop/hadoop-2.7.5
ENV PATH $PATH:$HADOOP_HOME/bin

CMD ["bin/zeppelin.sh"]

USER root:root
ENTRYPOINT ["/zeppelin/bin/startZeppelin.py"]