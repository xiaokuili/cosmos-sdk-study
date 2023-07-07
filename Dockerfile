# 使用 Ubuntu 22.04.2 LTS 作为基础镜像
FROM ubuntu:22.04

# 创建一个新的目录来存放二进制文件
RUN mkdir -p /app
WORKDIR /app

# 将本地的二进制文件复制到容器中
COPY ./datafactory/build /app
COPY ./run.sh /app

# 运行二进制文件
CMD ["/app.run.sh"]
