package command

import (
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	jobimage "minik8s/jobserver/image"
	ctlutils "minik8s/kubectl/utils"
	"minik8s/utils"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var jobCommand = &cobra.Command{
	Use:     "job",
	Short:   "Manage the job",
	Long:    `Manage the job. The job is used in gpu.`,
	Run:     RunJob,
	Args:    cobra.RangeArgs(1, 2),
	Example: `kubectl job create -f ./example_job.yaml`,
}

func RunJob(cmd *cobra.Command, args []string) {
	// Do Stuff Here
	if len(args) == 0 {
		fmt.Println("Please input the job options")
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
		job := apiobjects.Job{}
		if err = yaml.Unmarshal(data, &job); err != nil {
			fmt.Println(err)
			return
		}
		imageUrl, err := jobimage.CreateImage(&job)
		if err != nil {
			fmt.Println(err)
			return
		}
		job.Status.ImageUrl = imageUrl
		job.Status.JobState = apiobjects.JobState_Pending
		err = createJobInCluster(&job)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "get":
		// TODO: 获取单个job的状态
	case "clean":
		// TODO: clean finished/failed job
	}
}

func createJobInCluster(job *apiobjects.Job) error {
	url := route.Prefix + route.JobPath
	_, err := utils.PostWithJson(url, job)
	return err
}

func GetJobInCluster(namespace, name string) (*apiobjects.Job, error) {
	// CHECK: 修改url命名规则时，这里会失效
	url := route.Prefix + route.JobPath + "/" + namespace + "/" + name
	job := &apiobjects.Job{}
	err := utils.GetUnmarshal(url, job)
	return job, err
}

func GetAllJobsInCluster() ([]*apiobjects.Job, error) {
	url := route.Prefix + route.JobPath
	jobs := []*apiobjects.Job{}
	err := utils.GetUnmarshal(url, &jobs)
	return jobs, err
}
