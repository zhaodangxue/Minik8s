//go:build release

package global

const Host = "106.15.10.160"

const ApiVersion = "v1"
const NFSdir = "/home/nfs"
const ApiserverMountDir = "/home/zbm/nfs"
const Nfsserver = "192.168.1.7"

// 系统命名空间，用于存放系统级别的资源。例如：Node、Binding等
const SystemNamespace = "system"
const DefaultNamespace = "default"
