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
