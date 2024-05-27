package apiobjects

type FunctionCtlInput struct {
	Object       `yaml:",inline"`
	FunctionSpec FunctionSpec `yaml:"spec"`
	BuildOptions BuildOptions `yaml:"buildOptions"`
}

type Function struct {
	Object
	Spec   FunctionSpec
	Status FunctionStatus
}

type FunctionSpec struct {
	// 最小副本数
	// 副本数量不会低于这个值
	// ScaleToZero: 如果这个字段为0，则表示副本数可以缩减到0
	// 负数值无效
	MinReplicas int `yaml:"minReplicas"`
	// 最大副本数
	// 副本数量不会高于这个值
	// TODO: 不能小于MinReplicas
	MaxReplicas int `yaml:"maxReplicas"`
	// 预期每个副本的QPS
	// 用于扩容和缩容逻辑，Gateway统计一段时间内Function的总QPS，然后除以副本数，得到每个副本的QPS
	// 如果高于targetQPSPerReplica，则扩容，如果低于targetQPSPerReplica，则缩容
	TargetQPSPerReplica int `yaml:"targetQPSPerReplica"`
}

type FunctionStatus struct {
	// 镜像地址
	// 镜像由ctl构建，然后推送到镜像仓库
	// 这个字段在ctl构建时会被填充
	ImageUrl string
	// 用于引用ReplicaSet，控制副本数
	// 一个Function对应一个ReplicaSet，由Gateway负责管理
	// 包括创建、删除、扩缩容等
	ScaleTargetRef ScaleTargetRef
	// 用于引用Service，提供服务
	// 一个Function对应一个Service，由Gateway负责管理
	ServiceRef ServiceRef
}

type BuildOptions struct {
	// 构建镜像时额外使用的命令
	ExtraCommands []string `yaml:"extraCommands"`
	// 要被拷贝到镜像中函数文件所在的目录
	// 会将这个目录下的所有文件拷贝到镜像中的 /function 目录下, 与watchdog放在一起
	FunctionFileDir string `yaml:"functionFileDir"`
}
