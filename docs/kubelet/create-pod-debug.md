### dial: permission denied

containerd.sock拥有者为root，不允许other user rw

目前解决方案：chmod a+rw ...

证明：连接cri rpc确实也是通过unix:///run/containerd/containerd.sock

### 下载pause镜像超时

修改/etc/containerd/config.toml

首先要将默认设置导入config文件 `containerd config default | sudo tee /etc/containerd/config.toml`

然后修改pause默认镜像 `sandbox_image = "registry.aliyuncs.com/google_containers/pause:3.8"`

### cni插件未初始化

RRO[0001] RunPodSandbox error:rpc error: code = Unknown desc = failed to setup network for sandbox "2cbef6782973f1d48e7a69504765a2112ffd83718b9e9f6f01eb7cb4fe7f7db9": cni plugin not initialized

需要安装cni插件二进制的文件，并配置网络，可以使用containerd仓库中的[install-cni](https://github.com/containerd/containerd/blob/main/script/setup/install-cni)

### pod启动成功，container启动时直接rpc EOF

实际上是containerd crash了

通过 `journalctl -u containerd.service`看日志，发现是在makeContainerName中访问了非法地址，最终发现是rpc参数中的SandboxConfig没有填，导致访问了nil
