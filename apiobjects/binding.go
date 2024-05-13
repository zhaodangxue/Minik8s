package apiobjects

/*
Binding对象用于将Pod绑定到Node。

创建时机：

- 当Pod被调度到Node上时，Scheduler会请求ApiServer创建一个Binding对象，将Pod绑定到Node上。

销毁时机：

- 当Pod被删除时，ApiServer会删除Binding对象。

- 当Pod被调度到其他Node上时，Scheduler会请求ApiServer删除旧的Binding对象，将Pod绑定到新的Node上。

- 当ApiServer发现节点下线时，会删除所有与该节点相关的Binding对象。

约束条件：

- 同一个Pod同时最多只能对应一个Binding对象。

	因此，Binding对象的Name字段应该是Pod的Path。
*/
type NodePodBinding struct {
	Node Node
	Pod  Pod
}

func (npd *NodePodBinding) GetBindingPath() string {
	return GetBindingPath(&npd.Pod)
}

func GetBindingPath(pod *Pod) string {
	return "/api" + "/binding" + "/" + pod.Namespace + "/" + pod.Name
}

func (npd *NodePodBinding) Name() string {
	return npd.GetBindingPath() + "/" + npd.Node.ObjectMeta.Name
}

type PVCPodBinding struct {
	PVCName      string
	PVCNamespace string
	Pods         []PodAbstract
}

func (ppb *PVCPodBinding) GetBindingPath() string {
	return "/api" + "/PVCbinding" + "/" + ppb.PVCNamespace + "/" + ppb.PVCName
}
