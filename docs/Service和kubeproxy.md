# Service和Kubeproxy

### Service功能

Service⽀持多个pod的通信，对外提供service的虚拟ip。⽤⼾能够通过虚拟ip访问Service，由minik8s将请求具体转发⾄对应的pod，使⽤IPVS控制流量的转发。Service可以看作⼀组pod的前端代理。Service通过selector筛选包含的pod，并将发往service的流量通过随机/round robin等策略负载均衡到这些Pod上。

在符合selector筛选条件的Pod更新时（如Pod加⼊和Pod被删除），service会动态更新（如将被删除的Pod移出管理和将新启动的Pod纳⼊管理）。

⽤⼾能够通过yaml配置⽂件来创建service

Service配置文件格式如下：

```
apiVersion: v1
kind: Service
metadata:
  name: HelloService
  namespace: default
spec:
  selector:
    app: hello
  type: ClusterIP
  ports:
    - name: HelloPort
      protocol: TCP
      port: 12345 # 对外暴露的端口
      targetPort: hello # 转发的端口的名字，pod对应的端口名字
```

**NodePort** **Service**

service在不同节点上开相同的端口，用户以同一端口号访问任意一个node都可以访问到该服务。

NodePort类型Service的配置文件

```
apiVersion: v1
kind: Service
metadata:
  name: HelloService2
  namespace: default
spec:
  selector:
    app: hello
  type: NodePort
  ports:
    - name: HelloPort
      protocol: TCP
      port: 23456 # 对外暴露的端口
      targetPort: hello # 转发的端口的名字，pod对应的端口名字
```

### 实现思路

#### Kubeproxy

Kubeproxy运行在各个节点上， 主要是监听master节点上的svccontroller，当监听到service或对应的endpoint变化时，修改本地的ipvs规则。同时Kubeproxy还会定期检查节点上已经创建的service规则和etcd中存储的规则，保持一致性

#### Master节点上的ServiceController

ServiceController会监听apiserver，当监听到service创建时，会为Service分配一个“持久化的”集群内的IP（Service的IP是持久化的，就是Service对应的Pod挂了也不会变），然后筛选符合条件的Pod，创建对应的endpoints，将Service信息和对应的endpoints信息发给apiserver，apiserver将这些信息都存储到etcd中，最后通知各个节点上的kubeproxy创建相应的IPVS规则，最终实现访问Service的虚拟IP转发到相应Pod上的功能

ServiceController同时也会监听Pod状态的变化（创建，删除等），并检查变化的Pod符合哪些已经创建的Service，然后修改相应的信息（包括etcd中存储的service对应的endpoints信息，各个节点上的ipvs规则等）

### Service注册过程

通过yaml文件注册，kubectl接收到注册指令以后，发送http请求给apiserver, apiserver将配置文件转换成service数据结构，通过Redis发送给ServiceController，ServiceController监听到Redis中有关Service创建的信息，读取Service信息，为Service分配“持久化的”集群内唯一的ClusterIP，然后筛选符合条件的Pod，创建对应的endpoint数据结构，最后通过http将完整的Service信息和endpoint信息发送给apiserver，apiserver将这些信息都存储到etcd中，然后通知各个节点上的kubeproxy创建相应的IPVS规则。

```
#service创建
./kubectl apply -f service.yaml

#查看已经创建的service
./kubectl get service

#删除service
./kubectl delete service [对应service的NAME]
```

#### ipvsadm

kubeproxy中用来创建ipvs规则的工具，运行在用户态，提供简单的CLI接口进行ipvs配置

kubeproxy在创建一个Service的ipvs规则时流程如下

1. 打开ipvs的conntrack

```undefined
sysctl net.ipv4.vs.conntrack=1
```

2. 添加Service的虚拟ip

```csharp
ipvsadm -A -t 10.10.0.1:8410 -s rr
```

3. 把虚拟ip地址添到本地flannel.1网卡

```csharp
ip addr add 10.10.0.1/24 dev flannel.1
```

4. 为虚拟ip添加endpoint（真正提供服务的节点，对应的pod的ip和开放的端口）

```cpp
ipvsadm -a -t 10.10.0.1:8410 -r 10.2.17.53:12345 -m
```

5. 添加SNAT功能

```shell
iptables -t nat -A POSTROUTING -m ipvs  --vaddr 10.10.0.1 --vport 8410 -j MASQUERADE
```

删除命令：`ipvsadm -D -t 10.10.0.1:8410`

查看所有规则：`ipvsadm -Ln`

##### 下载ipvsadm

`apt install ipvsadm`