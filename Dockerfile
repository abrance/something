FROM golang:latest

WORKDIR /something

COPY . /something
# 拷贝凭证文件到 Docker 容器中, 这里要确保 .git-credentials 存在
COPY .git-credentials /root/.git-credentials

RUN git config --global credential.helper 'store --file=/root/.git-credentials' && \
    go mod tidy -x

# 构建应用
#RUN make release
RUN go build -o main cmd/main.go

# EXPOSE 20111

# 指定容器启动时执行的命令
CMD ./main
