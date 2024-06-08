# Prometheus

### Prometheus功能

Minik8s使用Promethus实现简单的⽇志和监控功能，以提⾼系统的可观测性。具体的功能如下：

1. 实现对集群中各节点的资源的监控。Minik8s能够采集集群中每个节点(主机)的运⾏指标，如CPU、内存、⽹络等信息，并通过Promethus UI显⽰各节点的资源。

2. 实现对⽤⼾程序的⽇志监控。⽤⼾可以通过Promethus Client Library在⾃⼰的程序中⾃定义需要采集的指标，主动暴露相应的信息。创建⼀个Pod运⾏该程序，集群中的Promethus能够监控到这些⽤⼾⾃定义的指标，并显⽰在Promethus UI。

3. 集群中Promethus的⾃动服务发现。集群中的Promethus⾃动发现Node的加⼊和退出、Pod的创建和删除，在不重启Prometheus服务的情况下动态地发现需要监控的Target实例信息。
   - 当新的节点加⼊集群时，集群中的Promethus可以⾃动监控新加⼊的节点。当节点退出集群时，⾃动停⽌对该节点的监控。
   - 当⽤⼾创建⼀个Pod，且其中的程序主动暴露⾃定义的Promethus指标时，集群中的Promethus可以⾃动监控Pod程序中⾃定义的指标。当Pod被删除时，⾃动停⽌对这些⽤⼾⾃定义指标的监控。

4.  利⽤Grafana为集群监控提供更丰富的可视化展⽰⽅式。

### 实现思路

1. 首先在master节点上部署Prometheus服务器

2. 在每个node上部署node_exporter

3. 在master节点的配置文件中加入node的IP和固定端口，具体配置如下

   ```
   # my global config
   global:
     scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
     evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
     # scrape_timeout is set to the global default (10s).
   
   # Alertmanager configuration
   alerting:
     alertmanagers:
       - static_configs:
           - targets:
             # - alertmanager:9093
   
   # Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
   rule_files:
     # - "first_rules.yml"
     # - "second_rules.yml"
   
   # A scrape configuration containing exactly one endpoint to scrape:
   # Here it's Prometheus itself.
   scrape_configs:
     # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
     - job_name: "prometheus"
   
       # metrics_path defaults to '/metrics'
       # scheme defaults to 'http'.
   
       static_configs:
         - targets: ["localhost:9090"]
   
     - job_name: "workers" #定义名称
   
       static_configs:
         - targets: ["192.168.1.14:9100"] //两个初始workers节点,9100为node固定的端口
         - targets: ["192.168.1.15:9100"]
   
   ```

4. 对于要监听的Pod来说，我们要求用户在创建相应的Pod时要在labels中加入“log: prometheus”，同时在容器的端口中指定某个开放的Port的name为prometheus，这样prometheus_controller就会自动将这个Pod加入配置文件中

   ```
   # my global config
   global:
     scrape_interval: 15s # Set the scrape interval to every 15 seconds. Default is every 1 minute.
     evaluation_interval: 15s # Evaluate rules every 15 seconds. The default is every 1 minute.
     # scrape_timeout is set to the global default (10s).
   
   # Alertmanager configuration
   alerting:
     alertmanagers:
       - static_configs:
           - targets:
             # - alertmanager:9093
   
   # Load rules once and periodically evaluate them according to the global 'evaluation_interval'.
   rule_files:
     # - "first_rules.yml"
     # - "second_rules.yml"
   
   # A scrape configuration containing exactly one endpoint to scrape:
   # Here it's Prometheus itself.
   scrape_configs:
     # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
     - job_name: "prometheus"
   
       # metrics_path defaults to '/metrics'
       # scheme defaults to 'http'.
   
       static_configs:
         - targets: ["localhost:9090"]
   
     - job_name: "workers" #定义名称
   
       static_configs:
         - targets: ["192.168.1.14:9100"] //两个初始workers节点,9100为node固定的端口
         - targets: ["192.168.1.15:9100"]
         - targets: ["10.10.17.100:2112"] //自己创建的需要监听的Pod
   ```

   

5. 通过PrometheusController来自动发现和修改需要监听的对象的变化，包括node和pod

   PrometheusController会定期获取集群中所有的node并筛选出需要监听的Pod，然后和正在监听的node数组和Pod数组进行比较，如果发现有变化，就对Prometheus服务进行热更新（修改配置文件，然后执行reload），以此来实现自动监控