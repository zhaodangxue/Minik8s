package image

import (
	"io"
	"minik8s/apiobjects"
	"os"
	"os/exec"

	//"strings"

	log "github.com/sirupsen/logrus"
)

const serverIp = "192.168.1.15"

// CreateImage to create image for function
func CreateImage(input apiobjects.FunctionCtlInput) (string, error) {
	// 1. create the image
	// 1.1 generate tmp dockerfile for the function from the basic dockerfile
	imageName := "Function-" + input.Name
	err := GenerateDockerfile(input)
	if err != nil {
		log.Error("[GenerateDockerfile] error")
		return "", err
	}

	// 1.2 create the image
	cmd := exec.Command("docker", "build", "-t", imageName, "/home/tmpdata/Dockerfile")
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

func GenerateDockerfile(input apiobjects.FunctionCtlInput) error {
	// 1.1 copy the basic dockerfile to tmp dockerfile for the function
	srcFile, err := os.Open("/home/imagedata/Dockerfile")
	if err != nil {
		log.Error("[CreateImage] open basic docker file error: ", err)
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile("/home/tmpdata/Dockerfile", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Error("[CreateImage] open tmp docker file error: ", err)
		return err
	}
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Error("[CreateImage] copy file error: ", err)
		return err
	}

	//add the extra commands
	dstFile.WriteString("\n")
	for _, command := range input.BuildOptions.ExtraCommands {
		dstFile.WriteString("\n")
		dstFile.WriteString(command + "\n")
	}
	dstFile.WriteString("\n")
	copyDir := "COPY " + input.BuildOptions.FunctionFileDir + " /function"
	dstFile.WriteString(copyDir + "\n")

	defer dstFile.Close()
	return nil
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
