package test

import (
	"fmt"
	"minik8s/apiserver/src/etcd"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtcd(t *testing.T) {
	fmt.Println("TestEtcd")
	var err error
	err = etcd.Put("test1", "test-1")
	assert.Nil(t, err)
}
