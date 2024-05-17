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
