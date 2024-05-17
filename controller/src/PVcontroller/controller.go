package PVcontroller

import (
	"encoding/json"
	"fmt"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/route"
	"minik8s/controller/api"
	"minik8s/global"
	"minik8s/listwatch"
	"minik8s/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

type Controller interface {
	Run()
}
type PVcontroller struct {
	initInfo          api.InitStruct
	ListFuncEnvelops  []api.ListFuncEnvelop
	WatchFuncEnvelops []api.WatchFuncEnvelop
}

func (c *PVcontroller) Init(init api.InitStruct) {
	c.initInfo = init
	c.ListFuncEnvelops = make([]api.ListFuncEnvelop, 0)
	c.WatchFuncEnvelops = make([]api.WatchFuncEnvelop, 0)
	c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
		Func:  c.handlePVupdate_API,
		Topic: global.PvRelevantTopic(),
	})
	c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
		Func:  c.handlePVCupdate_API,
		Topic: global.PvcRelevantTopic(),
	})
	c.WatchFuncEnvelops = append(c.WatchFuncEnvelops, api.WatchFuncEnvelop{
		Func:  c.handlePodUpdate_API,
		Topic: global.PodRelevantTopic(),
	})
}
func (c *PVcontroller) GetListFuncEnvelops() []api.ListFuncEnvelop {
	return c.ListFuncEnvelops
}
func (c *PVcontroller) GetWatchFuncEnvelops() []api.WatchFuncEnvelop {
	return c.WatchFuncEnvelops
}

// 在每一次对PV和PVC进行操作时，都需要检查PV和PVC能否绑定
func (c *PVcontroller) GetAllPVCFromApiserver() (pvcs []*apiobjects.PersistentVolumeClaim) {
	err := utils.GetUnmarshal(route.Prefix+route.PVCPath, &pvcs)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (c *PVcontroller) GetAllPVSWithNamespaceFromApiserver(namespace string) (pvs []*apiobjects.PersistentVolume) {
	err := utils.GetUnmarshal(route.Prefix+route.PVPath+"/"+namespace, &pvs)
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (c *PVcontroller) GetPVSWithSpecifiedStorageClassAndSpecifiedCapacity(pvs []*apiobjects.PersistentVolume, storageclassname string, capacity int) (res []*apiobjects.PersistentVolume) {
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
func (c *PVcontroller) SelectOnePVToBound(pvs []*apiobjects.PersistentVolume, accessmode string, pvc *apiobjects.PersistentVolumeClaim) bool {
	var res bool
	res = false
	for _, pv := range pvs {
		for _, mode := range pv.Spec.AccessModes {
			if mode == accessmode {
				pv.Status = apiobjects.PVBound
				pv.Spec.PVCBinding.PVCname = pvc.ObjectMeta.Name
				pv.Spec.PVCBinding.PVCnamespace = pvc.ObjectMeta.Namespace
				pvc.PVBinding.PVname = pv.ObjectMeta.Name
				pvc.PVBinding.PVnamespace = pv.ObjectMeta.Namespace
				pvc.PVBinding.PVcapacity = pv.Spec.Capacity.Storage
				pvc.PVBinding.PVpath = pv.Spec.NFS.Path
				url := pv.GetObjectPath() + "/" + pv.Spec.StorageClassName
				utils.PostWithJson(route.Prefix+url, pv)
				res = true
				break
			}
		}
	}
	return res
}
func (c *PVcontroller) DynamicAllocatePV(pvc *apiobjects.PersistentVolumeClaim) error {
	pv := apiobjects.PersistentVolume{}
	pv.ObjectMeta.Namespace = pvc.ObjectMeta.Namespace
	pv.ObjectMeta.UID = utils.NewUUID()
	pv.ObjectMeta.Name = "pv-" + pv.ObjectMeta.UID
	pv.CreationTimestamp = time.Now()
	pv.Spec.StorageClassName = pvc.Spec.StorageClassName
	pv.Spec.Capacity.Storage = pvc.Spec.Resources.Requests.Storage
	pv.Spec.AccessModes = pvc.Spec.AccessModes
	pv.Status = apiobjects.PVBound
	pv.Spec.VolumeMode = "Filesystem"
	pv.Spec.PersistentVolumeReclaimPolicy = "Retain"
	pv.TypeMeta.ApiVersion = "v1"
	pv.TypeMeta.Kind = "PersistentVolume"
	pv.Spec.PVCBinding.PVCname = pvc.ObjectMeta.Name
	pv.Spec.PVCBinding.PVCnamespace = pvc.ObjectMeta.Namespace
	pv.Dynamic = 1 // 1表示动态分配
	url_pv := route.PVDynamicAllocate
	_, err := utils.PostWithJson(route.Prefix+url_pv, pv)
	if err != nil {
		return err
	}
	pvc.PVBinding.PVname = pv.ObjectMeta.Name
	pvc.PVBinding.PVnamespace = pv.ObjectMeta.Namespace
	pvc.PVBinding.PVcapacity = pv.Spec.Capacity.Storage
	pvc.PVBinding.PVpath = global.NFSdir + "/" + pv.ObjectMeta.Name
	return nil
}
func (c *PVcontroller) CheckPVBoundPVC() {
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
		if c.SelectOnePVToBound(pvs_filtered, pvc.Spec.AccessModes[0], pvc) {
			pvc.Status = apiobjects.PVCBound
		} else {
			err = c.DynamicAllocatePV(pvc)
			if err != nil {
				fmt.Println(err)
			}
			pvc.Status = apiobjects.PVCBound
		}
		url := pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
		utils.PostWithJson(route.Prefix+url, pvc)
	}
}
func (c *PVcontroller) CancelPVBoundPVC(pvc *apiobjects.PersistentVolumeClaim) error {
	pv := apiobjects.PersistentVolume{}
	url := route.Prefix + route.PVPath + "/" + pvc.PVBinding.PVnamespace + "/" + pvc.PVBinding.PVname
	err := utils.GetUnmarshal(url, &pv)
	if err != nil {
		return err
	}
	if pv.Name == "" {
		return fmt.Errorf("pv not found")
	}
	pv.Status = apiobjects.PVAvailable
	pv.Spec.PVCBinding.PVCname = ""
	pv.Spec.PVCBinding.PVCnamespace = ""
	url = pv.GetObjectPath() + "/" + pv.Spec.StorageClassName
	_, err = utils.PostWithJson(route.Prefix+url, pv)
	if err != nil {
		return err
	}
	return nil
}
func (c *PVcontroller) PVUpdateHandler(data string) error {
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
func (c *PVcontroller) PVDeleteHandler(data string) error {
	pv := apiobjects.PersistentVolume{}
	err := json.Unmarshal([]byte(data), &pv)
	if err != nil {
		return err
	}
	fmt.Printf("delete pv name: %s namespace: %s uuid: %s", pv.ObjectMeta.Name, pv.ObjectMeta.Namespace, pv.ObjectMeta.UID)
	return nil
}
func (c *PVcontroller) PVCreateHandler(data string) error {
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
func (c *PVcontroller) handlePVupdate(msg *redis.Message) {
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
func (c *PVcontroller) handlePVupdate_API(controller api.Controller, message apiobjects.TopicMessage) (err error) {
	switch message.ActionType {
	case apiobjects.Create:
		err = c.PVCreateHandler(message.Object)
	case apiobjects.Update:
		err = c.PVUpdateHandler(message.Object)
	case apiobjects.Delete:
		err = c.PVDeleteHandler(message.Object)
	}
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (c *PVcontroller) PVCCreateHandler(data string) error {
	pvc := apiobjects.PersistentVolumeClaim{}
	err := json.Unmarshal([]byte(data), &pvc)
	if err != nil {
		return err
	}
	pvc.Status = apiobjects.PVCAvailable
	url_pvc := pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
	utils.PostWithJson(route.Prefix+url_pvc, pvc)
	c.CheckPVBoundPVC()
	return nil
}
func (c *PVcontroller) PVCDeleteHandler(data string) error {
	pvc := apiobjects.PersistentVolumeClaim{}
	err := json.Unmarshal([]byte(data), &pvc)
	if err != nil {
		return err
	}
	if pvc.Status == apiobjects.PVCBound {
		err = c.CancelPVBoundPVC(&pvc)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *PVcontroller) PVCUpdateHandler(data string) error {
	pvc := apiobjects.PersistentVolumeClaim{}
	err := json.Unmarshal([]byte(data), &pvc)
	if err != nil {
		return err
	}
	pvc.Status = apiobjects.PVCAvailable
	url_pvc := pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
	utils.PostWithJson(route.Prefix+url_pvc, pvc)
	c.CheckPVBoundPVC()
	return nil
}
func (c *PVcontroller) handlePVCupdate(msg *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), &topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		err = c.PVCCreateHandler(topicMessage.Object)
	case apiobjects.Update:
		err = c.PVCUpdateHandler(topicMessage.Object)
	case apiobjects.Delete:
		err = c.PVCDeleteHandler(topicMessage.Object)
	}
	if err != nil {
		fmt.Println(err)
	}
}
func (c *PVcontroller) handlePVCupdate_API(controller api.Controller, message apiobjects.TopicMessage) (err error) {
	switch message.ActionType {
	case apiobjects.Create:
		err = c.PVCCreateHandler(message.Object)
	case apiobjects.Update:
		err = c.PVCUpdateHandler(message.Object)
	case apiobjects.Delete:
		err = c.PVCDeleteHandler(message.Object)
	}
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (c *PVcontroller) PodCreateHandler(data string) error {
	pod := apiobjects.Pod{}
	err := json.Unmarshal([]byte(data), &pod)
	if err != nil {
		return err
	}
	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		PVCname := volume.PersistentVolumeClaim.ClaimName
		PVCnamespace := volume.PersistentVolumeClaim.ClaimNamespace
		url := route.Prefix + route.PVCPath + "/" + PVCnamespace + "/" + PVCname
		pvc := &apiobjects.PersistentVolumeClaim{}
		err = utils.GetUnmarshal(url, pvc)
		if err != nil {
			return err
		}
		if pvc.ObjectMeta.Name == "" {
			return fmt.Errorf("pvc not found")
		}
		if pvc.Status != apiobjects.PVCBound {
			return fmt.Errorf("pvc not bound")
		}
		var PodAbstract apiobjects.PodAbstract
		PodAbstract.Podname = pod.ObjectMeta.Name
		PodAbstract.Namespace = pod.ObjectMeta.Namespace
		var PodBinding []apiobjects.PodAbstract
		PodBinding = pvc.PodBinding
		PodBinding = append(PodBinding, PodAbstract)
		pvc.PodBinding = PodBinding
		url = pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
		_, err = utils.PostWithJson(route.Prefix+url, pvc)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *PVcontroller) PodUpdateHandler(data string) error {
	pod := apiobjects.Pod{}
	err := json.Unmarshal([]byte(data), &pod)
	if err != nil {
		return err
	}
	pvcs := c.GetAllPVCFromApiserver()
	for _, pvc := range pvcs {
		var flag bool
		var PodAbstracts []apiobjects.PodAbstract
		flag = false
		if pvc.Status == apiobjects.PVCAvailable {
			continue
		}
		if len(pvc.PodBinding) == 0 {
			continue
		}
		for _, podBinding := range pvc.PodBinding {
			if podBinding.Podname == pod.ObjectMeta.Name && podBinding.Namespace == pod.ObjectMeta.Namespace {
				flag = true
			}
		}
		if flag == true {
			for _, podBinding := range pvc.PodBinding {
				if podBinding.Podname != pod.ObjectMeta.Name || podBinding.Namespace != pod.ObjectMeta.Namespace {
					PodAbstracts = append(PodAbstracts, podBinding)
				}
			}
		}
		pvc.PodBinding = PodAbstracts
		_, err = utils.PostWithJson(route.Prefix+pvc.GetObjectPath()+"/"+pvc.Spec.StorageClassName, pvc)
		if err != nil {
			return err
		}
	}
	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		PVCname := volume.PersistentVolumeClaim.ClaimName
		PVCnamespace := volume.PersistentVolumeClaim.ClaimNamespace
		url := route.Prefix + route.PVCPath + "/" + PVCnamespace + "/" + PVCname
		pvc := &apiobjects.PersistentVolumeClaim{}
		err = utils.GetUnmarshal(url, pvc)
		if err != nil {
			return err
		}
		if pvc.Name == "" {
			return fmt.Errorf("pvc not found")
		}
		if pvc.Status != apiobjects.PVCBound {
			return fmt.Errorf("pvc not bound")
		}
		var PodAbstract apiobjects.PodAbstract
		PodAbstract.Podname = pod.ObjectMeta.Name
		PodAbstract.Namespace = pod.ObjectMeta.Namespace
		var PodBinding []apiobjects.PodAbstract
		PodBinding = pvc.PodBinding
		PodBinding = append(PodBinding, PodAbstract)
		pvc.PodBinding = PodBinding
		url = pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
		_, err = utils.PostWithJson(url, pvc)
		if err != nil {
			return err
		}

	}
	return nil
}
func (c *PVcontroller) PodDeleteHandler(data string) error {
	pod := apiobjects.Pod{}
	err := json.Unmarshal([]byte(data), &pod)
	if err != nil {
		return err
	}
	for _, volume := range pod.Spec.Volumes {
		if volume.PersistentVolumeClaim == nil {
			continue
		}
		PVCname := volume.PersistentVolumeClaim.ClaimName
		PVCnamespace := volume.PersistentVolumeClaim.ClaimNamespace
		url := route.Prefix + route.PVCPath + "/" + PVCnamespace + "/" + PVCname
		pvc := &apiobjects.PersistentVolumeClaim{}
		err = utils.GetUnmarshal(url, pvc)
		if err != nil {
			return err
		}
		if pvc.Name == "" {
			return fmt.Errorf("pvc not found")
		}
		if pvc.Status != apiobjects.PVCBound {
			return fmt.Errorf("pvc not bound")
		}
		var PodAbstracts []apiobjects.PodAbstract
		for _, podBinding := range pvc.PodBinding {
			if podBinding.Podname != pod.ObjectMeta.Name || podBinding.Namespace != pod.ObjectMeta.Namespace {
				PodAbstracts = append(PodAbstracts, podBinding)
			}
		}
		pvc.PodBinding = PodAbstracts
		url = pvc.GetObjectPath() + "/" + pvc.Spec.StorageClassName
		_, err = utils.PostWithJson(route.Prefix+url, pvc)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *PVcontroller) handlePodUpdate(msg *redis.Message) {
	topicMessage := apiobjects.TopicMessage{}
	err := json.Unmarshal([]byte(msg.Payload), &topicMessage)
	if err != nil {
		fmt.Println(err)
	}
	switch topicMessage.ActionType {
	case apiobjects.Create:
		err = c.PodCreateHandler(topicMessage.Object)
	case apiobjects.Update:
		err = c.PodUpdateHandler(topicMessage.Object)
	case apiobjects.Delete:
		err = c.PodDeleteHandler(topicMessage.Object)
	}
	if err != nil {
		fmt.Println(err)
	}
}
func (c *PVcontroller) handlePodUpdate_API(controller api.Controller, message apiobjects.TopicMessage) (err error) {
	switch message.ActionType {
	case apiobjects.Create:
		err = c.PodCreateHandler(message.Object)
	case apiobjects.Update:
		err = c.PodUpdateHandler(message.Object)
	case apiobjects.Delete:
		err = c.PodDeleteHandler(message.Object)
	}
	if err != nil {
		fmt.Println(err)
	}
	return
}
func (c *PVcontroller) Run() {
	go listwatch.Watch(global.PvRelevantTopic(), c.handlePVupdate)
	go listwatch.Watch(global.PvcRelevantTopic(), c.handlePVCupdate)
	listwatch.Watch(global.PodRelevantTopic(), c.handlePodUpdate)
}
func New() Controller {
	return &PVcontroller{}
}
