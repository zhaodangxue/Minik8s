package job

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func JobCreateHandler(c *gin.Context) {
	var job apiobjects.Job
	var action = apiobjects.Create
	err := utils.ReadUnmarshal(c.Request.Body, &job)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	if job.Name == "" {
		c.String(200, "Please input the job name")
		return
	}
	url := job.GetObjectPath()
	val, _ := etcd.Get(url)
	if val != "" {
		action = apiobjects.Update
		utils.Info("JobCreateHandler: job already exists, replacing")
	}
	jobJson, _ := json.Marshal(job)
	err = etcd.Put(url, string(jobJson))
	if err != nil {
		c.String(500, "Create job failed")
		return
	}
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = action
	topicMessage.Object = string(jobJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.JobTopic(), string(topicMessageJson))
	c.String(200, "the event is created")
}

func JobGetAllHandler(c *gin.Context) {
	url := route.JobPath
	jobJsons, err := etcd.Get_prefix(url)
	if err != nil {
		c.String(500, "Get jobs failed")
		return
	}
	jobs := []*apiobjects.Job{}
	for _, jobJson := range jobJsons {
		job := &apiobjects.Job{}
		err := json.Unmarshal([]byte(jobJson), job)
		if err != nil {
			utils.Error("JobGetAllHandler: ", err)
			continue
		}
		jobs = append(jobs, job)
	}
	c.JSON(200, jobs)
}

func JobUpdateAllHandler(c *gin.Context) {
	jobs := []*apiobjects.Job{}
	utils.ReadUnmarshal(c.Request.Body, &jobs)
	utils.Debug("JobUpdateAllHandler: get jobs:", jobs)
	for _, job := range jobs {
		jobJson, _ := json.Marshal(job)
		err := etcd.Put(job.GetObjectPath(), string(jobJson))
		if err != nil {
			utils.Error("JobUpdateAllHandler: ", err)
			c.String(500, "Update job failed")
			return
		}
		utils.Info("JobUpdateAllHandler: job ", job.Name, " is updated")
	}
	c.String(200, "the job is updated")
}
