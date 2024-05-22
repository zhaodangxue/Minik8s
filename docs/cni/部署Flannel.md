##### 在node2 上部署etcd和etcdctl

```
获取安装包
wget https://storage.googleapis.com/etcd/v3.4.13/etcd-v3.4.13-linux-amd64.tar.gz
tar -zxvf etcd-v3.4.13-linux-amd64.tar.gz -C /usr/local/
ls /usr/local/etcd-v3.4.13-linux-amd64/

软链接
ln -s /usr/local/etcd-v3.4.13-linux-amd64/etcd* /usr/local/bin/

启动etcd
/usr/local/etcd-v3.4.13-linux-amd64/etcd

配置 etcd
mkdir /etc/etcd

vim /etc/etcd/etcd.conf
#[member]
ETCD_NAME="etcd"
ETCD_DATA_DIR="/data/etcd/default.etcd"
ETCD_LISTEN_CLIENT_URLS="http://192.168.45.23:2379,http://127.0.0.1:2379"
ETCD_ADVERTISE_CLIENT_URLS=http://192.168.45.23:2379,http://127.0.0.1:2379
ETCD_ENABLE_V2=true 

配置文件解析
    ETCD_NAME 节点名称
    ETCD_DATA_DIR 数据目录
    ETCD_LISTEN_CLIENT_URLS 客户端访问监听地址
    ETCD_ADVERTISE_CLIENT_URLS 客户端通告地址
    ETCD_ENABLE_V2 ETCD 3.4 版本 ETCDCTL_API=3 和 --enable-v2=false 成为了默认配置,
        如要使用 v 2 版本, 需要 ETCD_ENABLE_V 2=true，否则会报错“404 page not found”


配置 etcd.service
vim /usr/lib/systemd/system/etcd.service
[Unit]
Description=Etcd Service
Documentation=https://coreos.com/etcd/docs/latest/
After=network.target

[Service]
Type=notify
ExecStart=/usr/local/bin/etcd
EnvironmentFile=-/etc/etcd/etcd.conf
Restart=on-failure
RestartSec=10
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
Alias=etcd3.service

启动服务
systemctl start etcd
systemctl enable etcd
ss -antup|grep etcd

测试
etcdctl --endpoints=http://192.168.1.12:2379 put foo "bar"
etcdctl --endpoints=http://192.168.1.12:2379 get foo
```

#### 将 flannel 网络的配置信息保存到 etcd

在node2上执行

```
export ETCDCTL_API=3 //所有node必须执行

node-2:~# etcdctl --endpoints "http://192.168.1.12:2379" put /coreos.com/network/config '{"NetWork":"10.2.0.0/16","SubnetMin":"10.2.1.0","SubnetMax": "10.2.20.0","Backend": {"Type": "vxlan"}}'
OK

node-2:~#etcdctl --endpoints "http://192.168.1.12:2379" get /coreos.com/network/config
/coreos.com/network/config
{"NetWork":"10.2.0.0/16","SubnetMin":"10.2.1.0","SubnetMax": "10.2.20.0","Backend": {"Type": "vxlan"}}
```

node1 上只要安装etcdctl

在node1上执行

```
export ETCDCTL_API=3 //通过apt安装，不执行这行代码会报错找不到key

node-1:# export ETCDCTL_API=3
node-1:#etcdctl --endpoints "http://192.168.1.12:2379" get /coreos.com/network/config
/coreos.com/network/config
{"NetWork":"10.2.0.0/16","SubnetMin":"10.2.1.0","SubnetMax": "10.2.20.0","Backend": {"Type": "vxlan"}}

node-1:/home# etcdctl --endpoints=http://192.168.1.12:2379 put foo "bar"
OK
node-1:/home# etcdctl --endpoints=http://192.168.1.12:2379 get foo
foo
bar
```



安装Flannel插件，由于高版本cni-plugins中已经删去了flannel插件,所以可以下载0.9.1的版本压缩包解压，将里面的flannel文件移动到存放cni-plugins的/opt/cni/bin中

然后去官网下载flanneld-amd64文件

```
wget https://github.com/flannel-io/flannel/releases/latest/download/flanneld-amd64 && chmod +x flanneld-amd64
```

node启动`./flanneld-amd64 -etcd-endpoints=http://192.168.1.12:2379 -iface=ens3`

这里ens3是主机上能和外界通信的网卡，如果不设置flannel也会自动找

```
node1上的
17: flannel.1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1400 qdisc noqueue state UNKNOWN group default 
    link/ether 2a:f1:30:e9:8b:2f brd ff:ff:ff:ff:ff:ff
    inet 10.2.6.0/32 scope global flannel.1
       valid_lft forever preferred_lft forever
    inet6 fe80::28f1:30ff:fee9:8b2f/64 scope link 
       valid_lft forever preferred_lft forever
       
node2上的
17: flannel.1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1400 qdisc noqueue state UNKNOWN group default 
    link/ether 82:39:3e:02:c8:82 brd ff:ff:ff:ff:ff:ff
    inet 10.2.10.0/32 scope global flannel.1
       valid_lft forever preferred_lft forever
    inet6 fe80::8039:3eff:fe02:c882/64 scope link 
       valid_lft forever preferred_lft forever
       
再结合放在etcd中的配置
{"NetWork":"10.2.0.0/16","SubnetMin":"10.2.1.0","SubnetMax": "10.2.20.0","Backend": {"Type": "vxlan"}}

good!

这时候查看网络设置，发现没有flannel，需要配置
root@node-1:/home# nerdctl network ls
NETWORK ID      NAME              FILE
                containerd-net    /etc/cni/net.d/10-containerd-net.conflist
                mynet             /etc/cni/net.d/10-mynet.conf
                lo                /etc/cni/net.d/99-loopback.conf
17f29b073143    bridge            /etc/cni/net.d/nerdctl-bridge.conflist
                host              
                none  
                
                
root@node-1:/home# nano /etc/cni/net.d/10-flannel.conflist            
{
  "name": "flannel",
  "cniVersion": "0.3.1",
  "plugins": [
    {
      "type": "flannel",
      "delegate": {
        "isDefaultGateway": true
      }
    },
    {
      "type": "portmap",
      "capabilities": {
        "portMappings": true
      }
    }
  ]
}

root@node-1:/home# nerdctl network ls
NETWORK ID      NAME              FILE
                containerd-net    /etc/cni/net.d/10-containerd-net.conflist
                flannel           /etc/cni/net.d/10-flannel.conflist
                mynet             /etc/cni/net.d/10-mynet.conf
                lo                /etc/cni/net.d/99-loopback.conf
17f29b073143    bridge            /etc/cni/net.d/nerdctl-bridge.conflist
                host              
                none 
出现了！good！
```

由于containerd-net网络在创建时会使用cni0，而flannel网络在创建容器时也会使用cni0，会出现冲突，所以这里先把mynet,lo ,containerd-net统统干掉.

###### 清掉所有容器
nerdctl ps -aq |xargs nerdctl rm 
###### 停掉cni0接口
ifconfig cni0 down 
###### 删除cni0接口
apt install -y bridge-utils
brctl delbr cni0

###### 移除containerd-net配置文件
mv /etc/cni/net.d/10-containerd-net.conflist /tmp/

```
root@node-1:/home# mv /etc/cni/net.d/10-containerd-net.conflist /tmp/
root@node-1:/home# mv /etc/cni/net.d/10-mynet.conf /tmp/
root@node-1:/home# mv /etc/cni/net.d/99-loopback.conf /tmp/
root@node-1:/home# ls  /etc/cni/net.d
10-flannel.conflist  nerdctl-bridge.conflist
root@node-1:/home# nerdctl ps -a
CONTAINER ID    IMAGE    COMMAND    CREATED    STATUS    PORTS    NAMES
root@node-1:/home# nerdctl network ls
NETWORK ID      NAME       FILE
                flannel    /etc/cni/net.d/10-flannel.conflist
17f29b073143    bridge     /etc/cni/net.d/nerdctl-bridge.conflist
                host       
                none 
```

```
root@node-1:/home# nerdctl run -d --net flannel --name flannel busybox:1.28
44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4

17: flannel.1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1400 qdisc noqueue state UNKNOWN group default 
    link/ether 2a:f1:30:e9:8b:2f brd ff:ff:ff:ff:ff:ff
    inet 10.2.6.0/32 scope global flannel.1
       valid_lft forever preferred_lft forever
    inet6 fe80::28f1:30ff:fee9:8b2f/64 scope link 
       valid_lft forever preferred_lft forever
18: cni0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN group default qlen 1000
    link/ether 5a:87:5f:42:e2:d3 brd ff:ff:ff:ff:ff:ff
    inet 10.2.6.1/24 brd 10.2.6.255 scope global cni0
       valid_lft forever preferred_lft forever
    inet6 fe80::5887:5fff:fe42:e2d3/64 scope link 
       valid_lft forever preferred_lft forever
这个cni0是flannel创建的
```

```
接下来就是非常抽象的容器启动，非常之脑瘫
root@node-1:/home# nerdctl exec -it 44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4 sh
FATA[0000] cannot exec in a stopped state: unknown      //提示容器没启动
root@node-1:/home# nerdctl ps -a
CONTAINER ID    IMAGE                             COMMAND    CREATED               STATUS                           PORTS    NAMES
44c47a1e635c    docker.io/library/busybox:1.28    "sh"       About a minute ago    Exited (0) About a minute ago             flannel
root@node-1:/home# nerdctl start 44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4
44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4 //你以为启动成功了，但是没有
root@node-1:/home# nerdctl exec -it 44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4 sh
FATA[0000] cannot exec in a stopped state: unknown      //继续提示容器没启动
root@node-1:/home# ctr task start 44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4//尝试用ctr启动，报错task已经存在
ctr: task 44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4: already exists
root@node-1:/home# ctr task delete 44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4
root@node-1:/home# ctr task start 44c47a1e635c04680b9c9e7b242f1070c86ad1c532d42739e6f2678a29c7ddd4 
//用ctr先删除再启动，成功，但会卡在这
```

```
node1进入用flannel网络创建的容器
root@node-1:~# nerdctl exec -it 44c47a1e635c sh
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: eth0@if21: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1400 qdisc noqueue 
    link/ether d2:da:30:e5:49:49 brd ff:ff:ff:ff:ff:ff
    inet 10.2.6.4/24 brd 10.2.6.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::d0da:30ff:fee5:4949/64 scope link 
       valid_lft forever preferred_lft forever

node2进入用flannel网络创建的容器
root@node-2:/home# nerdctl exec -it 712940a26f75 sh
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: eth0@if38: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1400 qdisc noqueue 
    link/ether 76:7e:d0:ea:3b:ef brd ff:ff:ff:ff:ff:ff
    inet 10.2.10.17/24 brd 10.2.10.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::747e:d0ff:feea:3bef/64 scope link 
       valid_lft forever preferred_lft forever

node2进入另一个用flannel网络创建的容器
root@node-2:/home# nerdctl exec -it 130d8fdea619 sh
/ # ip a s
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: eth0@if40: <BROADCAST,MULTICAST,UP,LOWER_UP,M-DOWN> mtu 1400 qdisc noqueue 
    link/ether c6:99:a0:1c:b3:d8 brd ff:ff:ff:ff:ff:ff
    inet 10.2.10.19/24 brd 10.2.10.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::c499:a0ff:fe1c:b3d8/64 scope link 
       valid_lft forever preferred_lft forever
/ # ping 10.2.10.19
PING 10.2.10.19 (10.2.10.19): 56 data bytes
64 bytes from 10.2.10.19: seq=0 ttl=64 time=0.051 ms
64 bytes from 10.2.10.19: seq=1 ttl=64 time=0.073 ms
64 bytes from 10.2.10.19: seq=2 ttl=64 time=0.070 ms
^C
--- 10.2.10.19 ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.051/0.064/0.073 ms
/ # ping 10.2.10.17
PING 10.2.10.17 (10.2.10.17): 56 data bytes
64 bytes from 10.2.10.17: seq=0 ttl=64 time=0.145 ms
64 bytes from 10.2.10.17: seq=1 ttl=64 time=0.135 ms
64 bytes from 10.2.10.17: seq=2 ttl=64 time=0.116 ms
^C
--- 10.2.10.17 ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.116/0.132/0.145 ms


发现同一个节点上的容器能ping通，不同节点上的不行???????????不知道为何
```

最后发现是flanneld这个进程被我关了

可以按照下面的指令检查flanneld进程是否正常

```
root@node-2:/home# ps aux | grep flanneld
root      580074  0.0  0.8 1710764 33120 pts/2   Sl+  10:20   0:00 ./flanneld-amd64 -etcd-endpoints=http://192.168.1.12:2379 -iface=ens3 //有这行说明没问题
root      580147  0.0  0.0   8160   720 pts/0    S+   10:20   0:00 grep --color=auto flanneld //只有这行说明就没有启动
```

两个node都再次执行`./flanneld-amd64 -etcd-endpoints=http://192.168.1.12:2379 -iface=ens3`

最后成功哩！！！！！！！

真是一场酣畅淋漓的吃屎啊



（配置flanneld服务，不用手动在前台执行）

```
sudo nano /etc/systemd/system/flanneld.service

[Unit]
Description=Flannel Network Fabric for Kubernetes
Documentation=https://github.com/flannel-io/flannel
After=network.target

[Service]
User=root
ExecStart=/home/flanneld-amd64 -etcd-endpoints=http://192.168.1.12:2379 -iface=ens3
          #(flanneld-amd64安装的位置，我是在home目录下）
Restart=on-failure
# 限制重启时间
RestartSec=5s
Type=notify
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
```

保存并关闭文件。然后，重新载入 systemd 配置以使新增加的服务生效：

```
sudo systemctl daemon-reload
```

启用该服务以使其在启动时自动运行，并立即手动启动服务：

```
sudo systemctl enable flanneld
sudo systemctl start flanneld
```

您可以使用以下命令来检查 `flanneld` 服务的状态，确保它正在运行：

```
sudo systemctl status flanneld
```
