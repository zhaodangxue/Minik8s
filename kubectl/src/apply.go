package command

import (
	"fmt"

	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	ctlutils "minik8s/kubectl/utils"
	"minik8s/utils"
	//"net/url"

	//"google.golang.org/grpc/internal/resolver/dns"
	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:                        "apply",
	Short:                      "apply a configuration to a resource by filename",
	Long:                       `Apply a configuration to a resource by filename. The resource will be created if it doesn't exist.`,
	RunE:                       RunApply,
	Example:                    "kubectl apply -f ./example.yaml",
	SuggestionsMinimumDistance: 1,
	SuggestFor:                 []string{"aply", "applying", "a"},
}

func RunApply(cmd *cobra.Command, args []string) error {
	err := RunApply_Cmd(filepath)
	return err
}

func RunApply_Cmd(file_path string) error {
	content, err := ctlutils.LoadFile(file_path)
	if err != nil {
		fmt.Println(err)
		return err
	}
	t := ctlutils.ParseApiObjectType(content)
	switch t {
	case ctlutils.Test:
		test := apiobjects.TestYaml{}
		err = yaml.Unmarshal(content, &test)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := route.Prefix + route.TestCtlPath
		utils.ApplyApiObject(url, test)
	case ctlutils.Pod:
		pod := apiobjects.Pod{}
		err = yaml.Unmarshal(content, &pod)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := route.Prefix + route.PodPath
		utils.ApplyApiObject(url, pod)
	case ctlutils.Pv:
		pv := apiobjects.PersistentVolume{}
		err = yaml.Unmarshal(content, &pv)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := route.Prefix + route.PVPath
		utils.ApplyApiObject(url, pv)
	case ctlutils.Pvc:
		pvc := apiobjects.PersistentVolumeClaim{}
		err = yaml.Unmarshal(content, &pvc)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := route.Prefix + route.PVCPath
		utils.ApplyApiObject(url, pvc)
		// case ctlutils.Node:
		// default:
	case ctlutils.Replicaset:
		replicaset := apiobjects.Replicaset{}
		err = yaml.Unmarshal(content, &replicaset)
		pod := apiobjects.ToPod(&replicaset.Spec.Template)
		fmt.Println("pod: ", pod)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := route.Prefix + route.ReplicasetPath
		utils.ApplyApiObject(url, replicaset)
	case ctlutils.Service:
		service := apiobjects.Service{}
		err = yaml.Unmarshal(content, &service)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := route.Prefix + route.ServiceApplyPath
		fmt.Println("appyly service: ", service)
		// TODO service格式是否符合要求
		utils.ApplyApiObject(url, service)
	case ctlutils.Dns:
		dns := apiobjects.DNSRecord{}
		err = yaml.Unmarshal(content, &dns)
		if err != nil {
			fmt.Println(err)
			return err
		}
		url := route.Prefix + "/api/dns/apply"
		fmt.Println("appyly service: ", dns)
		//TODO service格式是否符合要求
		utils.ApplyApiObject(url, dns)
	}
	return nil
}
