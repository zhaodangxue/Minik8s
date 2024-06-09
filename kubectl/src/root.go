package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var example = "kubectl apply | kubectl get | kubectl describe"
var testpath string
var filepath string
var namespace string
var rootCmd = &cobra.Command{
	Use:     "kubectl",
	Short:   "kubectl is a command line tool for interacting with Kubernetes clusters",
	Long:    `kubectl controls the Kubernetes cluster manager. For example: kubectl apply -f ./example.yaml; kubectl describe pod examplePod`,
	Version: "v1.0.0",
	Run:     RunRoot,
	Example: example,
}

func init() {
	rootCmd.Flags().StringVarP(&testpath, "testpath", "t", "", "this is a test path")
	applyCmd.Flags().StringVarP(&filepath, "filepath", "f", "", "input a yaml filepath")
	getCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "input a namespace")
	describeCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "input a namespace")
	deleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "input a namespace")
	wfCommand.Flags().StringVarP(&filepath, "filepath", "f", "", "input a json filepath")
	wfCommand.Flags().StringVarP(&namespace, "namespace", "n", "default", "input a namespace")
	funcCommand.Flags().StringVarP(&filepath, "filepath", "f", "", "input a yaml filepath")
	funcCommand.Flags().StringVarP(&namespace, "namespace", "n", "default", "input a namespace")
	jobCommand.Flags().StringVarP(&filepath, "filepath", "f", "", "input a yaml filepath") // use in create
	jobCommand.Flags().StringVarP(&namespace, "namespace", "n", "default", "input a namespace")	// use in get
	eventCommand.Flags().StringVarP(&filepath, "filepath", "f", "", "input a json filepath")
	eventCommand.Flags().StringVarP(&namespace, "namespace", "n", "default", "input a namespace")
	applyCmd.MarkFlagRequired("filepath")
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(schedCommand)
	rootCmd.AddCommand(wfCommand)
	rootCmd.AddCommand(funcCommand)
	rootCmd.AddCommand(jobCommand)
	rootCmd.AddCommand(eventCommand)
}

func RunRoot(cmd *cobra.Command, args []string) {
	fmt.Println("this root")
	fmt.Println("testpath:", testpath)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
