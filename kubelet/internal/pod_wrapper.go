package internal

import (
	"minik8s/apiobjects"

	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

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
