package command

import (
	"fmt"
	"strings"

	ctlutils "minik8s/kubectl/utils"

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
	err := RunGet_Cmd(apiObjType, apiObjName)
	if err != nil {
		fmt.Println(err)
	}
}

func RunGet_Cmd(apiObjType string, apiObjName string) error {
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
	case "replicaset":
		err = ctlutils.PrintReplicasetTable(namespace)
	case "node":
		err = ctlutils.PrintNodeTable()
	case "service":
		err = ctlutils.PrintServiceTable()
	case "dns":
		err = ctlutils.PrintDNSTable()
	case "hpa":
		err = ctlutils.PrintHPATable(namespace)
	case "function":
		err = ctlutils.PrintFunctionTable()
	case "job":
		jobs, err := GetAllJobsInCluster();
		if err != nil {
			return err
		}
		err = ctlutils.PrintJobTable(jobs)
	}
	return err
}
