package apiobjects

type PersistentVolumeClaim struct {
	Object     `yaml:",inline"`
	Spec       PersistentVolumeClaimSpec `yaml:"spec"`
	PodBinding []PodAbstract             `yaml:"-"`
	PVBinding  PersistentVolumeAbstract  `yaml:"-"`
	Status     PVCPhase                  `yaml:"-"`
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
type PodAbstract struct {
	Podname   string
	Namespace string
}
type PersistentVolumeAbstract struct {
	PVname      string
	PVnamespace string
	PVcapacity  string
	PVpath      string
}
type PVCPhase string

const (
	PVCAvailable PVCPhase = "Available"
	PVCBound     PVCPhase = "Bound"
)
