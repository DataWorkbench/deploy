FROM flink:1.12.3-scala_2.11
RUN cp -rf /opt/flink/opt/* /opt/flink/lib/
RUN sed -i 's#http://deb.debian.org#https://mirrors.163.com#g' /etc/apt/sources.list
RUN apt-get update
RUN apt-get install -y python3 python3-pip
RUN ln /usr/bin/python3 /usr/bin/python
RUN pip3 install -i https://mirrors.aliyun.com/pypi/simple/ apache-flink

USER root:root
