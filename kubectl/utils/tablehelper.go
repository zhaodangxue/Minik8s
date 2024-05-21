package ctlutils

import (
	"minik8s/apiobjects"
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
	tbl := table.New("NAME", "NameSpace","TYPE", "CLUSTER-IP", "EXTERNAL-IP", "PORT(S)")
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
		var clusterIP string
		var externalIP string
		var ports string
		for _, p := range svc.Spec.Ports {
			ports += p.Name + ":" + strconv.FormatInt(int64(p.Port),10) + "/" + string(p.Protocol) + ", "
		}
		if svc.Status.ClusterIP == "" {
			clusterIP = "<none>"
		} else {
			clusterIP = svc.Status.ClusterIP
		}
		externalIP = "<none>"
		tbl.AddRow(svc.Data.Name, svc.Data.Namespace, svc.Spec.Type, clusterIP, externalIP, ports)
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
			paths += path.PathName + ":" + path.Service + " " + path.Address + ":"+ strconv.Itoa(path.Port) + ", "
		}
		tbl.AddRow(dnsRecord.Name, dnsRecord.NameSpace, dnsRecord.Host, paths)
	}
    tbl.Print()
	return nil
}
