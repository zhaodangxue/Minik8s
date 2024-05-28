# use basic image
# FROM ubuntu:20.04
# RUN apt-get update && apt-get install -y python3 python3-pip
FROM python:3.10-slim

COPY ./function /function
# set workdir
WORKDIR /function

EXPOSE 8080

# install dependencies from requirements.txt
RUN pip3 install -i https://pypi.tuna.tsinghua.edu.cn/simple --default-timeout=60 --no-cache-dir -r ./requirements.txt

# set env
ENV PATH="/usr/local/bin:${PATH}"

# run watchdog.py, start the function server
CMD ["python3", "watchdog.py"]