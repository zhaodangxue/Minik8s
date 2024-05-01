package apiobjects

/*
Binding对象用于将Pod绑定到Node。

创建时机：

- 当Pod被调度到Node上时，Scheduler会请求ApiServer创建一个Binding对象，将Pod绑定到Node上。

销毁时机：

- 当Pod被删除时，ApiServer会删除Binding对象。

- 当Pod被调度到其他Node上时，Scheduler会请求ApiServer删除旧的Binding对象，将Pod绑定到新的Node上。

- 当ApiServer发现节点下线时，会删除所有与该节点相关的Binding对象。
*/
type NodePodBinding struct {
	Object
	Node ObjectRef
	Pod  ObjectRef
}