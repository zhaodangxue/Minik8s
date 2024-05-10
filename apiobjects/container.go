package apiobjects

type Container struct {
	Name         string          `yaml:"name"`
	Image        string          `yaml:"image"`
	Ports        []ContainerPort `yaml:"ports"`
	VolumeMounts []VolumeMount   `yaml:"volumeMounts"`
}

type ContainerPort struct {
	Name          string `yaml:"name"`
	ContainerPort int32  `yaml:"containerPort"`
	HostPort      int32  `yaml:"hostPort"`
}

type VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}
