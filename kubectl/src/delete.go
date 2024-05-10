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
func RunDelete(cmd *cobra.Command, args []string) {
	apiObjType := args[0]
	apiObjName := args[1]
	np := namespace
	var err error
	switch apiObjType {
	case "test":
		var val string
		val, err = deleteTest()
		if err == nil {
			fmt.Println(val)
		}
	case "pod":
		fmt.Printf("Delete pod Name:%s, Namespace:%s", apiObjName, np)
	default:
		fmt.Println("delete: not support this type")
	}
	if err != nil {
		fmt.Println(err)
	}
}
