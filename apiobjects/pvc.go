package apiobjects

type PersistentVolumeClaim struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata"`
	Spec       PersistentVolumeClaimSpec `yaml:"spec"`
}
type PersistentVolumeClaimSpec struct {
	AccessModes      []string  `yaml:"accessModes"`
	Resources        Resources `yaml:"resources"`
	StorageClassName string    `yaml:"storageClassName"`
}
type Resources struct {
	Requests Requests `yaml:"requests"`
}
type Requests struct {
	Storage string `yaml:"storage"`
}
