package command

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	ctlutils "minik8s/kubectl/utils"
	"minik8s/utils"

	"github.com/spf13/cobra"
)

var wfCommand = &cobra.Command{
	Use:                        "workflow",
	Short:                      "Manage the workflow",
	Long:                       `Manage the workflow. The workflow is used to determine how functions are scheduled in a specified workflow.`,
	Run:                        RunWf,
	Args:                       cobra.RangeArgs(1, 2),
	SuggestionsMinimumDistance: 1,
	SuggestFor:                 []string{"wf", "workflows"},
	Example:                    "kubectl workflow create -f ./example.json",
}

func AddWorkflowToApiServer(workflow apiobjects.Workflow) error {
	// Do Stuff Here
	url := route.Prefix + route.WorkflowPath
	_, err := utils.PostWithJson(url, workflow)
	return err
}
func DeleteWorkflowFromApiServer(workflowName string) error {
	// Do Stuff Here
	url := route.Prefix + route.WorkflowPath + "/" + namespace + "/" + workflowName
	_, err := utils.Delete(url)
	return err
}
func RunWf(cmd *cobra.Command, args []string) {
	// Do Stuff Here
	options := args[0]
	var err error
	switch options {
	case "create":
		var data []byte
		data, err = ctlutils.LoadFile(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		workflow := apiobjects.Workflow{}
		if err = json.Unmarshal(data, &workflow); err != nil {
			fmt.Println(err)
			return
		}
		err = AddWorkflowToApiServer(workflow)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		if len(args) == 1 {
			fmt.Println("Please specify the workflow name")
			return
		}
		workflowName := args[1]
		err = DeleteWorkflowFromApiServer(workflowName)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Deleting workflow", workflowName)
	}
}
