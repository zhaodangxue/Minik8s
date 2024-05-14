package PV_test

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/apiserver"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"minik8s/controller/src/PVcontroller"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreatePV(t *testing.T) {
	fmt.Println("TestCreatePV")
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := PVcontroller.New()
	go controller.Run()
	fmt.Println("apiServer is running")
	etcd.Clear()
	pv := apiobjects.PersistentVolume{}
	pv.Kind = "PersistentVolume"
	pv.Dynamic = 0
	pv.ApiVersion = "v1"
	pv.ObjectMeta.Name = "pv"
	pv.ObjectMeta.Namespace = "default"
	pv.ObjectMeta.CreationTimestamp = time.Now()
	pv.ObjectMeta.UID = utils.NewUUID()
	pv.Spec.Capacity.Storage = "1Gi"
	pv.Spec.AccessModes = []string{"ReadWriteMany"}
	pv.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv.Spec.VolumeMode = "Filesystem"
	pv.Spec.StorageClassName = "slow"
	pv.Spec.NFS.Server = global.Nfsserver
	pv.Spec.NFS.Path = global.NFSdir + "/aaa"
	pvJson, err := json.Marshal(pv)
	assert.Nil(t, err)
	etcd.Put(pv.GetObjectPath()+"/"+pv.Spec.StorageClassName, string(pvJson))
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
}
func TestCreatePVCTOBoundSpecifiedPV(t *testing.T) {
	fmt.Println("TestCreatePVCTOBoundSpecifiedPV")
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := PVcontroller.New()
	go controller.Run()
	fmt.Println("apiServer is running")
	etcd.Clear()
	pv1 := apiobjects.PersistentVolume{}
	pv1.Kind = "PersistentVolume"
	pv1.ApiVersion = "v1"
	pv1.ObjectMeta.Name = "pv1"
	pv1.ObjectMeta.Namespace = "default"
	pv1.Dynamic = 0
	pv1.ObjectMeta.CreationTimestamp = time.Now()
	pv1.ObjectMeta.UID = utils.NewUUID()
	pv1.Spec.Capacity.Storage = "3Gi"
	pv1.Spec.AccessModes = []string{"ReadWriteMany"}
	pv1.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv1.Spec.VolumeMode = "Filesystem"
	pv1.Spec.StorageClassName = "default"
	pv1.Spec.NFS.Server = global.Nfsserver
	pv1.Spec.NFS.Path = global.NFSdir + "/" + pv1.ObjectMeta.Name
	pvJson, err := json.Marshal(pv1)
	assert.Nil(t, err)
	etcd.Put(pv1.GetObjectPath()+"/"+pv1.Spec.StorageClassName, string(pvJson))
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	time.Sleep(2 * time.Second)
	pv2 := apiobjects.PersistentVolume{}
	pv2.Kind = "PersistentVolume"
	pv2.ApiVersion = "v1"
	pv2.ObjectMeta.Name = "pv2"
	pv2.ObjectMeta.Namespace = "default"
	pv2.Dynamic = 0
	pv2.ObjectMeta.CreationTimestamp = time.Now()
	pv2.Spec.Capacity.Storage = "1Gi"
	pv2.ObjectMeta.UID = utils.NewUUID()
	pv2.Spec.AccessModes = []string{"ReadWriteMany"}
	pv2.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv2.Spec.VolumeMode = "Filesystem"
	pv2.Spec.StorageClassName = "slow"
	pv2.Spec.NFS.Server = global.Nfsserver
	pv2.Spec.NFS.Path = global.NFSdir + "/" + pv2.ObjectMeta.Name
	pvJson, err = json.Marshal(pv2)
	assert.Nil(t, err)
	etcd.Put(pv2.GetObjectPath()+"/"+pv2.Spec.StorageClassName, string(pvJson))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	time.Sleep(2 * time.Second)
	pv3 := apiobjects.PersistentVolume{}
	pv3.Kind = "PersistentVolume"
	pv3.ApiVersion = "v1"
	pv3.ObjectMeta.Name = "pv3"
	pv3.ObjectMeta.Namespace = "default"
	pv3.Dynamic = 0
	pv3.ObjectMeta.CreationTimestamp = time.Now()
	pv3.Spec.Capacity.Storage = "3Gi"
	pv3.ObjectMeta.UID = utils.NewUUID()
	pv3.Spec.AccessModes = []string{"ReadWriteMany"}
	pv3.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv3.Spec.VolumeMode = "Filesystem"
	pv3.Spec.StorageClassName = "slow"
	pv3.Spec.NFS.Server = global.Nfsserver
	pv3.Spec.NFS.Path = global.NFSdir + "/" + pv3.ObjectMeta.Name
	pvJson, err = json.Marshal(pv3)
	assert.Nil(t, err)
	etcd.Put(pv3.GetObjectPath()+"/"+pv3.Spec.StorageClassName, string(pvJson))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	time.Sleep(2 * time.Second)
	pvc := apiobjects.PersistentVolumeClaim{}
	pvc.Kind = "PersistentVolumeClaim"
	pvc.ApiVersion = "v1"
	pvc.ObjectMeta.Name = "pvc"
	pvc.ObjectMeta.Namespace = "default"
	pvc.ObjectMeta.CreationTimestamp = time.Now()
	pvc.ObjectMeta.UID = utils.NewUUID()
	pvc.Spec.Resources.Requests.Storage = "2Gi"
	pvc.Spec.AccessModes = []string{"ReadWriteMany"}
	pvc.Spec.StorageClassName = "slow"
	var pvcJson []byte
	pvcJson, err = json.Marshal(pvc)
	assert.Nil(t, err)
	etcd.Put(pvc.GetObjectPath()+"/"+pvc.Spec.StorageClassName, string(pvcJson))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvcJson)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
	val, _ := etcd.Get(pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName)
	var pvc1 apiobjects.PersistentVolumeClaim
	err = json.Unmarshal([]byte(val), &pvc1)
	assert.Nil(t, err)
	assert.Equal(t, apiobjects.PVCPhase("Bound"), pvc1.Status)
	assert.Equal(t, "pv3", pvc1.PVBinding.PVname)
	assert.Equal(t, "default", pvc1.PVBinding.PVnamespace)
	assert.Equal(t, "3Gi", pvc1.PVBinding.PVcapacity)
	assert.Equal(t, pv3.Spec.NFS.Path, pvc1.PVBinding.PVpath)
}
func TestPVCDelete(t *testing.T) {
	fmt.Println("TestPVCDelete")
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := PVcontroller.New()
	go controller.Run()
	fmt.Println("apiServer is running")
	etcd.Clear()
	pv1 := apiobjects.PersistentVolume{}
	pv1.Kind = "PersistentVolume"
	pv1.ApiVersion = "v1"
	pv1.ObjectMeta.Name = "pv1"
	pv1.ObjectMeta.Namespace = "default"
	pv1.Dynamic = 0
	pv1.ObjectMeta.CreationTimestamp = time.Now()
	pv1.ObjectMeta.UID = utils.NewUUID()
	pv1.Spec.Capacity.Storage = "3Gi"
	pv1.Spec.AccessModes = []string{"ReadWriteMany"}
	pv1.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv1.Spec.VolumeMode = "Filesystem"
	pv1.Spec.StorageClassName = "default"
	pv1.Spec.NFS.Server = global.Nfsserver
	pv1.Spec.NFS.Path = global.NFSdir + "/" + pv1.ObjectMeta.Name
	pvJson, err := json.Marshal(pv1)
	assert.Nil(t, err)
	etcd.Put(pv1.GetObjectPath()+"/"+pv1.Spec.StorageClassName, string(pvJson))
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	time.Sleep(2 * time.Second)
	pv2 := apiobjects.PersistentVolume{}
	pv2.Kind = "PersistentVolume"
	pv2.ApiVersion = "v1"
	pv2.ObjectMeta.Name = "pv2"
	pv2.ObjectMeta.Namespace = "default"
	pv2.Dynamic = 0
	pv2.ObjectMeta.CreationTimestamp = time.Now()
	pv2.Spec.Capacity.Storage = "1Gi"
	pv2.ObjectMeta.UID = utils.NewUUID()
	pv2.Spec.AccessModes = []string{"ReadWriteMany"}
	pv2.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv2.Spec.VolumeMode = "Filesystem"
	pv2.Spec.StorageClassName = "slow"
	pv2.Spec.NFS.Server = global.Nfsserver
	pv2.Spec.NFS.Path = global.NFSdir + "/" + pv2.ObjectMeta.Name
	pvJson, err = json.Marshal(pv2)
	assert.Nil(t, err)
	etcd.Put(pv2.GetObjectPath()+"/"+pv2.Spec.StorageClassName, string(pvJson))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	time.Sleep(2 * time.Second)
	pv3 := apiobjects.PersistentVolume{}
	pv3.Kind = "PersistentVolume"
	pv3.ApiVersion = "v1"
	pv3.ObjectMeta.Name = "pv3"
	pv3.ObjectMeta.Namespace = "default"
	pv3.Dynamic = 0
	pv3.ObjectMeta.CreationTimestamp = time.Now()
	pv3.Spec.Capacity.Storage = "3Gi"
	pv3.ObjectMeta.UID = utils.NewUUID()
	pv3.Spec.AccessModes = []string{"ReadWriteMany"}
	pv3.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv3.Spec.VolumeMode = "Filesystem"
	pv3.Spec.StorageClassName = "slow"
	pv3.Spec.NFS.Server = global.Nfsserver
	pv3.Spec.NFS.Path = global.NFSdir + "/" + pv3.ObjectMeta.Name
	pvJson, err = json.Marshal(pv3)
	assert.Nil(t, err)
	etcd.Put(pv3.GetObjectPath()+"/"+pv3.Spec.StorageClassName, string(pvJson))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvJson)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PvRelevantTopic(), string(topicMessageJson))
	time.Sleep(2 * time.Second)
	pvc := apiobjects.PersistentVolumeClaim{}
	pvc.Kind = "PersistentVolumeClaim"
	pvc.ApiVersion = "v1"
	pvc.ObjectMeta.Name = "pvc"
	pvc.ObjectMeta.Namespace = "default"
	pvc.ObjectMeta.CreationTimestamp = time.Now()
	pvc.ObjectMeta.UID = utils.NewUUID()
	pvc.Spec.Resources.Requests.Storage = "2Gi"
	pvc.Spec.AccessModes = []string{"ReadWriteMany"}
	pvc.Spec.StorageClassName = "slow"
	var pvcJson []byte
	pvcJson, err = json.Marshal(pvc)
	assert.Nil(t, err)
	etcd.Put(pvc.GetObjectPath()+"/"+pvc.Spec.StorageClassName, string(pvcJson))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvcJson)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
	val, _ := etcd.Get(pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName)
	var pvc1 apiobjects.PersistentVolumeClaim
	err = json.Unmarshal([]byte(val), &pvc1)
	assert.Nil(t, err)
	assert.Equal(t, apiobjects.PVCPhase("Bound"), pvc1.Status)
	assert.Equal(t, "pv3", pvc1.PVBinding.PVname)
	assert.Equal(t, "default", pvc1.PVBinding.PVnamespace)
	assert.Equal(t, "3Gi", pvc1.PVBinding.PVcapacity)
	assert.Equal(t, pv3.Spec.NFS.Path, pvc1.PVBinding.PVpath)
	val, err = utils.Delete(route.Prefix + pvc1.GetObjectPath())
	assert.Nil(t, err)
	fmt.Println(val)
	time.Sleep(3 * time.Second)
}
func TestDynamicAllocatePV(t *testing.T) {
	fmt.Println("TestDynamicAllocatePV")
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := PVcontroller.New()
	go controller.Run()
	fmt.Println("apiServer is running")
	etcd.Clear()
	pvc := apiobjects.PersistentVolumeClaim{}
	pvc.Kind = "PersistentVolumeClaim"
	pvc.ApiVersion = "v1"
	pvc.ObjectMeta.Name = "pvc"
	pvc.ObjectMeta.Namespace = "default"
	pvc.ObjectMeta.CreationTimestamp = time.Now()
	pvc.ObjectMeta.UID = utils.NewUUID()
	pvc.Spec.Resources.Requests.Storage = "2Gi"
	pvc.Spec.AccessModes = []string{"ReadWriteMany"}
	pvc.Spec.StorageClassName = "slow"
	var pvcJson []byte
	pvcJson, err := json.Marshal(pvc)
	assert.Nil(t, err)
	etcd.Put(pvc.GetObjectPath()+"/"+pvc.Spec.StorageClassName, string(pvcJson))
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvcJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
}
func TestPodUpdate(t *testing.T) {
	fmt.Println("TestPodUpdate")
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := PVcontroller.New()
	go controller.Run()
	fmt.Println("apiServer is running")
	etcd.Clear()
	pvc := apiobjects.PersistentVolumeClaim{}
	pvc.Kind = "PersistentVolumeClaim"
	pvc.ApiVersion = "v1"
	pvc.ObjectMeta.Name = "pvc"
	pvc.ObjectMeta.Namespace = "default"
	pvc.ObjectMeta.CreationTimestamp = time.Now()
	pvc.ObjectMeta.UID = utils.NewUUID()
	pvc.Spec.Resources.Requests.Storage = "2Gi"
	pvc.Spec.AccessModes = []string{"ReadWriteMany"}
	pvc.Spec.StorageClassName = "slow"
	var pvcJson []byte
	pvcJson, err := json.Marshal(pvc)
	assert.Nil(t, err)
	etcd.Put(pvc.GetObjectPath()+"/"+pvc.Spec.StorageClassName, string(pvcJson))
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvcJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
	pod1 := apiobjects.Pod{}
	pod1.ApiVersion = "v1"
	pod1.Kind = "Pod"
	pod1.ObjectMeta.Name = "pod1"
	pod1.ObjectMeta.Namespace = "default"
	pod1.ObjectMeta.CreationTimestamp = time.Now()
	pod1.ObjectMeta.UID = utils.NewUUID()
	volume1 := apiobjects.Volume{}
	volume1.Name = "volume1"
	var pvcSource apiobjects.PersistentVolumeClaimVolumeSource
	pvcSource.ClaimName = "pvc"
	pvcSource.ClaimNamespace = "default"
	volume1.PersistentVolumeClaim = &pvcSource
	var volumes []apiobjects.Volume
	volumes = append(volumes, volume1)
	pod1.Spec.Volumes = volumes
	pod1Json, err := json.Marshal(pod1)
	assert.Nil(t, err)
	etcd.Put(pod1.GetObjectPath(), string(pod1Json))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pod1Json)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PodRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
	pod2 := apiobjects.Pod{}
	pod2.ApiVersion = "v1"
	pod2.Kind = "Pod"
	pod2.ObjectMeta.Name = "pod2"
	pod2.ObjectMeta.Namespace = "default"
	pod2.ObjectMeta.CreationTimestamp = time.Now()
	pod2.ObjectMeta.UID = utils.NewUUID()
	var volume2 apiobjects.Volume
	volume2.Name = "volume2"
	var pvcSource2 apiobjects.PersistentVolumeClaimVolumeSource
	pvcSource2.ClaimName = "pvc"
	pvcSource2.ClaimNamespace = "default"
	volume2.PersistentVolumeClaim = &pvcSource2
	volumes = nil
	volumes = append(volumes, volume2)
	pod2.Spec.Volumes = volumes
	pod2Json, err := json.Marshal(pod2)
	assert.Nil(t, err)
	etcd.Put(pod2.GetObjectPath(), string(pod2Json))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pod2Json)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PodRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
}
func TestPodDelete(t *testing.T) {
	fmt.Println("TestPodUpdate")
	apiServer := apiserver.New()
	go apiServer.RUN()
	time.Sleep(3 * time.Second)
	controller := PVcontroller.New()
	go controller.Run()
	fmt.Println("apiServer is running")
	etcd.Clear()
	pvc := apiobjects.PersistentVolumeClaim{}
	pvc.Kind = "PersistentVolumeClaim"
	pvc.ApiVersion = "v1"
	pvc.ObjectMeta.Name = "pvc"
	pvc.ObjectMeta.Namespace = "default"
	pvc.ObjectMeta.CreationTimestamp = time.Now()
	pvc.ObjectMeta.UID = utils.NewUUID()
	pvc.Spec.Resources.Requests.Storage = "2Gi"
	pvc.Spec.AccessModes = []string{"ReadWriteMany"}
	pvc.Spec.StorageClassName = "slow"
	var pvcJson []byte
	pvcJson, err := json.Marshal(pvc)
	assert.Nil(t, err)
	etcd.Put(pvc.GetObjectPath()+"/"+pvc.Spec.StorageClassName, string(pvcJson))
	var topicMessage apiobjects.TopicMessage
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pvcJson)
	topicMessageJson, _ := json.Marshal(topicMessage)
	listwatch.Publish(global.PvcRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
	pod1 := apiobjects.Pod{}
	pod1.ApiVersion = "v1"
	pod1.Kind = "Pod"
	pod1.ObjectMeta.Name = "pod1"
	pod1.ObjectMeta.Namespace = "default"
	pod1.ObjectMeta.CreationTimestamp = time.Now()
	pod1.ObjectMeta.UID = utils.NewUUID()
	volume1 := apiobjects.Volume{}
	volume1.Name = "volume1"
	var pvcSource apiobjects.PersistentVolumeClaimVolumeSource
	pvcSource.ClaimName = "pvc"
	pvcSource.ClaimNamespace = "default"
	volume1.PersistentVolumeClaim = &pvcSource
	var volumes []apiobjects.Volume
	volumes = append(volumes, volume1)
	pod1.Spec.Volumes = volumes
	pod1Json, err := json.Marshal(pod1)
	assert.Nil(t, err)
	etcd.Put(pod1.GetObjectPath(), string(pod1Json))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pod1Json)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PodRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
	pod2 := apiobjects.Pod{}
	pod2.ApiVersion = "v1"
	pod2.Kind = "Pod"
	pod2.ObjectMeta.Name = "pod2"
	pod2.ObjectMeta.Namespace = "default"
	pod2.ObjectMeta.CreationTimestamp = time.Now()
	pod2.ObjectMeta.UID = utils.NewUUID()
	var volume2 apiobjects.Volume
	volume2.Name = "volume2"
	var pvcSource2 apiobjects.PersistentVolumeClaimVolumeSource
	pvcSource2.ClaimName = "pvc"
	pvcSource2.ClaimNamespace = "default"
	volume2.PersistentVolumeClaim = &pvcSource2
	volumes = nil
	volumes = append(volumes, volume2)
	pod2.Spec.Volumes = volumes
	pod2Json, err := json.Marshal(pod2)
	assert.Nil(t, err)
	etcd.Put(pod2.GetObjectPath(), string(pod2Json))
	topicMessage.ActionType = apiobjects.Create
	topicMessage.Object = string(pod2Json)
	topicMessageJson, _ = json.Marshal(topicMessage)
	listwatch.Publish(global.PodRelevantTopic(), string(topicMessageJson))
	time.Sleep(3 * time.Second)
	val, err := utils.Delete(route.Prefix + pod1.GetObjectPath())
	assert.Nil(t, err)
	fmt.Println(val)
	time.Sleep(3 * time.Second)
}
func TestVolumeCreate(t *testing.T) {
	volume1 := apiobjects.Volume{}
	volume1.Name = "volume1"
	var pvcSource apiobjects.PersistentVolumeClaimVolumeSource
	pvcSource.ClaimName = "pvc"
	volume1.PersistentVolumeClaim = &pvcSource
}
