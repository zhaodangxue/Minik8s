#### 安装CNI插件

首先执行 sudo apt install container ,然后安装cni插件的二进制文件，并配置网络，可以使用containerd仓库中的

[install-cni]: https://github.com/containerd/containerd/blob/main/script/setup/install-cni

默认情况下，当安装好cni插件和工具后，使用nerdctl启动了一个容器之后，宿主机上会出现一个cni0的网桥

```powershell
# ip a s
......
3: cni0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default qlen 1000
    link/ether 36:af:7a:4a:d6:12 brd ff:ff:ff:ff:ff:ff
    inet 10.66.0.1/16 brd 10.66.255.255 scope global cni0
       valid_lft forever preferred_lft forever
    inet6 fe80::34af:7aff:fe4a:d612/64 scope link
       valid_lft forever preferred_lft forever

```

#### 创建容器

在执行下一步之前，首先要创建一个容器，在minik8s\kubelet\tests目录下运行go test，即可创建一个默认的busybox镜像的容器，运行过程中的问题可以参考create-pod-debug.md文件。

#### 进入容器

```powershell
#ctr -n k8s.io container ls 这行指令可以看到我们以及创建出来的所有容器，每个容器都有对应的pause容器， -n k8s.io指定了命名空间为k8s.io,我们创建的容器都在这个namespace下，如果直接执行ctr container ls将什么都看不到，因为default的命名空间下是什么都没有的

#ctr -n k8s.io task ls 这行指令可以看到所有正在运行的容器，这里我们应该能看到一个容器在运行，是我们刚刚创建的容器的pause容器
```

这里注意在创建一个容器时，容器的pause容器是默认运行的，但我们创建的容器本身是不会运行的，需要我们手动运行

```powershell
# ctr -n k8s.io task start [容器的ID]
```

然后就可以进入容器了

```powershell
# ctr -n k8s.io tasks exec --exec-id $RANDOM -t [容器的ID] sh  注意别写成pause容器ID，pause容器中什么都没有
```

然后在容器中执行以下命令

```powershell
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: eth0@if8: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1500 qdisc noqueue 
    link/ether 76:cd:1f:fd:1b:50 brd ff:ff:ff:ff:ff:ff
    inet 10.88.0.6/16 brd 10.88.255.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 2001:4860:4860::6/64 scope global 
       valid_lft forever preferred_lft forever
    inet6 fe80::74cd:1fff:fefd:1b50/64 scope link 
       valid_lft forever preferred_lft forever
       
在容器中ping容器宿主机IP地址
/ # ping -c 2 192.168.10.164
PING 192.168.10.164 (192.168.10.164): 56 data bytes
64 bytes from 192.168.10.164: seq=0 ttl=64 time=0.132 ms
64 bytes from 192.168.10.164: seq=1 ttl=64 time=0.044 ms

--- 192.168.10.164 ping statistics ---
2 packets transmitted, 2 packets received, 0% packet loss
round-trip min/avg/max = 0.044/0.088/0.132 ms
```

按照相同的步骤再创建一个容器，同一主机上的两个容器即可相互ping 通。

```
在容器中开启httpd服务
/ # echo "containerd net web test" > /tmp/index.html
/ # httpd -h /tmp

/ # wget -O - -q 127.0.0.1
containerd net web test
/ # exit

在宿主机访问容器提供的httpd服务
[root@localhost scripts]# curl http://10.88.0.2
containerd net web test
```

#### 安装nerdctl工具

使用效果与 docker 命令的语法一致，github 下载链接：https://github.com/containerd/nerdctl/releases

精简 (nerdctl–linux-amd64.tar.gz): 只包含 nerdctl

完整 (nerdctl-full–linux-amd64.tar.gz): 包含 containerd, runc, and CNI 等依赖

原文链接：https://blog.csdn.net/m0_37843156/article/details/128277966

1）安装 nerdctl（精简版）

先去官网下载压缩包，然后将压缩包放到 /opt/nerdctl目录下，压tar -xf nerdctl-版本号-linux-amd64.tar.gz

最后执行

```powershell
# ln -s /opt/nerdctl/nerdctl /usr/local/bin/nerdctl
```

使用以下命令可以看到所有在运行的容器

```powershell
# nerdctl -n k8s.io container ls

CONTAINER ID    IMAGE          COMMAND            CREATED        STATUS    PORTS    NAME            
975eb5a7882e registry.aliyuncs.com/google_containers/pause:3.8  "/pause"  6 hours ago Up                 
ab55fa6740a0    docker.io/library/busybox:latest "sh"   6 hours ago    Up                           
c28a749c37a9    registry.aliyuncs.com/google_containers/pause:3.8 "/pause" 3 hours ago Up                 
d1c85f8d2d99    docker.io/library/busybox:latest  "sh"   3 hours ago    Up 
```

