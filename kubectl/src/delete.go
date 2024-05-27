package command

import (
	"fmt"

	"minik8s/apiserver/src/route"
	"minik8s/utils"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:                        "delete",
	Short:                      "delete a resource",
	Long:                       `Delete a resource by given source type and name.`,
	Run:                        RunDelete,
	Example:                    "kubectl delete pod xxx -n default",
	SuggestionsMinimumDistance: 1,
	SuggestFor:                 []string{"del", "d"},
	Args:                       cobra.MinimumNArgs(2),
}

func deleteTest() (string, error) {
	url := route.Prefix + route.TestCtlPath
	val, err := utils.Delete(url)
	return val, err
}

func deleteSpecifiedPod(np, apiObjName string) (string, error) {
	url := route.Prefix + route.PodPath + "/" + np + "/" + apiObjName
	val, err := utils.Delete(url)
	return val, err
}

func deleteSpecifiedPV(np, apiObjName string) (string, error) {
	url := route.Prefix + route.PVPath + "/" + np + "/" + apiObjName
	val, err := utils.Delete(url)
	return val, err
}

func deleteSpecifiedPVC(np, apiObjName string) (string, error) {
	url := route.Prefix + route.PVCPath + "/" + np + "/" + apiObjName
	val, err := utils.Delete(url)
	return val, err
}

func deleteSpecifiedService(np, apiObjName string) (string, error) {
	url := route.Prefix + "/api/service/cmd/delete/" + np + "/" + apiObjName
	val, err := utils.Delete(url)
	return val, err
}

func deleteSpecifiedReplicaset(np, apiObjName string) (string, error) {
	url := route.Prefix + "/api/replicaset/" + np + "/" + apiObjName
	val, err := utils.Delete(url)
	return val, err
}

func deleteSpecifiedDNS(np, apiObjName string) (string, error) {
	url := route.Prefix + "/api/dns/delete/" + np + "/" + apiObjName
	val, err := utils.Delete(url)
	return val, err
}
func deleteSpecifiedHPA(np, apiObjName string) (string, error) {
	url := route.Prefix + route.HorizontalPodAutoscalerPath + "/" + np + "/" + apiObjName
	val, err := utils.Delete(url)
	return val, err
}
func RunDelete(cmd *cobra.Command, args []string) {
	apiObjType := args[0]
	apiObjName := args[1]
	err := RunDelete_Cmd(apiObjType, apiObjName)
	if err != nil {
		fmt.Println(err)
	}
}

func RunDelete_Cmd(apiObjType, apiObjName string) error {
	np := namespace
	var err error
	var val string
	switch apiObjType {
	case "test":
		_, err = deleteTest()
	case "pod":
		val, err = deleteSpecifiedPod(np, apiObjName)
		fmt.Println(val)
	case "pv":
		val, err = deleteSpecifiedPV(np, apiObjName)
		fmt.Println(val)
	case "pvc":
		val, err = deleteSpecifiedPVC(np, apiObjName)
		fmt.Println(val)
	case "replicaset":
		var val string
		val, err = deleteSpecifiedReplicaset(np, apiObjName)
		if err == nil {
			fmt.Println(val)
		}
	case "service":
		val, err = deleteSpecifiedService(np, apiObjName)
		if err == nil {
			fmt.Println(val)
		}
	case "dns":
		val, err = deleteSpecifiedDNS(np, apiObjName)
		if err == nil {
			fmt.Println(val)
		}
	case "hpa":
		val, err = deleteSpecifiedHPA(np, apiObjName)
		if err == nil {
			fmt.Println(val)
		}
	default:
		fmt.Println("delete: not support this type")
	}
	return err
}
