//go:build release

package global

const Host = "192.168.1.12"

const ApiVersion = "v1"
const NFSdir = "/home/nfs"
const ApiserverMountDir = "/mnt/nfs"
const WorkerMountDir = "/mnt/.k8s-volume"
const Nfsserver = "192.168.1.14"

// 系统命名空间，用于存放系统级别的资源。例如：Node、Binding等
const SystemNamespace = "system"
const DefaultNamespace = "default"
const ReplicasetLabel = "app.kubernetes.io/replicas"
