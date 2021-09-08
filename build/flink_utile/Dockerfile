FROM flink:1.12.3-scala_2.11

RUN cp -r /opt/flink/opt/* /opt/flink/lib/

RUN sed -i 's#deb.debian.org#mirrors.aliyun.com#g' /etc/apt/sources.list
RUN apt-get update && \
    apt-get install -y python3 python3-pip && \
    ln /usr/bin/python3 /usr/bin/python && \
    apt-get clean
RUN pip3 install --no-cache-dir -i https://mirrors.aliyun.com/pypi/simple/ apache-flink==1.12.3

USER root:root
