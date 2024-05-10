package tests

import (
	"minik8s/kubelet/internal"
	"testing"
)

func TestGetPodInfo(t *testing.T) {
	sandboxs, err := internal.ListPods()
	if err != nil {
		t.Error("ListPods error:", err)
	}

	for _, sandbox := range sandboxs {
		pod, err := internal.GetPodInfo(sandbox.Id)
		if err != nil {
			t.Error("GetPodInfo error:", err)
		}
		t.Log("Pod info:", pod)
	}
}
