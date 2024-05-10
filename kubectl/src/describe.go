package command

import (
	"fmt"

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

func RunDescribe(cmd *cobra.Command, args []string) {
	apiObjType := args[0]
	apiObjName := args[1]
	np := namespace
	var err error
	switch apiObjType {
	case "pod":
		err = nil
		fmt.Printf("describe pod Name:%s, Namespace:%s", apiObjName, np)
	default:
		fmt.Println("describe: not support this type")
	}
	if err != nil {
		fmt.Println(err)
	}
}
