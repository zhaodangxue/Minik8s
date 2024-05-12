package PVcontroller

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"

	"github.com/go-redis/redis/v8"
)

type Controller interface {
	Run()
}
type pvcontroller struct {
}

// 在每一次对PV和PVC进行操作时，都需要检查PV和PVC能否绑定
func (c *pvcontroller) GetAllPVCFromApiserver() (pvcs []*apiobjects.PersistentVolumeClaim) {
	err := utils.GetUnmarshal(route.Prefix+route.PVCPath, &pvcs)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (c *pvcontroller) GetAllPVSWithNamespaceFromApiserver(namespace string) (pvs []*apiobjects.PersistentVolume) {
	err := utils.GetUnmarshal(route.Prefix+route.PVPath+"/"+namespace, &pvs)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (c *pvcontroller) GetPVSWithSpecifiedStorageClassAndSpecifiedCapacity(pvs []*apiobjects.PersistentVolume, storageclassname string, capacity int) (res []*apiobjects.PersistentVolume) {
	for _, pv := range pvs {
		if (pv.Spec.StorageClassName == storageclassname) && (pv.Status == apiobjects.PVAvailable) {
			size_pv, err := utils.GetStorageCapacity(pv.Spec.Capacity.Storage)
			if err != nil {
				fmt.Println(err)
			}
			if size_pv >= capacity {
				res = append(res, pv)
			}
		}
	}
	return
}
func (c *pvcontroller) SelectOnePVToBound(pvs []*apiobjects.PersistentVolume, accessmode string, PVCName string, PVCNamespace string) bool {
	var res bool
	res = false
	for _, pv := range pvs {
		for _, mode := range pv.Spec.AccessModes {
			if mode == accessmode {
				pv.Status = apiobjects.PVBound
				pv.Spec.PVCBinding.PVCname = PVCName
				pv.Spec.PVCBinding.PVCnamespace = PVCNamespace
				url := pv.GetObjectPath() + "/" + pv.Spec.StorageClassName
				utils.PostWithJson(route.Prefix+url, pv)
				res = true
				break
			}
		}
	}
	return res
}
func (c *pvcontroller) CheckPVBoundPVC() {
	pvcs := c.GetAllPVCFromApiserver()
	for _, pvc := range pvcs {
		if pvc.Status == apiobjects.PVCBound {
			continue
		}
		np := pvc.ObjectMeta.Namespace
		size_pvc, err := utils.GetStorageCapacity(pvc.Spec.Resources.Requests.Storage)
		if err != nil {
			fmt.Println(err)
		}
		pvs := c.GetAllPVSWithNamespaceFromApiserver(np)
		pvs_filtered := c.GetPVSWithSpecifiedStorageClassAndSpecifiedCapacity(pvs, pvc.Spec.StorageClassName, size_pvc)
		if len(pvs_filtered) == 0 {
			continue
		}
		if c.SelectOnePVToBound(pvs_filtered, pvc.Spec.AccessModes[0], pvc.ObjectMeta.Name, pvc.ObjectMeta.Namespace) {
			pvc.Status = apiobjects.PVCBound
		}
		url := pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
		utils.PostWithJson(route.Prefix+url, pvc)
	}
}
func (c *pvcontroller) PVUpdateHandler(data string) error {
	return nil
}
func (c *pvcontroller) PVDeleteHandler(data string) error {
	return nil
}
func (c *pvcontroller) PVCreateHandler(data string) error {
	pv := apiobjects.PersistentVolume{}
	err := json.Unmarshal([]byte(data), &pv)
	if err != nil {
		return err
	}
	pv.Status = apiobjects.PVAvailable
	url_pv := pv.GetObjectPath() + "/" + pv.Spec.StorageClassName
	utils.PostWithJson(route.Prefix+url_pv, pv)
	c.CheckPVBoundPVC()
	return nil
}
func (c *pvcontroller) handlePVupdate(msg *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), &topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		err = c.PVCreateHandler(topicMessage.Object)
	case apiobjects.Update:
		err = c.PVUpdateHandler(topicMessage.Object)
	case apiobjects.Delete:
		err = c.PVDeleteHandler(topicMessage.Object)
	}
	if err != nil {
		fmt.Println(err)
	}
}
func (c *pvcontroller) handlePVCupdate(msg *redis.Message) {

}
func (c *pvcontroller) Run() {
	go listwatch.Watch(global.PvRelevantTopic(), c.handlePVupdate)
	listwatch.Watch(global.PvcRelevantTopic(), c.handlePVCupdate)
}
