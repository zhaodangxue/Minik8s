package test_utils

import (
	"minik8s/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetStorageCapacity(t *testing.T) {
	size, err := utils.GetStorageCapacity("10Gi")
	assert.Nil(t, err)
	assert.Equal(t, 10*1024*1024*1024, size)
	size, err = utils.GetStorageCapacity("4Mi")
	assert.Nil(t, err)
	assert.Equal(t, 4*1024*1024, size)
}
