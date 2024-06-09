package internal

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/utils"

	"github.com/redis/go-redis/v9"
)

func unmarshallTopicMessage(message *redis.Message, object interface{}, action *apiobjects.ActionType) (err error) {
	messageData := message.Payload
	topicMessage := apiobjects.TopicMessage{}
	err = json.Unmarshal([]byte(messageData), &topicMessage)
	if err != nil {
		return
	}
	*action = topicMessage.ActionType
	err = json.Unmarshal([]byte(topicMessage.Object), object)
	return
}

func JobMessageHandler(message *redis.Message) {
	var action apiobjects.ActionType
	var job apiobjects.Job
	err := unmarshallTopicMessage(message, &job, &action)
	if err != nil {
		utils.Error("Error unmarshalling job message: ", err)
		return
	}
	switch action {
	case apiobjects.Create:
		handleJobCreate(&job)
	default:
		utils.Warn("Unknown action type: ", action)
	}
}

func PodStateMessageHandler(message *redis.Message) {
	var action apiobjects.ActionType
	var pod apiobjects.Pod
	err := unmarshallTopicMessage(message, &pod, &action)
	if err != nil {
		utils.Error("Error unmarshalling pod state message: ", err)
		return
	}
	switch action {
	case apiobjects.Update:
		jobPath, ok := pod.Labels["job"]
		if !ok {
			utils.Debug("Ignoring Pod ", pod.ObjectMeta.Name, ": Pod has no job label: job")
			return
		}
		job, ok := Jobs().Load(jobPath)
		if !ok {
			utils.Debug("Ignoring Pod ", pod.ObjectMeta.Name, ": Job not found: ", jobPath)
			return
		}
		handlePodStateUpdate(job.(*apiobjects.Job), &pod)
	default:
		utils.Debug("Not supported podState action type", action)
	}
}
