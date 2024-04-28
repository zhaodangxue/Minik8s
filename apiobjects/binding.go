package apiobjects

type NodePodBinding struct {
	Object
	Node ObjectRef
	Pod  ObjectRef
}
