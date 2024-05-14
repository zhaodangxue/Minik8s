package command

import (
	"fmt"
	ctlutils "minik8s/kubectl/utils"
	"strings"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:                        "get",
	Short:                      "get a resource",
	Long:                       `Get a resource by given source type and namespace.`,
	Run:                        RunGet,
	Example:                    "kubectl get pod xxx -n default",
	SuggestionsMinimumDistance: 1,
	SuggestFor:                 []string{"ge", "g"},
	Args:                       cobra.MinimumNArgs(1),
}

func RunGet(cmd *cobra.Command, args []string) {
	apiObjType := args[0]
	var apiObjName string
	if len(args) > 1 {
		apiObjName = args[1]
	}
	np := namespace
	var err error
	apiObjType = strings.ToLower(apiObjType)
	switch apiObjType {
	case "test":
		err = ctlutils.PrintTestStatusTable()
	case "pod":
		err = ctlutils.PrintPodStatusTable(np)
	case "pv":
		err = ctlutils.PrintPVTable(namespace)
	case "pvc":
		err = ctlutils.PrintPVCTable(namespace)
	case "node":
		fmt.Printf("Name:%s", apiObjName)
	}
	if err != nil {
		fmt.Println(err)
	}
}
func RunGet_test(apiObjType string, apiObjName string) error {
	np := namespace
	var err error
	apiObjType = strings.ToLower(apiObjType)
	switch apiObjType {
	case "test":
		err = ctlutils.PrintTestStatusTable()
	case "pod":
		err = ctlutils.PrintPodStatusTable(np)
	case "pv":
		err = ctlutils.PrintPVTable(namespace)
	case "pvc":
		err = ctlutils.PrintPVCTable(namespace)
	case "node":
		fmt.Printf("get node Name:%s", apiObjName)
	}
	return err
}
