package apiobjects

type Replicaset struct {
	Object `yaml:",inline"`
	Spec   ReplicasetSpec `yaml:"spec"`
}
type ReplicasetSpec struct {
	Replicas int           `yaml:"replicas"`
	Selector LabelSelector `yaml:"selector"`
	Template PodTemplate   `yaml:"template"`
	Ready    int           `yaml:"-"`
}
type LabelSelector struct {
	MatchLabels map[string]string `yaml:"matchLabels"`
}
type PodTemplate struct {
	Metadata ObjectMeta `yaml:"metadata"`
	Spec     PodSpec    `yaml:"spec"`
}

func ToPod(podTemplate *PodTemplate) *Pod {
	// podstate :=PodState{}
	// return &Pod{
	// 	Object{
	// 		TypeMeta: TypeMeta{
	// 			ApiVersion: "v1",
	// 			Kind: "Pod",
	// 		},
	// 		Metadata: podTemplate.Metadata,
	// 	},
	// 	Spec: podTemplate.Spec,
	// 	Status: podstate,

	// }
	var pod Pod
	pod.TypeMeta.ApiVersion = "v1"
	pod.TypeMeta.Kind = "Pod"
	pod.ObjectMeta = podTemplate.Metadata
	pod.Spec = podTemplate.Spec
	return &pod
}
