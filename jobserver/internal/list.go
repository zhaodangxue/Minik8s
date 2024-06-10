package internal

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	"sync"
)

// 先获取所有的job
// 然后获取所有的pod
func FullCheck() {
	// 全量更新Job
	jobs, err := getAllJobs()
	if err != nil {
		utils.Error(err)
		return
	}
	jobsMap := &sync.Map{}
	for _, job := range jobs {
		jobsMap.Store(job.GetObjectPath(), job)
	}
	utils.Debug("Jobs: ", jobsMap)
	ReplaceJobs(jobsMap)

	// 全量更新Pod
	pods, err := getAllPods()
	if err != nil {
		utils.Error(err)
		return
	}
	for _, pod := range pods {
		handlePodStateUpdate(pod)
	}
}

func getAllJobs() (jobs []*apiobjects.Job, err error) {
	err = utils.GetUnmarshal(route.Prefix + route.JobPath, &jobs)
	return
}

func getAllPods() (pods []*apiobjects.Pod, err error) {
	err = utils.GetUnmarshal(route.Prefix + route.PodPath, &pods)
	return
}
