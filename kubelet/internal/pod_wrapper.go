package internal

import (
	"fmt"
	"minik8s/apiobjects"

	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type PodWrapper struct {
	Pod          apiobjects.Pod
	PodSandboxId string
}

func SandboxStateToPodPhase(state cri.PodSandboxState) apiobjects.PodPhase {
	switch state {
	case cri.PodSandboxState_SANDBOX_READY:
		return apiobjects.PodPhase_POD_RUNNING
	case cri.PodSandboxState_SANDBOX_NOTREADY:
		return apiobjects.PodPhase_POD_PENDING
	default:
		return apiobjects.PodPhase_POD_UNKNOWN
	}
}

// 合并Pod状态，同时输出Pod状态是否有变化
// 最终输出的状态直接写入origin
func MergePodStates(origin *apiobjects.Pod, new *apiobjects.Pod) (is_diff bool, err error) {
	is_diff = false
	if origin == nil || new == nil {
		err = fmt.Errorf("origin or new pod is nil")
		return
	}

	if origin.Status.PodPhase != new.Status.PodPhase {
		origin.Status.PodPhase = new.Status.PodPhase
		is_diff = true
	}

	if origin.Status.PodIP != new.Status.PodIP {
		origin.Status.PodIP = new.Status.PodIP
		is_diff = true
	}

	origin.CreationTimestamp = new.CreationTimestamp

	// Check container status
	for i := range origin.Spec.Containers {
		isContainerDiff, err_1 := MergeContainerStates(&origin.Spec.Containers[i], &new.Spec.Containers[i])
		if err_1 != nil {
			err = err_1
			return
		}
		is_diff = is_diff || isContainerDiff
	}

	return
}

func MergeContainerStates(origin *apiobjects.Container, new *apiobjects.Container) (is_diff bool, err error) {
	is_diff = false
	if origin == nil || new == nil {
		err = fmt.Errorf("origin or new container is nil")
		return
	}

	if origin.Status.State != new.Status.State {
		origin.Status.State = new.Status.State
		is_diff = true
	}

	return
}
