# Gitlab

目前，我们在Node-1上部署了gitlab，并设置其自动推送到gitee仓库

## Gitlab 启动

Gitlab对机器的性能有一定的要求，我们在其文档所述的2核4G的机器上没有跑起来，建议采用4核，8G以上的配置。

最简单的配置方式是通过docker compose创建docker容器启动

```
version: '3.6'
services:
  web:
    image: 'registry.gitlab.cn/omnibus/gitlab-jh:latest'
    restart: always
    hostname: 'gitlab.unit1112.com'
    environment:
      GITLAB_OMNIBUS_CONFIG: |
        external_url 'https://gitlab.unit1112.com'
        # Add any other gitlab.rb configuration here, each on its own line
    ports:
      - '80:80'
      - '443:443'
      - '23:22'
    volumes:
      - '$GITLAB_HOME/config:/etc/gitlab'
      - '$GITLAB_HOME/logs:/var/log/gitlab'
      - '$GITLAB_HOME/data:/var/opt/gitlab'
    shm_size: '256m'
```

通过docker ps观察容器状态，通常会卡在health(starting)，代表gitlab正在启动。在我们的配置下（4核 16G），每次重新创建容器，启动的过程通常需要1～2m。

启动之后，就可以通过https://your.hostname访问gitlab。

gitlab在开始时会有一个默认的管理员账户root，其初始密码记录在/etc/gitlab/config/initial_root_password中。默认密码24h后会删除，建议在登陆后立即修改密码。

其他用户可以通过注册账户的方式加入。默认情况下，注册新用户经过管理员同意，可以在设置中配置。

## git 连接 Gitlab

这里的问题在于，由于服务器自己也需要向外提供ssh服务，gitlab的ssh端口通常不能开在22，这时如果要通过ssh拉取仓库，如何指定仓库地址就成了一个问题。

笔者没有发现一种可以直接在仓库地址中指定远程端口号的方法，一种比较曲折的解决方式是在.ssh/config中进行配置

```
Host my-gitlab
    HostName xx.xx.xx.xx
    Port 23
```

这样就可以直接使用Host后跟的别名取代仓库url

```
git remote set-url origin git@my-gitlab:username/repository.git
```

## Gitlab Runner踩坑

这里遇到的主要问题在于runner与gitlab进行连接时，由于是通过https连接，需要gitlab提供证书

gitlab在创建时会在/etc/gitlab/config/ssl目录下创建自签名证书，其文件名为hostname.crt。

在执行gitlab-runner register时，可以通过--tls-ca-file选项指定上述证书文件，信任gitlab的自签名证书。

然而问题在于，若域名为ip地址，必须在证书的SAN（subject alternative name）字段中指明，gitlab在创建自签名证书时，没有考虑到域名为ip地址的情况。

因而gitlab-runner在连接时就会报对应的错误。

解决这个问题有两个途径

- 通过openssl，创建一个带有SAN字段的证书，替换原有证书，然后执行gitlab-ctl reconfigure。
- 不使用ip访问，在配置文件中指定gitlab的hostname形如your.gitlab.hostname，通过修改/etc/hosts将your.gitlab.hostname解析到正确的ip，绕开这个问题。

我们选择了后一种。

## Gitlab Runner CI/CD

### log不完整

gitlab runner在执行job时，在gitlab web界面中显示的log是不完整的，常常缺失后半部分，原因不明。

解决方法是在.gitlab-ci.yml中定义的job中添加after-scripts字段

```
your-job: 
  stage: test
  after_script:
    - sleep 6
  script:
    - xxx
```

通过sleep 6拖延时间，可以缓解这个问题。
