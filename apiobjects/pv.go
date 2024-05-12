package apiobjects

type PersistentVolume struct {
	Object `yaml:",inline"`
	Spec   PersistentVolumeSpec `yaml:"spec"`
	Status PVPhase              `yaml:"-"`
}
type PersistentVolumeSpec struct {
	Capacity                      capacity                      `yaml:"capacity"`
	VolumeMode                    string                        `yaml:"volumeMode"`
	AccessModes                   []string                      `yaml:"accessModes"`
	PersistentVolumeReclaimPolicy string                        `yaml:"persistentVolumeReclaimPolicy"`
	StorageClassName              string                        `yaml:"storageClassName"`
	MountOptions                  []string                      `yaml:"mountOptions"`
	NFS                           NFS                           `yaml:"nfs"`
	PVCBinding                    PersistentVolumeClaimAbstract `yaml:"-"` // 用于绑定PVC
	PVPath                        string                        `yaml:"-"` // 用于绑定PV
}
type capacity struct {
	Storage string `yaml:"storage"`
}
type NFS struct {
	Server string `yaml:"server"`
	Path   string `yaml:"path"`
}
type PVPhase string

const (
	PVAvailable PVPhase = "Available"
	PVBound     PVPhase = "Bound"
	PVFailed    PVPhase = "Failed"
)

type PersistentVolumeClaimAbstract struct {
	PVCname      string
	PVCnamespace string
}
