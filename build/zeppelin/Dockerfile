FROM apache/zeppelin:0.9.0
COPY startZeppelin.py /zeppelin/bin/
COPY interpreter.json /opt/zeppelin/conf/
ADD flink-1.12.3-bin-scala_2.11.tgz /zeppelin/flink/
ADD 1.12_lib /zeppelin/flink/1.12_lib/

CMD ["bin/zeppelin.sh"]

USER root:root
ENTRYPOINT ["/zeppelin/bin/startZeppelin.py"]
