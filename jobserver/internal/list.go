package internal

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	"sync"
)

// 先获取所有的job
// 然后获取所有的pod
// 最后推送所有的job到cluster
func ClusterFullCheck() {
	// 全量更新Job
	jobs, err := getAllJobs()
	if err != nil {
		utils.Error(err)
		return
	}
	newJobsMap := &sync.Map{}
	// CHECK: 需要大锁吗
	for _, job := range jobs {
		oldJob, ok := Jobs().Load(job.GetObjectPath())
		utils.Debug("Job in cluster:", job)
		if ok {
			utils.Debug("Job in cache:", oldJob)
			job.Status = oldJob.(*apiobjects.Job).Status
		}
		newJobsMap.Store(job.GetObjectPath(), job)
	}
	ReplaceJobs(newJobsMap)

	// 全量更新Pod
	pods, err := getAllPods()
	if err != nil {
		utils.Error(err)
		return
	}
	for _, pod := range pods {
		handlePodStateUpdate(pod)
	}

	// 获取所有job的新状态
	jobStatusFullCheck()

	// 推送所有job到cluster
	utils.PutWithJson(route.Prefix+route.JobPath, Jobs())
}

func getAllJobs() (jobs []*apiobjects.Job, err error) {
	err = utils.GetUnmarshal(route.Prefix+route.JobPath, &jobs)
	return
}

func getAllPods() (pods []*apiobjects.Pod, err error) {
	err = utils.GetUnmarshal(route.Prefix+route.PodPath, &pods)
	return
}

func jobStatusFullCheck() {
	jobs := Jobs()
	jobs.Range(func(key, value interface{}) bool {
		job := value.(*apiobjects.Job)
		collectJobStatus(job)
		return true
	})
}
