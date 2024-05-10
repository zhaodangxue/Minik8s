package apiobjects

type Container struct {
	Name         string            `yaml:"name"`
	Image        string            `yaml:"image"`
	Ports        []ContainerPort   `yaml:"ports"`
	VolumeMounts []VolumeMount     `yaml:"volumeMounts"`
	Labels       map[string]string `yaml:"labels"`
	Status       ContainerStatus   `yaml:"-"`
}

type ContainerState int32

const (
	ContainerState_CONTAINER_CREATED ContainerState = 0
	ContainerState_CONTAINER_RUNNING ContainerState = 1
	ContainerState_CONTAINER_EXITED  ContainerState = 2
	ContainerState_CONTAINER_UNKNOWN ContainerState = 3
)

type ContainerStatus struct {
	State      ContainerState
	CreatedAt  int64
	StartedAt  int64
	FinishedAt int64
	ExitCode   int32
	Reason     string
	Message    string
}

type ContainerPort struct {
	ContainerPort int32 `yaml:"containerPort"`
	HostPort      int32 `yaml:"hostPort"`
}

type VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}
