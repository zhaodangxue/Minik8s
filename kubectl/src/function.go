package command

import (
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	ctlutils "minik8s/kubectl/utils"
	image "minik8s/serverless/image"
	"minik8s/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var funcCommand = &cobra.Command{
	Use:   "func",
	Short: "Manage the function",
	Long:  `Manage the function. The function is used in serveless.`,
	Run:   RunFunc,
	Args:  cobra.RangeArgs(1, 2),
}

func AddFunctionToApiServer(function apiobjects.Function) error {
	// Do Stuff Here
	url := route.Prefix + route.FunctionPath
	_, err := utils.PostWithJson(url, function)
	return err
}
func DeleteFunctionFromApiServer(name string) error {
	// Do Stuff Here
	url := route.Prefix + route.FunctionPath + "/" + namespace + "/" + name
	_, err := utils.Delete(url)
	return err
}
func RunFunc(cmd *cobra.Command, args []string) {
	// Do Stuff Here
	if len(args) == 0 {
		fmt.Println("Please input the func options")
		return
	}
	options := args[0]
	switch options {
	case "create":
		if filepath == "" {
			fmt.Println("Please input the yaml filepath")
			return
		}
		var data []byte
		data, err := ctlutils.LoadFile(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		FunctionCtlInput := apiobjects.FunctionCtlInput{}
		if err = yaml.Unmarshal(data, &FunctionCtlInput); err != nil {
			fmt.Println(err)
			return
		}
		imageUrl, err := image.CreateImage(FunctionCtlInput)
		if err != nil {
			fmt.Println(err)
			return
		}
		function := apiobjects.Function{}
		function.Object = FunctionCtlInput.Object
		function.Spec = FunctionCtlInput.FunctionSpec
		var FunctionStatus apiobjects.FunctionStatus
		FunctionStatus.ImageUrl = imageUrl
		function.Status = FunctionStatus
		err = AddFunctionToApiServer(function)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Create function success")
	case "delete":
		if len(args) != 2 {
			fmt.Println("Please input the function name")
			return
		}
		name := args[1]
		err := DeleteFunctionFromApiServer(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Delete function success")
	}
}
