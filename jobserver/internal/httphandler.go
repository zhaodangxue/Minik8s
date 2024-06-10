package internal

import (
	"io"
	"minik8s/apiobjects"

	"github.com/gin-gonic/gin"
)

func JobRunHandler(c *gin.Context){
	requestBody, err :=io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, gin.H{"error": "Error reading request body"})
		return
	}
	jobName := c.Param("name")
	namespace := c.Param("namespace")
	jobPath := "/api/job/" + namespace + "/" + jobName
	job, ok := Jobs().Load(jobPath)
	if !ok {
		c.JSON(400, gin.H{"error": "Job not found"})
		return
	}
	res, err := callJobRun(job.(*apiobjects.Job), string(requestBody))
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, res)
}

func JobGetStatusHandler(c *gin.Context){
	jobName := c.Param("name")
	namespace := c.Param("namespace")
	jobPath := "/api/job/" + namespace + "/" + jobName
	if job, ok := Jobs().Load(jobPath); ok {
		collectJobStatus(job.(*apiobjects.Job))
		c.JSON(200, job)
		return
	}
	c.JSON(400, gin.H{"error": "Job not found"})
}
