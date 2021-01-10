FROM apache/zeppelin:0.9.0
COPY startZeppelin.py /zeppelin/bin/
COPY interpreter.json /zeppelin/conf/
ADD flink-1.11.2-bin-scala_2.11.tgz /zeppelin/flink/

CMD ["bin/zeppelin.sh"]

USER root:root
ENTRYPOINT ["/zeppelin/bin/startZeppelin.py"]
