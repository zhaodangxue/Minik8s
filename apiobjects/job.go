package apiobjects

/*
Job 任务

提供异步任务的抽象

要求用户提供一个python脚本，脚本名字为`job.py`

job.py必须包含两个函数
- def run(): 任务的入口函数
- def get_status(): 任务的状态查询函数

基础环境中的`watchdog.py`会调用这两个函数
  - run()函数会被`watchdog.py`调用，用于开始任务的执行。
    run()会在一个子线程中执行，不会阻塞`watchdog.py`的执行，用户可以任意实现run()函数为阻塞或非阻塞。
  - get_status()函数会被`watchdog.py`调用，用于查询任务的状态。返回值为一个json字符串，必须包含一个status字段，表示任务的状态。

status字段的值必须是以下值之一
- "pending": 任务等待执行
- "running": 任务正在运行
- "success": 任务成功完成
- "failed": 任务失败
*/
type Job struct {
	Object
	Spec   JobSpec   `yaml:"spec"`
	Status JobStatus `yaml:"-"`
}

type JobSpec struct {
	// 构建任务时额外使用的命令
	BuildOptions BuildOptions `yaml:"buildOptions"`
}

type JobStatus struct {
	// 任务对应的镜像url
	ImageUrl string `yaml:"-"`
	// 任务的状态
	JobState JobState
	// 任务的输出(json字符串)
	Output string
	// 任务所在的pod的引用
	PodRef PodRef
	// Pod的状态(cache，注意一致性)
	PodIp string
}

type JobState string

const (
	// pending: 任务等待执行，创建job后的初始状态
	JobState_Pending JobState = "pending"
	JobState_Running JobState = "running"
	JobState_Success JobState = "success"
	JobState_Failed  JobState = "failed"
)
