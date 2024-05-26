package apiobjects

type ScaleTargetRef Object
type HorizontalPodAutoscaler struct {
	Object `yaml:",inline"`
	Spec   HorizontalPodAutoscalerSpec `yaml:"spec"`
	Stat   HorizontalPodAutoscalerStat `yaml:"-"`
}
type HorizontalPodAutoscalerSpec struct {
	MinReplicas    int            `yaml:"minReplicas"`
	MaxReplicas    int            `yaml:"maxReplicas"`
	ScaleTargetRef ScaleTargetRef `yaml:"scaleTargetRef"`
	Metrics        Metrics        `yaml:"metrics"`
	ScaleInterval  int            `yaml:"scaleInterval"`
}
type Metrics struct {
	CPUUtilizationPercentage int     `yaml:"CPUUtilizationPercentage"`
	MemoryUtilizationUsage   float32 `yaml:"MemoryUtilizationUsage"`
}
type HorizontalPodAutoscalerStat struct {
	CurrnentReplicaseCPUUsage int     `yaml:"-"`
	CurrentReplicaseMemUsage  float32 `yaml:"-"`
}
