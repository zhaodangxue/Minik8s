//go:build dev

package global

const Host = "127.0.0.1"
const EtcdAndRedisHost = "106.15.10.160"

const ApiVersion = "v1"
const NFSdir = "/home/nfs"
const ApiserverMountDir = "/home/zbm/k8s-nfs"
const WorkerMountDir = "/home/zbm/.k8s-volume"
const Nfsserver = "127.0.0.1"

// 系统命名空间，用于存放系统级别的资源。例如：Node、Binding等
const SystemNamespace = "system"
const DefaultNamespace = "default"
const ReplicasetLabel = "app.kubernetes.io/replicas"
