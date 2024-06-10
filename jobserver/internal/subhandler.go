package internal

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
)

func handleJobCreate(job *apiobjects.Job){
	utils.Info("Handling job create: ", job)
	// 创建对应的pod
	pod := apiobjects.GetPodFromJob(job)
	utils.Info("Creating pod: ", pod)
	_, err := utils.PostWithJson(route.Prefix + route.PodPath, pod)
	if err != nil {
		utils.Error("Error creating pod: ", err)
		return
	}
	// 缓存job
	Jobs().Store(job.GetObjectPath(), job)
}

func handlePodStateUpdate(pod *apiobjects.Pod){
	jobPath, ok := pod.Labels["job"]
	if !ok {
		utils.Debug("Ignoring Pod ", pod.ObjectMeta.Name, ": Pod has no job label: job")
		return
	}
	iJob, ok := Jobs().Load(jobPath)
	if !ok {
		utils.Warn("Ignoring Pod ", pod.ObjectMeta.Name, ": Job not found: ", jobPath)
		return
	}
	utils.Info("Handling pod state update: ", pod)
	job := iJob.(*apiobjects.Job)
	job.Status.PodIp = pod.Status.PodIP
	job.Status.PodRef = pod.GetObjectRef()
}
