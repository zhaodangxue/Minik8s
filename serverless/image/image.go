package image

import (
	"minik8s/apiobjects"
	"minik8s/utils"
	"os"
	"os/exec"

	//"strings"

	log "github.com/sirupsen/logrus"
)

// CreateImage to create image for function
func CreateImage(input apiobjects.FunctionCtlInput) (string, error) {
	// 1. create the image
	// 1.1 generate tmp dockerfile for the function from the basic dockerfile
	imageName := "function-" + input.Name
	dstFilePath, err := PrepareBuildEnv(input)
	if err != nil {
		log.Error("[GenerateDockerfile] error")
		return "", err
	}

	// 1.2 create the image
	cmd := exec.Command("docker", "build", "-t", imageName, dstFilePath)
	err = cmd.Run()
	if err != nil {
		log.Error("[CreateImage] create image error: ", err)
		return "", err
	}

	cmd = exec.Command("docker", "tag", imageName, serverIp+":5000/"+imageName+":latest")
	err = cmd.Run()
	if err != nil {
		log.Error("[CreateImage] tag image error: ", err)
		return "", err
	}

	// 2. save the image to the registry
	err = SaveImage(imageName)
	if err != nil {
		log.Error("[CreateImage] save image error: ", err)
		return "", err
	}
	return imageName, nil
}

func PrepareBuildEnv(input apiobjects.FunctionCtlInput) (buildPath string, err error) {
	// 1.1 copy the basic dockerfile to tmp dockerfile for the function
	buildPath = baseDir + "/buildenv"
	// 删除原有的文件
	err = os.RemoveAll(buildPath)
	if err != nil {
		log.Error("[CreateImage] remove old build path error: ", err)
		return
	}
	// 拷贝基础镜像
	err = utils.CopyDir(baseDir+"/imagedata", buildPath)
	if err != nil {
		log.Error("[CreateImage] copy base image error: ", err)
		return
	}
	// 拷贝用户文件
	err = utils.CopyDir(input.BuildOptions.FunctionFileDir, buildPath+input.BuildOptions.FunctionFileDir)
	if err != nil {
		log.Error("[CreateImage] copy user file error: ", err)
		return
	}

	// 打开Dockerfile文件
	dstFile, err := os.OpenFile(buildPath+"/Dockerfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer dstFile.Close()
	if err != nil {
		log.Error("[CreateImage] open tmp docker file error: ", err)
		return
	}

	// 加入用户自定义的命令
	dstFile.WriteString("\n")
	for _, command := range input.BuildOptions.ExtraCommands {
		dstFile.WriteString("\n")
		dstFile.WriteString(command + "\n")
	}

	// 加入默认文件
	dstFile.WriteString("\n")
	copyDir := "COPY " + input.BuildOptions.FunctionFileDir + " /function"
	dstFile.WriteString(copyDir + "\n")

	return
}

// save the image to the registry
func SaveImage(name string) error {
	imageName := serverIp + ":5000/" + name + ":latest"

	//push the image into the registry
	cmd := exec.Command("docker", "push", imageName)
	err := cmd.Run()
	if err != nil {
		log.Error("[SaveImage] push image error: ", err)
		return err
	}
	return nil
}
