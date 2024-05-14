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
	etcd.Clear()
	err = etcd.Put("test1", "test-1")
	assert.Nil(t, err)
	err = etcd.Put("test2", "test-2")
	assert.Nil(t, err)
	err = etcd.Put("test3", "test-3")
	assert.Nil(t, err)
	str, err := etcd.Get("test1")
	assert.Nil(t, err)
	assert.Equal(t, "test-1", str)
	err = etcd.Delete("test1")
	assert.Nil(t, err)
	str, err = etcd.Get("test1")
	assert.NotEqual(t, "test-1", str)
	assert.Nil(t, err)
	err = etcd.Clear()
	assert.Nil(t, err)
	str, err = etcd.Get("test2")
	assert.NotEqual(t, "test-2", str)
	assert.Nil(t, err)
}
