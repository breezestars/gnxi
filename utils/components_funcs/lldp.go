package components_funcs

import (
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/tidwall/gjson"
	"os/exec"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"fmt"
)

func InitLldp(device *gostruct.Device) error {

	cmd := "lldpctl -f json"
	json, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}

	fmt.Println(string(json))
	fmt.Println(gjson.Get(string(json), "lldp.interface"))
	intfs := gjson.Get(string(json), "lldp.interface").Array()
	for k,v := range intfs{
		fmt.Println(k,", ",v)
	}

	//inf, err := lldp.NewInterface("eth0")
	//if err != nil {
	//	return err
	//}
	//inf.Name = ygot.String("eth0")
	//
	//for j := 0; j < 32; j++ {
	//	index := "Ethernet" + strconv.Itoa(4*j)
	//
	//	inf, err = lldp.NewInterface(index)
	//	if err != nil {
	//		return err
	//	}
	//	inf.Name = ygot.String(index)
	//
	//}

	return nil
}
