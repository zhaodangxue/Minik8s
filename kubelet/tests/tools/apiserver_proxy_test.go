package tests

import (
	"minik8s/kubelet/internal"
	"testing"
)

func TestGetAllBindings(t *testing.T) {
	bindings, err := internal.GetAllBindings()
	if err != nil {
		t.Error(err)
	}
	t.Log(bindings)
}
