package cri_proxy

import (
	"minik8s/apiobjects"
	"minik8s/global"
	"minik8s/kubelet/internal"
	"minik8s/utils"

	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

type PodStatus cri.PodSandboxStatusResponse

func GetPodStatus(sandboxId string) (response *PodStatus, err error) {

	// Parameters
	ctx := getContext()

	runtimeServiceClient, err := getRuntimeServiceClient()
	if err != nil {
		utils.Error("getRuntimeServiceClient error:", err)
		return
	}

	// Get pod sandbox status
	podStatusRequest := &cri.PodSandboxStatusRequest{PodSandboxId: sandboxId, Verbose: true}
	response_raw, err := runtimeServiceClient.PodSandboxStatus(ctx, podStatusRequest)
	response = (*PodStatus)(response_raw)
	if err != nil {
		utils.Error("ListPodSandbox error:", err)
		return
	}

	utils.Debug("Pod sandbox status:", response)
	return
}

type ContainerStatus cri.ContainerStatusResponse

func GetContainerStatuses(containerId string) (container *ContainerStatus, err error) {

	// Parameters
	ctx := getContext()

	runtimeServiceClient, err := getRuntimeServiceClient()
	if err != nil {
		utils.Error("getRuntimeServiceClient error:", err)
		return
	}

	// Get container status
	containerStatusRequest := &cri.ContainerStatusRequest{ContainerId: containerId}
	response_raw, err := runtimeServiceClient.ContainerStatus(ctx, containerStatusRequest)
	container = (*ContainerStatus)(response_raw)
	if err != nil {
		utils.Error("ContainerStatus error:", err)
		return
	}

	utils.Debug("Container status:", container)
	return
}

func ListPods() (sandboxs []*cri.PodSandbox, err error) {

	// Parameters
	ctx := getContext()

	runtimeServiceClient, err := getRuntimeServiceClient()
	if err != nil {
		utils.Error("getRuntimeServiceClient error:", err)
		return
	}

	// List pod sandboxes
	listRequest := &cri.ListPodSandboxRequest{}
	response, err := runtimeServiceClient.ListPodSandbox(ctx, listRequest)
	if err != nil {
		utils.Error("ListPodSandbox error:", err)
		return
	}
	utils.Debug("Pod sandboxes:", response.Items)
	sandboxs = response.Items

	return
}

func GetAllPods() (pods []*PodStatus, err error) {
	sandboxs, err := ListPods()
	if err != nil {
		utils.Error("ListPods error:", err)
		return
	}

	for _, sandbox := range sandboxs {
		pod, err := GetPodStatus(sandbox.Id)
		if err != nil {
			utils.Error("GetPodInfo error:", err)
			continue
		}
		pods = append(pods, pod)
	}

	return
}

func GetObjectRef(pod *PodStatus) (objectRef *apiobjects.ObjectRef) {
	objectRef = &apiobjects.ObjectRef{
		TypeMeta: apiobjects.TypeMeta{
			ApiVersion: global.ApiVersion,
			Kind:       "Pod",
		},
		Name:      pod.Status.Metadata.Name,
		Namespace: pod.Status.Metadata.Namespace,
		UID:       pod.Status.Metadata.Uid,
	}
	return
}

func UpdatePodStatus(pod *apiobjects.Pod) {
	// 更新Pod状态
	podStatus, err := GetPodStatus(pod.Status.SandboxId)
	if err != nil {
		utils.Error("GetPodInfo error:", err)
		return
	}
	pod.Status.PodPhase = internal.SandboxStateToPodPhase(podStatus.Status.State)
	pod.Status.PodIP = podStatus.Status.Network.Ip

	// 更新Container状态
	for i := range pod.Spec.Containers {
		container := &pod.Spec.Containers[i]
		if container.Status == nil {
			continue
		}
		containerStatus, err := GetContainerStatuses(container.Status.Id)
		if err != nil {
			utils.Error("GetContainerStatus error:", err)
			continue
		}
		container.Status.State = apiobjects.ContainerState(containerStatus.Status.State)
		container.Status.CreatedAt = containerStatus.Status.CreatedAt
		container.Status.StartedAt = containerStatus.Status.StartedAt
		container.Status.FinishedAt = containerStatus.Status.FinishedAt
		container.Status.ExitCode = containerStatus.Status.ExitCode
		container.Status.Reason = containerStatus.Status.Reason
		container.Status.Message = containerStatus.Status.Message
	}
}

func UpdateContainerStatus(container *apiobjects.Container) {
	// 更新Container状态
	ctx := getContext()

	runtimeServiceClient, err := getRuntimeServiceClient()
	if err != nil {
		utils.Error("getRuntimeServiceClient error:", err)
		return
	}

	containerRequest := &cri.ContainerStatusRequest{ContainerId: container.Status.Id}
	response, err := runtimeServiceClient.ContainerStatus(ctx, containerRequest)
	if err != nil {
		utils.Error("ContainerStatus error:", err)
		return
	}

	container.Status.State = apiobjects.ContainerState(response.Status.State)
	container.Status.CreatedAt = response.Status.CreatedAt
	container.Status.StartedAt = response.Status.StartedAt
	container.Status.FinishedAt = response.Status.FinishedAt
	container.Status.ExitCode = response.Status.ExitCode
	container.Status.Reason = response.Status.Reason
	container.Status.Message = response.Status.Message
}
