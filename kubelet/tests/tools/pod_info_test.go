package tests

import (
	cri "minik8s/kubelet/internal/cri_proxy"
	"testing"
)

func TestGetPodInfo(t *testing.T) {
	sandboxs, err := cri.ListPods()
	if err != nil {
		t.Error("ListPods error:", err)
	}

	for _, sandbox := range sandboxs {
		pod, err := cri.GetPodStatus(sandbox.Id)
		if err != nil {
			t.Error("GetPodInfo error:", err)
		}
		t.Log("Pod info:", pod)
	}
}
