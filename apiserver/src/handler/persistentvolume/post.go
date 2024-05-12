package persistentvolume_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/utils"

	"github.com/gin-gonic/gin"
)

func PVApplyDetailHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	pvName := c.Param("name")
	storageclass := c.Param("storageclass")
	url := "/api/persistentvolume" + "/" + namespace + "/" + pvName + "/" + storageclass
	pv := apiobjects.PersistentVolume{}
	err := utils.ReadUnmarshal(c.Request.Body, &pv)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	pvJson, _ := json.Marshal(pv)
	etcd.Put(url, string(pvJson))
	c.String(200, "ok")
}
func PVCApplyDetailHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	pvcName := c.Param("name")
	storageclass := c.Param("storageclass")
	url := "/api/persistentvolumeclaim" + "/" + namespace + "/" + pvcName + "/" + storageclass
	pvc := apiobjects.PersistentVolumeClaim{}
	err := utils.ReadUnmarshal(c.Request.Body, &pvc)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	pvcJson, _ := json.Marshal(pvc)
	etcd.Put(url, string(pvcJson))
	c.String(200, "ok")
}
