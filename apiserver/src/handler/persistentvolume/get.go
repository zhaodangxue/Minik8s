package persistentvolume_handler

import (
	"encoding/json"
	"minik8s/apiobjects"
	"minik8s/apiserver/src/etcd"
	"minik8s/apiserver/src/route"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PVCSGetHandler(c *gin.Context) {
	var pvcs []*apiobjects.PersistentVolumeClaim
	values, err := etcd.Get_prefix(route.PVCPath)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	for _, value := range values {
		var pvc apiobjects.PersistentVolumeClaim
		err := json.Unmarshal([]byte(value), &pvc)
		if err != nil {
			c.String(200, err.Error())
			return
		}
		pvcs = append(pvcs, &pvc)
	}
	c.JSON(http.StatusOK, pvcs)
}
func PVGetSpecifiedHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	pvName := c.Param("name")
	url := route.PVPath + "/" + namespace + "/" + pvName
	val, _ := etcd.Get_prefix(url)
	var pv apiobjects.PersistentVolume
	if len(val) == 0 {
		c.JSON(http.StatusOK, pv)
		return
	}
	err := json.Unmarshal([]byte(val[0]), &pv)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSON(http.StatusOK, pv)
}
func PVCGetSpecifiedHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	pvcName := c.Param("name")
	url := route.PVCPath + "/" + namespace + "/" + pvcName
	val, _ := etcd.Get_prefix(url)
	var pvc apiobjects.PersistentVolumeClaim
	if len(val) == 0 {
		c.JSON(http.StatusOK, pvc)
		return
	}
	err := json.Unmarshal([]byte(val[0]), &pvc)
	if err != nil {
		c.String(200, err.Error())
		return
	}
	c.JSON(http.StatusOK, pvc)
}
