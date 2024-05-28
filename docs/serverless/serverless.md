# Serverless Gateway

参考OpenFaaS的实现，[文档](https://docs.openfaas.com/architecture/invocations/)

## HTTP调用模型

全异步，轮询获取结果

HTTP发往Serverless Gateway，gateway收到后，返回给用户一个用于获取结果的凭证call id，目前使用uuid

// TODO：考虑 call id collision

Serveless Gateway将call id和参数传给函数，对于返回值，判断其处于workflow中的何种环节，从而决定是更新call id对应的返回值，还是进一步调用workflow中的下一个函数。

用户可以通过一个带call id的获取结果请求，获取workflow的最终结果。

## Event调用模型

## Function

组成：对应一个service和一个replica set

replica set当中是function pod

### Function pod

成分：

- 固定的成分：Function Daemon（Watchdog）
  - 开端口监听从中心调度来的请求，发送返回值给中心
  - HealthReport，以便service知道何时可以发送流量给它
- 运动的成分：python
  - python标准：支持哪些输入、输出

能力：

- 能同时处理多个请求

## Gateway

Serverless gateway要做的事情

- 监听serverless创建，构建function对应的image
- 创建function对应的replicaset和service
- 像hpa一样控制replicaset的pod数量
- 监听用户请求，负载均衡转发给正在运行的function pod，如果没有，就创建一个
- 支持DAG

创建image
