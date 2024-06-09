package kubectl_test

import (
	command "minik8s/kubectl/src"
	testing "testing"
)

func TestGetAll(t *testing.T) {
	_, err := command.GetAllJobsInCluster()
	if err != nil {
		t.Error(err)
	}
}
