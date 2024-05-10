package apiobjects

type PersistentVolume struct {
	TypeMeta   `yaml:",inline"`
	ObjectMeta `yaml:"metadata"`
	Spec       PersistentVolumeSpec `yaml:"spec"`
}
type PersistentVolumeSpec struct {
	Capacity                      capacity `yaml:"capacity"`
	VolumeMode                    string   `yaml:"volumeMode"`
	AccessModes                   []string `yaml:"accessModes"`
	PersistentVolumeReclaimPolicy string   `yaml:"persistentVolumeReclaimPolicy"`
	StorageClassName              string   `yaml:"storageClassName"`
	MountOptions                  []string `yaml:"mountOptions"`
	NFS                           NFS      `yaml:"nfs"`
}
type capacity struct {
	Storage string `yaml:"storage"`
}
type NFS struct {
	Server string `yaml:"server"`
	Path   string `yaml:"path"`
}
