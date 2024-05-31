package tests

import (
	"minik8s/apiobjects"
	cri "minik8s/kubelet/internal/cri_proxy"
	"minik8s/utils"
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
					Name:  "viewer",
					Image: "strm/helloworld-http",
					Ports: []apiobjects.ContainerPort{
						{
							ContainerPort: 80,
							HostPort:      9091,
						},
					},
				},
			},
			MemLimit: 200 * 1024 * 1024, // 200MB
			CpuLimit: 0.5,	// 50%
		},
		Status: apiobjects.PodState{},
	}

	err := cri.CreatePod(&pod)

	if err != nil {
		t.Error(err)
	}

}
func TestPodCreateWithVolume(t *testing.T) {
	var volume apiobjects.Volume
	volume.Name = "test-volume"
	var hostPath apiobjects.HostPathVolumeSource
	volume.HostPath = &hostPath
	volume.HostPath.Path = "/home/zbm/k8s"
	var volumes []apiobjects.Volume
	volumes = append(volumes, volume)
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
					VolumeMounts: []apiobjects.VolumeMount{
						{
							Name:      "test-volume",
							MountPath: "/home",
						},
					},
				},
			},
			Volumes: volumes,
		},
		Status: apiobjects.PodState{},
	}
	cri.CreatePod(&pod)
}
func TestNFSMount(t *testing.T) {
	server_ip := "192.168.20.128"
	server_path := "/home/zbm/nfs"
	DirName := "test-nfs"
	uuid := utils.NewUUID()
	cri.NFSMountLocal(server_ip, server_path, uuid, DirName)
}
