package ipvs

import (
	"fmt"
	"os/exec"

	"github.com/mqliang/libipvs"
	//"github.com/coreos/go-iptables/iptables"
)

var handle *libipvs.IPVSHandle

func init() {
	handle ,err := libipvs.New()
	if handle == nil {
		fmt.Println(err.Error())
	}

	_, err = exec.Command("sysctl", []string{"net.ipv4.vs.conntrack=1"}...).CombinedOutput()

	if err != nil {
		fmt.Println(err.Error())
	}
	

}



