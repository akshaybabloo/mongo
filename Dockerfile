FROM mongo:4-bionic

ARG go_version=1.15.6
ENV PATH="/usr/local/go/bin:$PATH"

ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow

RUN apt-get update && \
    apt-get install -y wget build-essential gnupg curl git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN wget -q https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${go_version}.linux-amd64.tar.gz && \
    rm go${go_version}.linux-amd64.tar.gz

COPY entrypoint.sh /entrypoint.sh
COPY . .

ENTRYPOINT ["/entrypoint.sh"]
