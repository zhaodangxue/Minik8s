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

// PodPhase是Pod的状态。
const (
	// Created意味着Pod已经在系统中被创建，但是还没有被所调度到的Node启动。
	PodCreated PodPhase = "Created"
	// PodPending 意味着Pod已经被调度到了一个Node上，并在Node上被创建，但尚处于NotReady状态。
	// 亦即还有一些条件没有满足，Pod中的容器还尚未全部被启动。
	PodPending PodPhase = "Pending"
	// PodRunning 对应Sandbox的Ready状态，意味着Pod中的容器已经被启动。
	PodRunning PodPhase = "Running"
	// PodSucceeded 意味着Pod中的所有容器都已经成功地终止，并且不会再重启。
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed 意味着Pod中的所有容器都已经终止，并且至少有一个容器是非正常终止的。
	PodFailed PodPhase = "Failed"
	// PodUnknown 意味着Pod的状态无法被获取。
	// Kubelet获取Pod的状态失败时，会将Pod的状态设置为Unknown。
	PodUnknown PodPhase = "Unknown"
)
