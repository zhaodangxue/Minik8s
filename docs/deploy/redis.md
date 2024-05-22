# Redis 跨节点连接失败 Debug

## 可能的原因

### 安全组

### Redis配置

Redis默认不开放远程连接，需要修改Redis的配置文件 `/etc/redis/redis.conf`。

首先需要注释这一行

```
bind 127.0.0.1 ::1
```

这代表只允许来自127.0.0.1的连接

如果使用不带密码的连接，还需要将保护模式关闭

```
protected-mode no
```
