package cri_proxy

import (
	"minik8s/apiobjects"
	"minik8s/utils"

	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func convertProtocol(protocol apiobjects.ContainerPortProtocol) (criProtocol cri.Protocol) {
	switch protocol {
	case apiobjects.Protocol_TCP:
		criProtocol = cri.Protocol_TCP
	case apiobjects.Protocol_UDP:
		criProtocol = cri.Protocol_UDP
	}
	return
}

func getSandboxConfig(pod *apiobjects.Pod) (sandboxConfig cri.PodSandboxConfig) {

	portMappings := make([]*cri.PortMapping, 0)
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			portMappings = append(portMappings, &cri.PortMapping{
				ContainerPort: port.ContainerPort,
				HostPort:      port.HostPort,
				Protocol:      convertProtocol(port.Protocol),
			})
		}
	}

	sandboxConfig = cri.PodSandboxConfig{
		Metadata: &cri.PodSandboxMetadata{
			Name:      pod.ObjectMeta.Name,
			Namespace: pod.ObjectMeta.Namespace,
			Uid:       pod.ObjectMeta.UID,
		},
		Hostname:     "",
		Labels:       pod.ObjectMeta.Labels,
		Annotations:  make(map[string]string),
		Linux:        &cri.LinuxPodSandboxConfig{},
		Windows:      nil,
		PortMappings: portMappings,
	}
	return
}

// 创建Pod
//
// 创建的同时会写入 Pod.Status.SandboxId 和 Container.Status.Id 字段
func CreatePod(pod *apiobjects.Pod) (err error) {

	// Parameters
	ctx := getContext()

	runtimeServiceClient, err := getRuntimeServiceClient()
	if err != nil {
		utils.Error("getRuntimeServiceClient error:", err)
		return
	}
	imageSeviceClient, err := getImageServiceClient()
	if err != nil {
		utils.Error("getImageServiceClient error:", err)
		return
	}

	// Create a pod sandbox
	sandboxConfig := getSandboxConfig(pod)
	runRequest := &cri.RunPodSandboxRequest{
		Config: &sandboxConfig,
	}
	response, err := runtimeServiceClient.RunPodSandbox(ctx, runRequest)
	//PodSandboxID = response.PodSandboxId
	if err != nil {
		utils.Error("RunPodSandbox error:", err)
		return
	}
	pod.Status.SandboxId = response.PodSandboxId
	utils.Info("Pod sandbox created with ID:", response.PodSandboxId)

	// Create containers
	for i, container := range pod.Spec.Containers {

		// Pull Image
		imageSpec := &cri.ImageSpec{
			Image: container.Image,
		}
		pullImageRequest := &cri.PullImageRequest{
			Image: imageSpec,
		}
		_, err = imageSeviceClient.PullImage(ctx, pullImageRequest)
		if err != nil {
			utils.Error("PullImage error:", err)
			return
		}
		utils.Info("Image pulled:", container.Image)

		// Create Container
		containerConfig := cri.ContainerConfig{
			Metadata: &cri.ContainerMetadata{
				Name: container.Name,
			},
			Image: &cri.ImageSpec{
				Image: container.Image,
			},
			Command:    []string{"/bin/sh", "-c", "sleep 1000"},
			Args:       nil,
			WorkingDir: "/root",
			Envs:       nil,
			Labels:     container.Labels,
			Mounts:     nil,
			Devices:    nil,
		}

		sandboxConfig.Metadata.Attempt = 1

		createContainerRequest := &cri.CreateContainerRequest{
			PodSandboxId:  pod.Status.SandboxId,
			Config:        &containerConfig,
			SandboxConfig: &sandboxConfig,
		}
		createContainerResponse, err1 := runtimeServiceClient.CreateContainer(ctx, createContainerRequest)

		err = err1
		if err != nil {
			utils.Error("CreateContainer error:", err)
			// TODO: fault tolerance
			return
		}
		if pod.Spec.Containers[i].Status == nil {
			pod.Spec.Containers[i].Status = &apiobjects.ContainerStatus{}
		}
		pod.Spec.Containers[i].Status.Id = createContainerResponse.ContainerId
		utils.Info("Container created with ID:", createContainerResponse)

		_, err = runtimeServiceClient.StartContainer(ctx, &cri.StartContainerRequest{ContainerId: createContainerResponse.ContainerId})
		if err != nil {
			utils.Error("StartContainer error:", err)
			// TODO: fault tolerance
			return
		}
		utils.Info("Container started with ID:", createContainerResponse.ContainerId)
	}

	return
}

func DeletePod(podSandboxId string) (err error) {
	ctx := getContext()

	runtimeServiceClient, err := getRuntimeServiceClient()
	if err != nil {
		utils.Error("getRuntimeServiceClient error:", err)
		return
	}

	stopRequest := &cri.StopPodSandboxRequest{PodSandboxId: podSandboxId}
	_, err = runtimeServiceClient.StopPodSandbox(ctx, stopRequest)
	if err != nil {
		utils.Error("StopPodSandbox error:", err)
		return
	}
	utils.Info("Pod sandbox stopped with ID:", podSandboxId)

	removeRequest := &cri.RemovePodSandboxRequest{PodSandboxId: podSandboxId}
	_, err = runtimeServiceClient.RemovePodSandbox(ctx, removeRequest)
	if err != nil {
		utils.Error("RemovePodSandbox error:", err)
		return
	}
	utils.Info("Pod sandbox removed with ID:", podSandboxId)
	return
}
