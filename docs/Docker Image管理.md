# Docker Image管理

##### 安装Docker

参考文章https://mirror.tuna.tsinghua.edu.cn/help/docker-ce/?eqid=8bc57c10006dcb2100000002647b03a4

```
export DOWNLOAD_URL="https://mirrors.tuna.tsinghua.edu.cn/docker-ce"
# 使用curl或者wget，执行一个即可
# 如您使用 curl
curl -fsSL https://get.docker.com/ | sh
# 如您使用 wget
wget -O- https://get.docker.com/ | sh
```



如果你过去安装过 docker，先删掉：

```
for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do apt-get remove $pkg; done
```

安装依赖

```
apt-get update
apt-get install ca-certificates curl gnupg
```

信任 Docker 的 GPG 公钥并添加仓库：

```
install -m 0755 -d /etc/apt/keyrings

curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://mirrors.tuna.tsinghua.edu.cn/docker-ce/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null
```

最后安装

```
apt-get update
apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

#查看docker版本信息
docker version
```



##### Docker registry容器 环境搭建

> 用于镜像的存储和分发

1. 拉取镜像

```shell
docker pull registry
```

2. 启动registry容器

```shell
docker run -d -p 5000:5000 --restart=always --name registry registry
```

3. 验证容器是否启动成功

```shell
docker ps
```

4. 验证docker registry的功能

```
docker pull testcontainers/helloworld
docker tag testcontainers/helloworld localhost:5000/helloworld
docker push localhost:5000/helloworld

# 拉取本地镜像
docker pull localhost:5000/helloworld:latest
# 运行对应的容器
docker run -it --rm localhost:5000/helloworld:latest
```

运行效果检验：

```shell
2024/05/27 09:01:25 DELAY_START_MSEC: 0
2024/05/27 09:01:25 Sleeping for 0 ms
2024/05/27 09:01:25 Starting server on port 8080
2024/05/27 09:01:25 Sleeping for 0 ms
2024/05/27 09:01:25 Starting server on port 8081
2024/05/27 09:01:25 Ready, listening on 8080 and 8081
```

5. 验证是否可以使用containerd运行镜像

```
# 从本地registry中拉取镜像
docker pull localhost:5000/helloworld:latest
# 保存镜像
docker save localhost:5000/helloworld:latest -o helloworld.tar
# 导入镜像
ctr i import helloworld.tar
# 查看镜像
ctr i ls
# 输出如下：
# REF                              TYPE    ...
# localhost:5000/helloworld:latest application/vnd.docker.distribution.manifest.v2+json
# 根据ref信息，运行镜像
ctr run --rm -t localhost:5000/helloworld:latest helloworld
```

运行效果检验：

```
2024/05/27 09:09:56 DELAY_START_MSEC: 0
2024/05/27 09:09:56 Sleeping for 0 ms
2024/05/27 09:09:56 Starting server on port 8080
2024/05/27 09:09:56 Sleeping for 0 ms
2024/05/27 09:09:56 Starting server on port 8081
2024/05/27 09:09:56 Ready, listening on 8080 and 8081
```



#### docker registry 对应的命令

1. 启动 Docker Registry 容器：确保 Docker Registry 容器正在运行。如果尚未启动，请使用以下命令启动容器：

   ```shell
   docker run -d -p 5000:5000 --restart=always --name registry registry:2
   ```

   这将在后台运行一个 Registry 容器，并将容器的 5000 端口映射到主机的 5000 端口。

2. 构建镜像并标记：使用 Docker CLI 构建一个新的镜像，并为该镜像添加 Registry 的地址和标签。例如，假设你有一个名为 `myimage` 的镜像，可以执行以下命令：

   ```shell
   docker build -t myimage .
   docker tag myimage localhost:5000/myimage:latest
   ```

   这将构建 `myimage` 镜像并为其添加 `localhost:5000` Registry 的地址和 `latest` 标签。

3. 推送镜像到 Registry：使用 `docker push` 命令将镜像推送到 Registry。执行以下命令：

   ```shell
   docker push localhost:5000/myimage:latest
   ```

   这将把镜像推送到 `localhost:5000` Registry。

4. 拉取镜像：使用 `docker pull` 命令从 Registry 拉取镜像。执行以下命令：

   ```shell
   docker pull localhost:5000/myimage:latest
   ```

   这将从 Registry 拉取 `myimage` 镜像的最新版本。

5. 运行容器：使用拉取的镜像运行一个容器来验证容器存储的功能。执行以下命令：

   ```shell
   docker run -it --rm localhost:5000/myimage:latest
   ```

   这将在容器中运行 `myimage` 镜像，并在终端中显示容器的输出。如果容器成功运行并显示预期的输出，说明容器存储功能正常。

6. 在另外一台机器上使用这台机器上的镜像

- 设置运行私有仓库的节点IP为可信任的

```shell
nano /etc/docker/daemon.json
#添加以下内容
{
  "insecure-registries": ["192.168.1.15:5000"]
}
```

- 重启服务

```shell
systemctl daemon-reload
systemctl restart docker
```

- 从registry中拉取镜像

```shell
docker pull 192.168.1.15:5000/helloworld:latest
```

