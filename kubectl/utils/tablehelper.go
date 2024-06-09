package ctlutils

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/utils"
	"strconv"

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
	tbl := table.New("NAME", "NAMESPACE", "UUID", "STATUS", "CREATION", "IP", "HOSTIP")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func NodeTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("NAME", "STATUS", "IP", "CREATION")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func PVCTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("NAME", "STATUS", "VOLUME", "CAPACITY", "ACCESSMODE", "CREATION")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func PVTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("NAME", "CAPACITY", "ACCESSMODE", "RECLAIM POLICY", "STATUS", "CREATION")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func ReplicasetTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("NAME", "DESIRED", "READY", "AGE")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func ServiceTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("NAME", "Selector", "TYPE", "CLUSTER-IP", "PORT(S)", "ENDPOINTS")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}

func DNSRecordTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Name", "Namespace", "Host", "Paths")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func HPARecordTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Name", "Namespace", "MinReplicas", "MaxReplicas", "TargetCPUUtilizationPercentage", "TargetMemoryUtilizationPercentage", "ScaleInterval")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func FunctionTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Name", "MinReplicas", "MaxReplicas", "TargetQPSPerReplica", "ImageUrl", "Creation")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	return tbl
}
func JobTbl() table.Table {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Name", "Namespace", "JobState", "Otuput", "ImageUrl", "PodIp", "PodPath")
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
func GetPVFromApiserver(namespace string) (pvs []*apiobjects.PersistentVolume, err error) {
	url := route.Prefix + route.PVPath + "/" + namespace
	err = utils.GetUnmarshal(url, &pvs)
	return
}
func GetPVCFromApiserver(namespace string) (pvcs []*apiobjects.PersistentVolumeClaim, err error) {
	url := route.Prefix + route.PVCPath + "/" + namespace
	err = utils.GetUnmarshal(url, &pvcs)
	return
}
func GetReplicasetFromApiserver(namespace string) (replicasets []*apiobjects.Replicaset, err error) {
	url := route.Prefix + route.ReplicasetPath + "/" + namespace
	err = utils.GetUnmarshal(url, &replicasets)
	return
}
func GetHPAFromApiserver(namespace string) (hpas []*apiobjects.HorizontalPodAutoscaler, err error) {
	url := route.Prefix + route.HorizontalPodAutoscalerPath + "/" + namespace
	err = utils.GetUnmarshal(url, &hpas)
	return
}
func GetNodeFromApiserver() (nodes []*apiobjects.Node, err error) {
	url := route.Prefix + route.NodePath
	err = utils.GetUnmarshal(url, &nodes)
	return
}
func GetFuncFromApiserver() (functions []*apiobjects.Function, err error) {
	url := route.Prefix + route.FunctionPath
	err = utils.GetUnmarshal(url, &functions)
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
		tbl.AddRow(pod.ObjectMeta.Name, pod.ObjectMeta.Namespace, pod.ObjectMeta.UID, pod.Status.PodPhase, pod.ObjectMeta.CreationTimestamp.Format("2006-01-02 15:04:05"), pod.Status.PodIP, pod.Status.HostIP)
	}
	tbl.Print()
	return nil
}
func PrintPVTable(namespace string) error {
	pvs, err := GetPVFromApiserver(namespace)
	if err != nil {
		return err
	}
	tbl := PVTbl()
	for _, pv := range pvs {
		var accessMode string
		for _, mode := range pv.Spec.AccessModes {
			if mode == "ReadWriteOnce" {
				accessMode += "RWO" + " "
			} else if mode == "ReadOnlyMany" {
				accessMode += "ROX" + " "
			} else if mode == "ReadWriteMany" {
				accessMode += "RWX" + " "
			}
		}
		tbl.AddRow(pv.ObjectMeta.Name, pv.Spec.Capacity.Storage, accessMode, pv.Spec.PersistentVolumeReclaimPolicy, pv.Status, pv.CreationTimestamp.Format("2006-01-02 15:04:05"))
	}
	tbl.Print()
	return nil
}
func PrintPVCTable(namespace string) error {
	pvcs, err := GetPVCFromApiserver(namespace)
	if err != nil {
		return err
	}
	tbl := PVCTbl()
	for _, pvc := range pvcs {
		var accessMode string
		if pvc.Spec.AccessModes[0] == "ReadWriteOnce" {
			accessMode = "RWO"
		} else if pvc.Spec.AccessModes[0] == "ReadOnlyMany" {
			accessMode = "ROX"
		} else if pvc.Spec.AccessModes[0] == "ReadWriteMany" {
			accessMode = "RWX"
		}
		tbl.AddRow(pvc.ObjectMeta.Name, pvc.Status, pvc.PVBinding.PVname, pvc.PVBinding.PVcapacity, accessMode, pvc.CreationTimestamp.Format("2006-01-02 15:04:05"))
	}
	tbl.Print()
	return nil
}

func PrintServiceTable() error {
	svcs := []*apiobjects.Service{}
	url := route.Prefix + route.GetAllServicesPath
	err := utils.GetUnmarshal(url, &svcs)
	if err != nil {
		return err
	}
	tbl := ServiceTbl()
	for _, svc := range svcs {
		edpts := []*apiobjects.Endpoint{}
		values, err := etcd.Get_prefix("/api/endpoint/" + svc.Data.Name + "/" + svc.Data.Namespace)
		if err != nil {
			fmt.Println(err)
		}
		PodIp := ""
		for _, value := range values {
			var endpoint apiobjects.Endpoint
			err := json.Unmarshal([]byte(value), &endpoint)
			if err != nil {
				fmt.Println(err)
			}
			PodIp += endpoint.Spec.DestIP + ":" + strconv.FormatInt(int64(endpoint.Spec.DestPort), 10) + ","
			edpts = append(edpts, &endpoint)
		}
		var clusterIP string
		ports := ""
		selector := ""
		for key, value := range svc.Spec.Selector {
			selector += key + ":" + value + ", "
		}

		for _, p := range svc.Spec.Ports {
			ports += p.Name + ":" + strconv.FormatInt(int64(p.Port), 10) + "/ TargetPort:" + p.TargetPort + ", "
		}
		if svc.Status.ClusterIP == "" {
			clusterIP = "<none>"
		} else {
			clusterIP = svc.Status.ClusterIP
		}

		if PodIp == "" {
			PodIp = "<none>"
		}

		tbl.AddRow(svc.Data.Name, selector, svc.Spec.Type, clusterIP, ports, PodIp)
	}
	tbl.Print()
	return nil
}

func PrintReplicasetTable(namespace string) error {
	replicasets, err := GetReplicasetFromApiserver(namespace)
	if err != nil {
		return err
	}
	tbl := ReplicasetTbl()
	for _, replicaset := range replicasets {
		tbl.AddRow(replicaset.Name, replicaset.Spec.Replicas, replicaset.Spec.Ready, replicaset.CreationTimestamp.Format("2006-01-02 15:04:05"))
	}
	tbl.Print()
	return nil
}

func PrintDNSTable() error {
	//TODO
	dnsRecords := []*apiobjects.DNSRecord{}
	url := route.Prefix + route.DnsGetAllPath
	err := utils.GetUnmarshal(url, &dnsRecords)
	if err != nil {
		return err
	}
	tbl := DNSRecordTbl()
	for _, dnsRecord := range dnsRecords {
		var paths string
		for _, path := range dnsRecord.Paths {
			paths += path.PathName + ":" + path.Service + " " + path.Address + ":" + strconv.Itoa(path.Port) + ", "
		}
		tbl.AddRow(dnsRecord.Name, dnsRecord.NameSpace, dnsRecord.Host, paths)
	}
	tbl.Print()
	return nil
}
func PrintHPATable(namespace string) error {
	hpas, err := GetHPAFromApiserver(namespace)
	if err != nil {
		return err
	}
	tbl := HPARecordTbl()
	for _, hpa := range hpas {
		strMem_HPA := fmt.Sprintf("%f", hpa.Stat.CurrentReplicaseMemUsage)
		expected_HPA := fmt.Sprintf("%f", hpa.Spec.Metrics.MemoryUtilizationUsage)
		tbl.AddRow(hpa.Name, hpa.Namespace, hpa.Spec.MinReplicas, hpa.Spec.MaxReplicas, strconv.Itoa(hpa.Stat.CurrnentReplicaseCPUUsage)+"/"+strconv.Itoa(hpa.Spec.Metrics.CPUUtilizationPercentage), strMem_HPA+"/"+expected_HPA, hpa.Spec.ScaleInterval)
	}
	tbl.Print()
	return nil
}
func PrintNodeTable() error {
	nodes, err := GetNodeFromApiserver()
	if err != nil {
		return err
	}
	tbl := NodeTbl()
	for _, node := range nodes {
		tbl.AddRow(node.Name, node.Status.State, node.Info.Ip, node.CreationTimestamp.Format("2006-01-02 15:04:05"))
	}
	tbl.Print()
	return nil
}
func PrintFunctionTable() error {
	functions, err := GetFuncFromApiserver()
	if err != nil {
		return err
	}
	tbl := FunctionTbl()
	for _, function := range functions {
		tbl.AddRow(function.Name, function.Spec.MinReplicas, function.Spec.MaxReplicas, function.Spec.TargetQPSPerReplica, function.Status.ImageUrl, function.CreationTimestamp.Format("2006-01-02 15:04:05"))
	}
	tbl.Print()
	return nil
}

func PrintJobTable(jobs []*apiobjects.Job) error {
	tbl := JobTbl()
	for _, job := range jobs {
		tbl.AddRow(job.Name, job.Namespace, job.Status.JobState, job.Status.Output, 
			job.Status.ImageUrl, job.Status.PodIp, job.Status.PodRef.GetObjectPath())
	}
	tbl.Print()
	return nil
}
