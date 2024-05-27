package prometheus

import (
	"testing"
)

func TestPrometheusEditor(t *testing.T) {
    addressList := []string{}
	addressList = append(addressList, "192.168.1.14:9100")
	addressList = append(addressList, "192.168.1.15:9100")
	addressList = append(addressList, "10.0.0.1:2112")

	GenerateProConfig(addressList)
}
