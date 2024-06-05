package cri_proxy

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/utils"
	"os/exec"

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

	const CpuPeriod = 1000000

	linux := cri.LinuxPodSandboxConfig{}
	linux.CgroupParent = "/kubelet/pod-" + pod.ObjectMeta.UID
	linux.Resources = &cri.LinuxContainerResources{
		CpuPeriod:  CpuPeriod,
		CpuQuota:   int64(float32(CpuPeriod) * pod.Spec.CpuLimit),
		MemoryLimitInBytes: pod.Spec.MemLimit,
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
		Linux:        &linux,
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
			Args:       nil,
			WorkingDir: "/root",
			Envs:       nil,
			Labels:     container.Labels,
			Mounts:     nil,
			Devices:    nil,
		}

		// VolumeMounts
		volumemounts := container.VolumeMounts
		if len(volumemounts) != 0 {
			volumes := pod.Spec.Volumes
			containerConfig.Mounts = MountContainer(volumes, volumemounts, pod.ObjectMeta.UID)
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
	utils.Info("DeletePod with ID:", podSandboxId)
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
func MountContainer(volumes []apiobjects.Volume, volumeMounts []apiobjects.VolumeMount, uuid string) (mounts []*cri.Mount) {
	//TODO
	var mountPath string
	var mountName string
	for _, volumeMount := range volumeMounts {
		mountPath = volumeMount.MountPath
		mountName = volumeMount.Name
		for _, volume := range volumes {
			if volume.Name == mountName {
				if volume.HostPath != nil {
					mounts = append(mounts, &cri.Mount{
						ContainerPath: mountPath,
						HostPath:      volume.HostPath.Path,
					})
				} else if volume.NFS != nil {
					server_ip := volume.NFS.Server
					server_path := volume.NFS.Path
					local_path := NFSMountLocal(server_ip, server_path, uuid, mountName)
					mounts = append(mounts, &cri.Mount{
						ContainerPath: mountPath,
						HostPath:      local_path,
					})
				} else if volume.PersistentVolumeClaim != nil {
					pvcName := volume.PersistentVolumeClaim.ClaimName
					pvcNamespace := volume.PersistentVolumeClaim.ClaimNamespace
					local_path := PVMountLocal(pvcName, pvcNamespace, uuid)
					mounts = append(mounts, &cri.Mount{
						ContainerPath: mountPath,
						HostPath:      local_path,
					})
				}
			}
		}
	}
	return
}
func NFSMountLocal(server_ip string, server_path string, uuid string, DirName string) (local_path string) {
	//TODO
	url := server_ip + ":" + server_path
	url_local := global.WorkerMountDir + "/" + uuid + "/" + DirName
	cmd_mkdir := exec.Command("mkdir", "-p", url_local)
	err_mkdir := cmd_mkdir.Run()
	if err_mkdir != nil {
		utils.Error("NFSMountLocal error:", err_mkdir)
		local_path = ""
		return
	}
	cmd := exec.Command("sudo", "mount", "-t", "nfs", url, url_local)
	//打印出执行的命令
	utils.Info("mount -t nfs", url, url_local)
	err := cmd.Run()
	if err != nil {
		utils.Error("NFSMountLocal error:", err)
		local_path = ""
		return
	}
	utils.Info("NFS mounted with url:", url)
	local_path = url_local
	return
}
func PVMountLocal(pvcName string, pvcNamespace string, uuid string) (local_path string) {
	url := route.Prefix + route.PVCPath + "/" + pvcNamespace + "/" + pvcName
	var pvc apiobjects.PersistentVolumeClaim
	err := utils.GetUnmarshal(url, &pvc)
	if err != nil {
		utils.Error("PVMountLocal error:", err)
		local_path = ""
		return
	}
	server_ip := global.Nfsserver
	server_path := pvc.PVBinding.PVpath
	DirName := pvc.Namespace + "-" + pvc.Name
	local_path = NFSMountLocal(server_ip, server_path, uuid, DirName)
	return
}
