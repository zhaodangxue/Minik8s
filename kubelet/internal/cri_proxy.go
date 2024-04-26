package internal

import (
	"context"
	"minik8s/apiobjects"

	"minik8s/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cri "k8s.io/cri-api/pkg/apis/runtime/v1"
)

func getContext() context.Context {
	return context.Background()
}

func CreatePod(pod apiobjects.Pod) (PodSandboxID string, err error) {

	// Parameters
	ctx := getContext();

	// Create a gRPC client connection
	conn, err := grpc.Dial("unix:///run/containerd/containerd.sock", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	// Create the runtime service client using the gRPC client connection
	runtimeServiceClient := cri.NewRuntimeServiceClient(conn)
	imageSeviceClient := cri.NewImageServiceClient(conn)


	// Create a pod sandbox
	sandboxConfig := cri.PodSandboxConfig{
		Metadata: &cri.PodSandboxMetadata{
			Name:      pod.ObjectMeta.Name,
			Namespace: pod.ObjectMeta.Namespace,
			Uid:       pod.UID,
		},
	}
	runRequest := &cri.RunPodSandboxRequest{
		Config: &sandboxConfig,
	}
	response, err := runtimeServiceClient.RunPodSandbox(ctx, runRequest)
	PodSandboxID = response.PodSandboxId
	if err != nil {
		utils.Error("RunPodSandbox error:", err)
		return
	}
	utils.Info("Pod sandbox created with ID:", PodSandboxID)

	// Create containers
	for _, container := range pod.Spec.Containers {

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
		}
		createContainerRequest := &cri.CreateContainerRequest{
			PodSandboxId: PodSandboxID,
			Config:       &containerConfig,
			SandboxConfig: &sandboxConfig,
		}
		containerID, err1 := runtimeServiceClient.CreateContainer(ctx, createContainerRequest)
		err = err1
		if err != nil {
			utils.Error("CreateContainer error:", err)
			return
		}
		utils.Info("Container created with ID:", containerID)
	}

	return 
}
