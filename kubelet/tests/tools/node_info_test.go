package tests

import (
	"minik8s/kubelet/internal"
	"testing"
)

func TestNodeInfo(t *testing.T) {

	nodeStats, err := internal.GetNodeStats()
	if err != nil {
		t.Error(err)
	}
	t.Log(nodeStats)

}
