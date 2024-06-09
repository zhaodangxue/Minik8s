package kubectl_test

import (
	ctlutils "minik8s/kubectl/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLetter(t *testing.T) {
	assert.Equal(t, true, ctlutils.IsLetter('a'))
	assert.Equal(t, true, ctlutils.IsLetter('z'))
	assert.Equal(t, true, ctlutils.IsLetter('A'))
	assert.Equal(t, true, ctlutils.IsLetter('Z'))
	assert.Equal(t, false, ctlutils.IsLetter('1'))
	assert.Equal(t, false, ctlutils.IsLetter(' '))
}
func TestParseApiObjectType(t *testing.T) {
	assert.Equal(t, ctlutils.Test, ctlutils.ParseApiObjectType([]byte("kind: test")))
	assert.Equal(t, ctlutils.Pod, ctlutils.ParseApiObjectType([]byte("kind: pod")))
	assert.Equal(t, ctlutils.Node, ctlutils.ParseApiObjectType([]byte("kind: node")))
	assert.Equal(t, ctlutils.Unknown, ctlutils.ParseApiObjectType([]byte("kind: unknown")))
}
func TestLoadFile(t *testing.T) {
	_, err := ctlutils.LoadFile("./test.yaml")
	assert.Nil(t, err)
}
func TestParseType(t *testing.T) {
	data, err := ctlutils.ParseType("./test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, ctlutils.Test, data)
}
