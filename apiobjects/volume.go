package apiobjects

type Volume struct {
	Name         string `yaml:"name"`
	VolumeSource `yaml:",inline"`
}
type VolumeSource struct {
	EmptyDir              *EmptyDirVolumeSource              `yaml:"emptyDir"`
	HostPath              *HostPathVolumeSource              `yaml:"hostPath"`
	NFS                   *NFSVolumeSource                   `yaml:"nfs"`
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `yaml:"persistentVolumeClaim"`
}
type EmptyDirVolumeSource struct {
}
type HostPathVolumeSource struct {
	Path string `yaml:"path"`
}
type NFSVolumeSource struct {
	Server      string `yaml:"server"`
	Path        string `yaml:"path"`
	BindingPath string `yaml:"-"`
}
type PersistentVolumeClaimVolumeSource struct {
	ClaimName      string `yaml:"claimName"`
	ClaimNamespace string `yaml:"claimNamespace"`
}
