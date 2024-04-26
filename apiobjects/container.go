package apiobjects

type Container struct {
	Name string
	Image string
	Ports []ContainerPort
	VolumeMounts []VolumeMount
}

type ContainerPort struct {
	ContainerPort int32
	HostPort int32
}

type VolumeMount struct {
	Name string
	MountPath string
}
