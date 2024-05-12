# 安装containerd

sudo apt install containerd

从https://github.com/containernetworking/plugins的release中下载最新的cni插件二进制文件，放入/opt/cni/bin

在/etc/cni/net.d下创建文件10-containerd-net.conflist
内容为
{
  "cniVersion": "1.0.0",
  "name": "containerd-net",
  "plugins": [
    {
      "type": "bridge",
      "bridge": "cni0",
      "isGateway": true,
      "ipMasq": true,
      "promiscMode": true,
      "ipam": {
        "type": "host-local",
        "ranges": [
          [{
            "subnet": "10.88.0.0/16"
          }],
          [{
            "subnet": "2001:4860:4860::/64"
          }]
        ],
        "routes": [
          { "dst": "0.0.0.0/0" },
          { "dst": "::/0" }
        ]
      }
    },
    {
      "type": "portmap",
      "capabilities": {"portMappings": true}
    }
  ]
}

注：目前以非root用户直接运行kubelet/tests/pod_create_test会出错，需要将/run/containerd/containerd.sock设为当前用户可读写

可以使用指令 `sudo chmod a+rw /run/containerd/containerd.sock`
