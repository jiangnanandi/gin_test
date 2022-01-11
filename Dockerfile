FROM ubuntu:latest

RUN apt update && apt install openssh-server sudo -y

# install golang
# 设置环境变量，所有操作都是非交互式的
ENV DEBIAN_FRONTEND noninteractive
ENV GO_USER=golang
ENV GO_LOG_DIR=/var/log/golang
# 这里的GOPATH路径是挂载的 app 项目的目录
ENV GOPATH=/home/golang/
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

#ENV GOLANG_VERSION 1.7.4
#ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
#ENV GOLANG_DOWNLOAD_SHA256 47fda42e46b4c3ec93fa5d4d4cc6a748aa3f9411a2a2b7e08e3a6d80d753ec8b

# 替换 sources.list 的配置文件，并复制配置文件到对应目录下面。
# 这里使用的AWS国内的源，也可以替换成其他的源（例如：阿里云的源）
#COPY sources.list /etc/apt/sources.list

# 安装基础工具
RUN sudo apt-get clean
RUN sudo rm -rf /var/lib/apt/lists/*
RUN sudo apt-get update
RUN sudo apt-get install -y vim wget curl git

# 使用apt方式安装golang
RUN sudo apt-get -y install golang

# 下载并安装golang
#RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
#	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
#	&& tar -C /usr/local -xzf golang.tar.gz \
#	&& rm golang.tar.gz

# 创建用户和创建目录
RUN set -x && useradd $GO_USER && mkdir -p $GO_LOG_DIR $GOPATH && chown $GO_USER:$GO_USER $GO_LOG_DIR $GOPATH

WORKDIR $GOPATH

# install sshd
RUN service ssh start

RUN echo 'root:123456' | chpasswd

RUN sed -i '$a PermitRootLogin yes' /etc/ssh/sshd_config

RUN sed -i '$a UsePAM no' /etc/ssh/sshd_config

EXPOSE 22

COPY . /home/app

RUN cd /home/app

RUN go get -d -v ./...

CMD ["/usr/sbin/sshd", "-D"]


# Use demo
# docker build -t go-remote .
# docker run -d -p 2020:22 CONTAINER_ID
# ssh root@127.0.0.1 -p 2020
