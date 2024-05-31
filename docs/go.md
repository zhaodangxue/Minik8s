# Go

## Go 安装

从[官网](go.dev)直接下载压缩包，添加bin目录到PATH下

## Go Get 卡住

设置proxy

```shell
go env -w GOPROXY=https://goproxy.cn
```

## GLIBC_XX Not Found

编译时关闭CGO，避免编译产物需要动态链接glibc

添加环境变量CGO_ENABLED=0
