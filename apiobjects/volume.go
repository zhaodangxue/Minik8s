package apiobjects

type Volume struct {
	Name string `yaml:"name"`
}
type VolumeSource struct {
	EmptyDir *EmptyDirVolumeSource `yaml:"emptyDir"`
	HostPath *HostPathVolumeSource `yaml:"hostPath"`
}
type EmptyDirVolumeSource struct {
}
type HostPathVolumeSource struct {
	Path string `yaml:"path"`
}
