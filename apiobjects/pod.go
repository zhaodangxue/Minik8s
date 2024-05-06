package apiobjects

type Pod struct {
	Object `yaml:",inline"`
	Spec   PodSpec  `yaml:"spec"`
	Status PodState `yaml:"-"`
}

type PodSpec struct {
	Containers   []Container       `yaml:"containers"`
	Volumes      []Volume          `yaml:"volumes"`
	NodeSelector map[string]string `yaml:"nodeSelector"`
}

type PodState struct {
	PodPhase PodPhase
	// TODO: conditions
	HostIP string
	PodIP  string
}

type PodPhase string

// These are the valid statuses of pods.
const (
	// Created意味着Pod已经在系统中被创建，但是还没有被所调度到的Node启动。
	PodCreated PodPhase = "Created"
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
)
