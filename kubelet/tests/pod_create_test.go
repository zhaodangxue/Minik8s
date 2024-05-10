package tests

import (
	"minik8s/apiobjects"
	"minik8s/kubelet/internal"
	"testing"

	"github.com/google/uuid"
)

func TestPodCreate(t *testing.T) {
	pod := apiobjects.Pod{
		Object: apiobjects.Object{
			TypeMeta: apiobjects.TypeMeta{
				ApiVersion: "v1",
				Kind:       "Pod",
			},
			ObjectMeta: apiobjects.ObjectMeta{
				Name:      "test-pod",
				Namespace: "default",
				Labels: map[string]string{
					"app": "test",
				},
				UID: uuid.New().String(),
			},
		},
		Spec: apiobjects.PodSpec{
			Containers: []apiobjects.Container{
				{
					Name:  "test-container",
					Image: "docker.io/library/busybox:latest",
				},
			},
		},
		Status: apiobjects.PodState{},
	}

	internal.CreatePod(pod)

}
