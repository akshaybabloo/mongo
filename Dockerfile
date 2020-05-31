FROM ubuntu:18.04

ARG go_version=1.14.3
ENV PATH="/usr/local/go/bin:$PATH"

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow

RUN apt-get update && \
    apt-get install -y wget build-essential gnupg curl git

RUN wget https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${go_version}.linux-amd64.tar.gz && \
    rm go${go_version}.linux-amd64.tar.gz

RUN wget -qO - https://www.mongodb.org/static/pgp/server-4.2.asc | apt-key add -
RUN touch /etc/apt/sources.list.d/mongodb-org-4.2.list
RUN echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu bionic/mongodb-org/4.2 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-4.2.list
RUN apt-get update && \
    apt-get install -y mongodb-org systemd && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY entrypoint.sh /entrypoint.sh
COPY . .

ENTRYPOINT ["/entrypoint.sh"]