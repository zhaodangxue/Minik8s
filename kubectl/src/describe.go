package command

import (
	"fmt"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"

	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:                        "describe",
	Short:                      "describe a resource",
	Long:                       `Describe a resource by given source type and name.`,
	Run:                        RunDescribe,
	Example:                    "kubectl describe pod xxx -n default",
	SuggestionsMinimumDistance: 1,
	SuggestFor:                 []string{"desc", "d"},
	Args:                       cobra.MinimumNArgs(2),
}

func GetDetailedHelp(name string, np string) (string, error) {
	url := route.Prefix + route.PodPath + "/" + np + "/" + name
	var binding apiobjects.NodePodBinding
	err := utils.GetUnmarshal(url, &binding)
	if err != nil {
		return "", err
	}
	pod := binding.Pod
	node := binding.Node
	var labels string
	for k, v := range pod.ObjectMeta.Labels {
		labels += k + "=" + v + "\n"
	}
	s1 := fmt.Sprintf("Name: %s\nNamespace: %s\n NodeName: %s\n NodeIP: %s\n Labels: %s\n PodIP: %s\n Status: %s\n UUID: %s\n", pod.Name, pod.Namespace, node.ObjectMeta.Name, node.Info.Ip, labels, pod.Status.PodIP, pod.Status.PodPhase, pod.ObjectMeta.UID)
	var s2 string
	s2 = "Containers:\n"
	for _, container := range pod.Spec.Containers {
		var ports string
		var volumnMounts string
		for _, port := range container.Ports {
			ports += fmt.Sprintf("%d/%s\n", port.ContainerPort, "TCP")
		}
		for _, volumnMount := range container.VolumeMounts {
			volumnMounts += fmt.Sprintf("%s/%s\n", volumnMount.Name, volumnMount.MountPath)
		}
		s2 += fmt.Sprintf("Name: %s\n Image: %s\n Ports: %s\n VolumnMounts: %s\n", container.Name, container.Image, ports, volumnMounts)
	}
	return s1 + s2, nil
}

func RunDescribe(cmd *cobra.Command, args []string) {
	apiObjType := args[0]
	apiObjName := args[1]
	err := RunDescribe_Cmd(apiObjType, apiObjName)
	if err != nil {
		fmt.Println(err)
	}
}

func RunDescribe_Cmd(apiObjType string, apiObjName string) error {
	var val string
	var err error
	np := namespace
	switch apiObjType {
	case "pod":
		val, err = GetDetailedHelp(apiObjName, np)
		fmt.Println(val)
	default:
		fmt.Println("describe: not support this type")
	}
	return err
}
