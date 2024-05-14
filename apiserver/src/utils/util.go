package apiserver_utils

import (
	"minik8s/global"
	"os"
)

func GeneratePVPath(name string) (string, error) {
	//这里接下来要做的是在/home/k8s-nfs目录下创建一个名为name的文件夹
	//然后返回这个文件夹的路径
	err := os.Mkdir(global.ApiserverMountDir+"/"+name, 0755)
	if err != nil {
		return "", err
	}
	return global.NFSdir + "/" + name, nil
}
func DeletePVPath(name string) error {
	err := os.RemoveAll(global.ApiserverMountDir + "/" + name)
	return err
}
