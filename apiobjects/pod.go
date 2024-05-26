package apiobjects

type Pod struct {
	Object `yaml:",inline"`
	Spec   PodSpec  `yaml:"spec"`
	Status PodState `yaml:"-"`
	Stats  PodStats `yaml:"-"`
}

type PodSpec struct {
	Containers   []Container       `yaml:"containers"`
	Volumes      []Volume          `yaml:"volumes"`
	NodeSelector map[string]string `yaml:"nodeSelector"`
}

func (c *Pod) AddLabel(key, value string) {
	if c.ObjectMeta.Labels == nil {
		c.ObjectMeta.Labels = make(map[string]string)
	}
	c.ObjectMeta.Labels[key] = value
}

type PodState struct {
	SandboxId string
	PodPhase  PodPhase
	// TODO: conditions
	HostIP string
	PodIP  string
}

type PodStats struct {
	CpuUsage    CpuUsage
	MemoryUsage MemoryUsage
}

type CpuUsage struct {
	// 每秒内在所有核上一共占用了CPU多少纳秒
	TotalNanos uint64
}

type MemoryUsage struct {
	// 工作集大小
	WorkingSetBytes uint64
	// 可用内存大小，即可以被Pod使用的内存大小
	AvailableBytes uint64
	// 实际使用的内存大小
	UsageBytes uint64
}

func (c *CpuUsage) GetCpuUsage() float32 {
	if c.TotalNanos == 0 {
		return 0
	}
	ret := float32(c.TotalNanos) / 1e9
	return ret
}
func (m *MemoryUsage) GetMemUsage() float32 {
	ret := float32(m.UsageBytes) / float32(1024) / float32(1024) / float32(1024)
	return ret
}

type PodPhase string

// PodPhase是Pod的状态。
const (
	// Created意味着Pod已经在系统中被创建，但是还没有被所调度到的Node启动。
	PodPhase_POD_CREATED PodPhase = "Created"
	// PodPhase_POD_PENDING 意味着Pod已经被调度到了一个Node上，并在Node上被创建，但尚处于NotReady状态。
	// 亦即还有一些条件没有满足，Pod中的容器还尚未全部被启动。
	PodPhase_POD_PENDING PodPhase = "Pending"
	// PodPhase_POD_RUNNING 对应Sandbox的Ready状态，意味着Pod中的容器已经被启动。
	PodPhase_POD_RUNNING PodPhase = "Running"
	// PodPhase_POD_SUCCEEDED 意味着Pod中的所有容器都已经成功地终止，并且不会再重启。
	PodPhase_POD_SUCCEEDED PodPhase = "Succeeded"
	// PodPhase_POD_FAILED 意味着Pod中的所有容器都已经终止，并且至少有一个容器是非正常终止的。
	PodPhase_POD_FAILED PodPhase = "Failed"
	// PodPhase_POD_UNKNOWN 意味着Pod的状态无法被获取。
	// Kubelet获取Pod的状态失败时，会将Pod的状态设置为Unknown。
	PodPhase_POD_UNKNOWN PodPhase = "Unknown"
)

type PodTemplate struct {
	Metadata ObjectMeta `yaml:"metadata"`
	Spec     PodSpec    `yaml:"spec"`
}

func ToPod(podTemplate *PodTemplate) *Pod {
	var pod Pod
	pod.TypeMeta.ApiVersion = "v1"
	pod.TypeMeta.Kind = "Pod"
	pod.ObjectMeta = podTemplate.Metadata
	pod.Spec = podTemplate.Spec
	return &pod
}
