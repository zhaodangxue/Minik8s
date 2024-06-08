# DNS

#### 主要功能

Minik8s中的DNS主要功能：

⽀持⽤⼾通过yaml配置⽂件对Service的域名和路径进⾏配置，使得集群内的⽤⼾可以直接通过域名与路径的80端⼝⽽不是虚拟IP来访问映射⾄其下的其他Service的特定端⼝。同时，集群内的Pod可以通过域名与路径的80端⼝访问到该Service。同时⽀持同⼀个域名下的多个⼦路径path对应多个Service。

DNS配置文件的详细格式

```
kind: Dns
apiVersion: v1
name: dns-test1
namespace: default
host: minik8s.com
paths:
  - service: dns-service
    pathName: path1
    port: 22222
  - service: dns-service2
    pathName: path2
    port: 34567
```

#### 具体实现思路

DNS的配置是通过yaml文件进行的，在具体操作的时候将其映射为一个`apiobject` (`DNSRecord`)，支持的字段及其含义和要求文档中的基本一致

当相应的yaml文件被解析后，会生成一个`DNSRecord`对象，该对象的host和nginx server ip的对应会被存储到etcd中，从而实现了域名到ip的映射，并通过`coreDNS`实现了域名解析的动态加载

而nginx server ip下不同path到具体的service的映射则是通过nginx的配置文件实现的，具体的配置文件示例如下

```
worker_processes  5;  ## Default: 1
error_log  ./error.log debug;
pid        ./nginx.pid;
worker_rlimit_nofile 8192;

events {
  worker_connections  4096;  ## Default: 1024
}
http {
    
    server {
        listen 192.168.1.12:80;
        server_name minik8s.com;

        
        location /path1/ {
            access_log /var/log/nginx/access.log;
            proxy_pass http://10.10.0.1:22222/;
        }
        
        location /path2/ {
            access_log /var/log/nginx/access.log;
            proxy_pass http://10.10.0.2:34567/;
        }
        
    }
    
}
```

当DNSRecord被更新后，nginx的配置文件也会同步热更新，从而实现了不同path到不同service的映射

###### 根据域名定位到IP和port的过程

1. 通过域名找到nginx server ip

以上面的nginx配置文件为例，我们要访问`http://minik8s.com:80/path1`，首先会通过`coreDNS`将`minik8s.com`解析为`nginx` server ip `192.168.1.12`，这个过程是通过`coreDNS`自动读取etcd中存储的host和固定的nginx server ip并解析来实现的

2. 根据不同的path找到service

在nginx中根据location的path名字找到对应的service ip和端口，比如在上面的例子中是`10.10.0.1:22222`，这个过程是通过nginx实现的

#### DNS注册的过程

通过yaml文件注册，kubectl接收到注册指令以后，发送http请求给apiserver, apiserver将配置文件转换成dnsrecord数据结构存储到etcd中，同时存储对应的host和固定的nginx server ip到etcd中，以供coreDNS来解析域名，最后读取etcd中所有已经创建的dnsrecord,对nginx配置进行热更新

```powershell
#DNS创建
./kubectl apply -f dnsrecord.yaml

#查看已经创建的DNS
./kubectl get dns

#删除DNS
./kubectl delete dns [对应DNS的NAME]
```

### coreDNS

CoreDNS只需在master节点上运行，是在控制平面上运行的。

> coreDNS对于DNS记录的存储

CoreDNS内部包含一个DNS记录存储后端。该后端被称为CoreDNS的存储插件。CoreDNS存储插件的主要作用是管理DNS记录的存储和检索，例如将DNS记录存储在etcd、Consul或文件系统中，并在必要时从这些后端中检索记录。当客户端查询DNS记录时，CoreDNS将首先从存储插件中检索记录，然后将记录返回给客户端。如果记录不存在，则返回一个相应的错误。此外，存储插件还支持动态DNS更新，允许客户端通过API向CoreDNS添加、删除和修改DNS记录。

#### 不同子路径对应不同的service

使用**nginx**做反向代理，将不同的path路由到不同的ip+port
