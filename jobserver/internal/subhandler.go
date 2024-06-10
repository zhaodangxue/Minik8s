package internal

import (
	"encoding/json"
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

func collectJobStatus(job *apiobjects.Job){
	utils.Info("Collecting job status: ", job)
	if job.Status.PodRef.Name == "" || job.Status.PodIp == "" {
		utils.Debug("Ignoring job ", job.ObjectMeta.Name, ": Job pod is not ready")
		return
	}
	joburl := "http://" + job.Status.PodIp + ":8080/status"
	jsonmap := make(map[string]interface{})
	err := utils.GetUnmarshal(joburl, &jsonmap)
	if err != nil {
		utils.Error("Error GetUnmarshal job status: ", err)
		return
	}
	job.Status.JobState = jsonmap["status"].(apiobjects.JobState)
	if job.Status.JobState == apiobjects.JobState_Success || job.Status.JobState == apiobjects.JobState_Failed {
		jsonByte, err := json.Marshal(jsonmap)
		if err != nil {
			utils.Error("Error marshalling job status: ", err)
			return
		}
		job.Status.Output = string(jsonByte)
	}
	return
}

func callJobRun(job *apiobjects.Job, jsonStr string) (responseJson map[string]interface{}, err error) {
	utils.Info("Calling job run: ", job)
	if job.Status.PodRef.Name == "" || job.Status.PodIp == "" {
		utils.Debug("Ignoring run job ", job.ObjectMeta.Name, ": Job pod is not ready")
		return
	}
	joburl := "http://" + job.Status.PodIp + ":8080/run"
	response, err := utils.PostWithString(joburl, jsonStr)
	if err != nil {
		return
	}
	err = utils.ReadUnmarshal(response.Body, &responseJson)
	return
}
