package internal

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"

	"github.com/redis/go-redis/v9"
)

func FunctionHandlerOnWatch(msg *redis.Message) {
	utils.Info("Function Handler On Watch")
	topicMessage := &apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		//调用HandleCreate
		function := &apiobjects.Function{}
		err2 := json.Unmarshal([]byte(topicMessage.Object), function)
		if err2 != nil {
			fmt.Println(err2)
		}
		functionJson, _ := json.Marshal(function)
		HandleCreate([]byte(functionJson))
	case apiobjects.Update:
		//调用HandleUpdate
		function := &apiobjects.Function{}
		err2 := json.Unmarshal([]byte(topicMessage.Object), function)
		if err2 != nil {
			fmt.Println(err2)
		}
		functionJson, _ := json.Marshal(function)
		HandleUpdate([]byte(functionJson))
	case apiobjects.Delete:
		function := &apiobjects.Function{}
		err2 := json.Unmarshal([]byte(topicMessage.Object), function)
		if err2 != nil {
			fmt.Println(err2)
		}
		functionJson, _ := json.Marshal(function)
		HandleDelete([]byte(functionJson))
	}
}

func HandleCreate(data []byte) {
	//utils.Info("Handle Function Create")
	function := &apiobjects.Function{}
	err := json.Unmarshal(data, function)
	if err != nil {
		fmt.Println(err)
	}
	utils.Info("Create function: ", function)

	labels := map[string]string{
		"app": "function-" + function.ObjectMeta.Name + "-label",
	}

	// 1. create the replicaset
	replicaset := apiobjects.Replicaset{
		Object: apiobjects.Object{
			TypeMeta: apiobjects.TypeMeta{
				ApiVersion: "v1",
				Kind:       "ReplicaSet",
			},
			ObjectMeta: apiobjects.ObjectMeta{
				Name:      "function-" + function.ObjectMeta.Name + "-rs",
				Namespace: function.ObjectMeta.Namespace,
			},
		},
		Spec: apiobjects.ReplicasetSpec{
			Replicas: function.Spec.MinReplicas,
			Selector: apiobjects.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiobjects.PodTemplate{
				Metadata: apiobjects.ObjectMeta{
					Labels: labels,
				},
				Spec: apiobjects.PodSpec{
					Containers: []apiobjects.Container{
						{
							Name: "function-" + function.ObjectMeta.Name + "-container",
							//serverIp+":5000/"+imageName+":latest"
							Image: "192.168.1.15:5000/" + function.Status.ImageUrl + ":latest",
							Ports: []apiobjects.ContainerPort{
								{
									Name:          "function-port",
									ContainerPort: 8080,
									HostPort:      8080,
								},
							},
						},
					},
				},
			},
		},
	}
	url := route.Prefix + route.ReplicasetPath
	utils.ApplyApiObject(url, replicaset)

	// 2. create the service
	service := apiobjects.Service{
		APIVersion: "v1",
		Data: apiobjects.MetaData{
			Name:      "function-" + function.ObjectMeta.Name + "-service",
			Namespace: function.ObjectMeta.Namespace,
		},
		Spec: apiobjects.ServiceSpec{
			Selector: labels,
			Type:     apiobjects.ServiceTypeClusterIP,
			Ports: []apiobjects.ServicePort{
				{
					Name:       "function-port",
					Protocol:   apiobjects.ProtocolTCP,
					Port:       8080,
					TargetPort: "function-port",
				},
			},
		},
	}
	url = route.Prefix + route.ServiceApplyPath
	utils.ApplyApiObject(url, service)
}

func HandleDelete(data []byte) {
	//utils.Info("Handle Function Delete")
	function := &apiobjects.Function{}
	err := json.Unmarshal(data, function)
	if err != nil {
		fmt.Println(err)
	}
	utils.Info("Delete function: ", function)

	// 1. delete the replicaset
	url := route.Prefix + route.ReplicasetPath + "/" + function.ObjectMeta.Namespace + "/function-" + function.ObjectMeta.Name + "-rs"
	_, err = utils.Delete(url)
	if err != nil {
		fmt.Println(err)
	}

	// 2. delete the service
	url = route.Prefix + "/api/service/cmd/delete/" + function.ObjectMeta.Namespace + "/function-" + function.ObjectMeta.Name + "-service"
	_, err = utils.Delete(url)
	if err != nil {
		fmt.Println(err)
	}

	// 3. delete the function
	url = route.Prefix + route.FunctionPath + "/" + function.ObjectMeta.Namespace + "/" + function.ObjectMeta.Name
	utils.Delete(url)
}

func HandleUpdate(data []byte) {
	utils.Info("Handle Function Update")
	//utils.Error("Handle Function Update Don't Support")

	function := &apiobjects.Function{}
	err := json.Unmarshal(data, function)
	if err != nil {
		fmt.Println(err)
	}
	utils.Info("Update function: ", function)
	labels := map[string]string{
		"app": "function-" + function.ObjectMeta.Name + "-label",
	}

	// 1. update the replicaset
	replicaset := apiobjects.Replicaset{
		Object: apiobjects.Object{
			TypeMeta: apiobjects.TypeMeta{
				ApiVersion: "v1",
				Kind:       "ReplicaSet",
			},
			ObjectMeta: apiobjects.ObjectMeta{
				Name:      "function-" + function.ObjectMeta.Name + "-rs",
				Namespace: function.ObjectMeta.Namespace,
			},
		},
		Spec: apiobjects.ReplicasetSpec{
			Replicas: function.Spec.MinReplicas,
			Selector: apiobjects.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiobjects.PodTemplate{
				Metadata: apiobjects.ObjectMeta{
					Labels: labels,
				},
				Spec: apiobjects.PodSpec{
					Containers: []apiobjects.Container{
						{
							Name: "function-" + function.ObjectMeta.Name + "-container",
							//serverIp+":5000/"+imageName+":latest"
							Image: "192.168.1.15:5000/" + function.Status.ImageUrl + ":latest",
							Ports: []apiobjects.ContainerPort{
								{
									Name:          "function-port",
									ContainerPort: 8080,
									HostPort:      8080,
								},
							},
						},
					},
				},
			},
		},
	}
	url := route.Prefix + route.ReplicasetPath
	utils.ApplyApiObject(url, replicaset)

	// // 2. create the service
	// service := apiobjects.Service{
	// 	APIVersion: "v1",
	// 	Data: apiobjects.MetaData{
	// 		Name:      "function-" + function.ObjectMeta.Name + "-service",
	// 		Namespace: function.ObjectMeta.Namespace,
	// 	},
	// 	Spec: apiobjects.ServiceSpec{
	// 		Selector: labels,
	// 		Type:     apiobjects.ServiceTypeClusterIP,
	// 		Ports: []apiobjects.ServicePort{
	// 			{
	// 				Name:       "function-port",
	// 				Protocol:   apiobjects.ProtocolTCP,
	// 				Port:       8080,
	// 				TargetPort: "function-port",
	// 			},
	// 		},
	// 	},
	// }
	// url = route.Prefix + route.ServiceApplyPath
	// utils.ApplyApiObject(url, service)

}
