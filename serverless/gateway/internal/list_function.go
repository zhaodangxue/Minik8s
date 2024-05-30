package internal

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	"sync/atomic"
)

func FunctionHandlerOnList() error {
	fmt.Println("Function Handler On List")
	var vals []string
	vals, _ = etcd.Get_prefix(route.FunctionPath)
	var functions []apiobjects.Function
	for _, val := range vals {
		function := apiobjects.Function{}
		err := json.Unmarshal([]byte(val), &function)
		if err != nil {
			return err
		}
		functions = append(functions, function)
	}
	for _, function := range functions {
		name := function.ObjectMeta.Name
		fmt.Println("Function name: " + name)
		if ServerlessGatewayInstance.functions[name] == nil {
			rs, _ := etcd.Get(route.ReplicasetPath + "/" + "default" + "/" + "function-" + name + "-rs")
			if rs == "" {
				fmt.Println("Replicaset not found with name: " + "function-" + name + "-rs")
				continue
			}
			replicaset := apiobjects.Replicaset{}
			err := json.Unmarshal([]byte(rs), &replicaset)
			if err != nil {
				return err
			}
			ServerlessGatewayInstance.functions[name] = &FunctionWrapper{
				Function:    &function,
				QPSCounter:  &atomic.Int64{},
				ScaleTarget: replicaset.Spec.Replicas,
			}
			fmt.Println("Function: " + name + " added to ServerlessGatewayInstance.functions")
		} else {
			var replicas int
			val, _ := etcd.Get(route.ReplicasetPath + "/" + "default" + "/" + "function-" + name + "-rs")
			if val == "" {
				fmt.Println("Replicaset not found with name: " + "function-" + name + "-rs")
				continue
			}
			replicaset := apiobjects.Replicaset{}
			err := json.Unmarshal([]byte(val), &replicaset)
			if err != nil {
				return err
			}
			replicas = replicaset.Spec.Replicas
			tmp := replicas
			TargetQps := function.Spec.TargetQPSPerReplica
			CurrentQps := ServerlessGatewayInstance.functions[name].QPSCounter.Load()
			if CurrentQps == 0 {
				replicas = 0
			} else {
				if CurrentQps/(int64(TargetQps)*30) > 1 {
					if replicas+1 <= function.Spec.MaxReplicas {
						replicas = replicas + 1
					}
				} else if CurrentQps/(int64(TargetQps)*30) < 1 {
					if replicas-1 >= function.Spec.MinReplicas {
						replicas = replicas - 1
					}
				}
			}
			//重置QPSCounter为0
			ServerlessGatewayInstance.functions[name].QPSCounter.Store(0)
			ServerlessGatewayInstance.functions[name].ScaleTarget = replicas
			if replicas != tmp {
				replicaset.Spec.Replicas = replicas
				url := route.Prefix + route.ReplicasetScale
				_, err = utils.PutWithJson(url, replicaset)
				if err != nil {
					fmt.Println("Failed to scale replicaset")
					return err
				}
				fmt.Println("Function: " + name + " scale ")
			}
		}
	}
	return nil
}
