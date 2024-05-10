package ctlutils

import (
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/utils"

	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func TestTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("ID", "Name", "Replicas", "Kind")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func PodTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("NAME", "NAMESPACE", "UUID", "STATUS", "CREATION")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func NodeTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("NAME", "STATUS", "ROLES", "AGE", "VERSION")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func GetTestFromApiserver() (testyaml *apiobjects.TestYaml, err error) {
	url := route.Prefix + route.TestCtlPath
	err = utils.GetUnmarshal(url, &testyaml)
	return
}
func GetPodFromApiserver(namespace string) (pods []*apiobjects.Pod, err error) {
	url := route.Prefix + route.PodPath + "/" + namespace
	err = utils.GetUnmarshal(url, &pods)
	return
}

//类似这样的GetPodFromApiserver函数等等,第一个（）里填的是查找所需的参数

func PrintTestStatusTable() error {
	Teststatus, err := GetTestFromApiserver()
	if err != nil {
		return err
	}
	tbl := TestTbl()
	tbl.AddRow(1, Teststatus.Spec.Name, Teststatus.Spec.Replicas, Teststatus.Kind)
	tbl.Print()
	return nil
}
func PrintPodStatusTable(namespace string) error {
	pods, err := GetPodFromApiserver(namespace)
	if err != nil {
		return err
	}
	tbl := PodTbl()
	for _, pod := range pods {
		tbl.AddRow(pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, pod.ObjectMeta.UID, pod.Status.PodPhase, pod.ObjectMeta.CreationTimestamp.Format("2006-01-02 15:04:05"))
	}
	tbl.Print()
	return nil
}
