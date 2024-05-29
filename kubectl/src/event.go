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

var eventCommand = &cobra.Command{
	Use:     "event",
	Short:   "Manage the event",
	Long:    `Manage the event. The event is used in serveless to triger workflow.`,
	Run:     RunEvent,
	Args:    cobra.RangeArgs(1, 2),
	Example: `kubectl event create -f ./event.json`,
}

func RunEvent(cmd *cobra.Command, args []string) {
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
		event := apiobjects.Event{}
		if err = json.Unmarshal(data, &event); err != nil {
			fmt.Println(err)
			return
		}
		err = AddEventToApiServer(event)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "delete":
		if len(args) < 2 {
			fmt.Println("Please input the event name")
			return
		}
		eventName := args[1]
		err := DeleteEventFromApiServer(eventName)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}

func AddEventToApiServer(event apiobjects.Event) error {
	// Do Stuff Here
	url := route.Prefix + route.EventPath
	_, err := utils.PostWithJson(url, event)
	return err
}

func DeleteEventFromApiServer(eventName string) error {
	// Do Stuff Here
	url := route.Prefix + route.EventPath + "/" + eventName
	_, err := utils.Delete(url)
	return err
}
