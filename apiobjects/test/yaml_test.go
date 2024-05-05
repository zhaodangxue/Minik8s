package api_test

import (
	"fmt"
	"minik8s/apiobjects"
	ctlutils "minik8s/kubectl/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestPodYaml(t *testing.T) {
	var pod apiobjects.Pod
	content, err := ctlutils.LoadFile("./pod-example.yaml")
	assert.Nil(t, err)
	tp := ctlutils.ParseApiObjectType(content)
	assert.Equal(t, ctlutils.Pod, tp)
	err = yaml.Unmarshal(content, &pod)
	assert.Nil(t, err)
	fmt.Println(pod)
}
func TestNodeYaml(t *testing.T) {
	var node apiobjects.Node
	content, err := ctlutils.LoadFile("./__node.yaml")
	assert.Nil(t, err)
	tp := ctlutils.ParseApiObjectType(content)
	assert.Equal(t, ctlutils.Node, tp)
	err = yaml.Unmarshal(content, &node)
	assert.Nil(t, err)
	fmt.Println(node)
}
func TestPVYaml(t *testing.T) {
	var pv apiobjects.PersistentVolume
	content, err := ctlutils.LoadFile("./persistent_volumn.yaml")
	assert.Nil(t, err)
	tp := ctlutils.ParseApiObjectType(content)
	assert.Equal(t, ctlutils.Pv, tp)
	err = yaml.Unmarshal(content, &pv)
	assert.Nil(t, err)
	fmt.Println(pv)
}
func TestPVCYaml(t *testing.T) {
	var pvc apiobjects.PersistentVolumeClaim
	content, err := ctlutils.LoadFile("./persistent_volumn_claim.yaml")
	assert.Nil(t, err)
	tp := ctlutils.ParseApiObjectType(content)
	assert.Equal(t, ctlutils.Pvc, tp)
	err = yaml.Unmarshal(content, &pvc)
	assert.Nil(t, err)
	fmt.Println(pvc)
}
