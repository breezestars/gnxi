package utils

import (
	"github.com/openconfig/ygot/ygot"
	"strconv"
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"fmt"
	"encoding/json"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"time"
	"strings"
)

func initLldp(device *gostruct.Device) error {

	lldp := &gostruct.Lldp{}
	inf, err := lldp.NewInterface("eth0")
	if err != nil {
		return err
	}
	inf.Name = ygot.String("eth0")

	for j := 0; j < 32; j++ {
		index := "Ethernet" + strconv.Itoa(4*j)

		inf, err = lldp.NewInterface(index)
		if err != nil {
			return err
		}
		inf.Name = ygot.String(index)

	}
	device.Lldp = lldp
	return nil
}

func initInterface(device *gostruct.Device) error {

	//cmd := exec.Command("/bin/sh", "-c", `show interfaces status | grep Ethernet | awk -F' ' '{print $1" "$2" "$3" "$4" "$5" "$6" "$7}'`)
	//cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	//cmd.Start()
	//cmd.Run()
	//cmd.Wait()

	console := `Ethernet0 49,50,51,52 N/A 9100 hundredGigE1 down up
Ethernet4 53,54,55,56 N/A 9100 hundredGigE2 down up
Ethernet8 57,58,59,60 N/A 9100 hundredGigE3 down up
Ethernet12 61,62,63,64 N/A 9100 hundredGigE4 down up
Ethernet16 65,66,67,68 N/A 9100 hundredGigE5 down up
Ethernet20 69,70,71,72 N/A 9100 hundredGigE6 down up
Ethernet24 73,74,75,76 N/A 9100 hundredGigE7 down up
Ethernet28 77,78,79,80 N/A 9100 hundredGigE8 down up
Ethernet32 33,34,35,36 N/A 9100 hundredGigE9 down up
Ethernet36 37,38,39,40 N/A 9100 hundredGigE10 down up
Ethernet40 41,42,43,44 N/A 9100 hundredGigE11 down up
Ethernet44 45,46,47,48 N/A 9100 hundredGigE12 down up
Ethernet48 81,82,83,84 N/A 9100 hundredGigE13 down up
Ethernet52 85,86,87,88 N/A 9100 hundredGigE14 down up
Ethernet56 89,90,91,92 N/A 9100 hundredGigE15 down up
Ethernet60 93,94,95,96 N/A 9100 hundredGigE16 down up
Ethernet64 97,98,99,100 N/A 9100 hundredGigE17 down up
Ethernet68 101,102,103,104 N/A 9100 hundredGigE18 down up
Ethernet72 105,106,107,108 N/A 9100 hundredGigE19 down up
Ethernet76 109,110,111,112 N/A 9100 hundredGigE20 down up
Ethernet80 17,18,19,20 N/A 9100 hundredGigE21 down up
Ethernet84 21,22,23,24 N/A 9100 hundredGigE22 down up
Ethernet88 25,26,27,28 N/A 9100 hundredGigE23 down up
Ethernet92 29,30,31,32 N/A 9100 hundredGigE24 down up
Ethernet96 113,114,115,116 N/A 9100 hundredGigE25 down up
Ethernet100 117,118,119,120 N/A 9100 hundredGigE26 down up
Ethernet104 121,122,123,124 N/A 9100 hundredGigE27 down up
Ethernet108 125,126,127,128 N/A 9100 hundredGigE28 down up
Ethernet112 1,2,3,4 N/A 9100 hundredGigE29 down up
Ethernet116 5,6,7,8 N/A 9100 hundredGigE30 down up
Ethernet120 9,10,11,12 N/A 9100 hundredGigE31 down up
Ethernet124 13,14,15,16 N/A 9100 hundredGigE32 down up`

	// 0:Interface
	// 1:Lanes
	// 2:Speed
	// 3:MTU
	// 4:Alias
	// 5:Oper
	// 6:Admin
	consoleArray := strings.Split(console, "\n")

	inf, err := device.NewInterface("eth0")
	if err != nil {
		return err
	}

	inf.Name = ygot.String("eth0")
	inf.Enabled = ygot.Bool(true)
	inf.AdminStatus = gostruct.OpenconfigInterfaces_Interface_AdminStatus_UP
	inf.OperStatus = gostruct.OpenconfigInterfaces_Interface_OperStatus_UP

	for j := 0; j < len(consoleArray); j++ {
		intDetail := strings.Split(consoleArray[j], " ")

		intName := strings.Split(consoleArray[j], " ")[0]

		mtu, err := strconv.Atoi(intDetail[3])
		if err != nil {
			return err
		}

		inf, err := device.NewInterface(intName)
		inf.Name = ygot.String(intName)
		inf.Mtu = ygot.Uint16(uint16(mtu))

		if intDetail[6] == "up" {
			inf.Enabled = ygot.Bool(true)
			inf.AdminStatus = gostruct.OpenconfigInterfaces_Interface_AdminStatus_UP
		} else {
			inf.Enabled = ygot.Bool(false)
			inf.AdminStatus = gostruct.OpenconfigInterfaces_Interface_AdminStatus_DOWN
		}

		if intDetail[5] == "up" {
			inf.OperStatus = gostruct.OpenconfigInterfaces_Interface_OperStatus_UP
		} else {
			inf.OperStatus = gostruct.OpenconfigInterfaces_Interface_OperStatus_DOWN
		}

	}
	return nil
}

func initInterfaceEthernet() {

}

func initInterfaceAggregate() {

}

func initInterfaceIp() {

}

func initVlan(device *gostruct.Device) error {

	// TODO: Should gather information from switch, here is still use static data

	// I use the Switch Serial Number to this forwarding instance unique name
	// The Serial Number is come from following command
	// show platform syseeprom | grep 0x23 | awk -F' ' '{print $5}'
	int, err := device.NewNetworkInstance("771232X1721087")
	if err != nil {
		return err
	}

	vlan1, err := int.NewVlan(uint16(1))
	if err != nil {
		return err
	}
	vlan1.Name = ygot.String("VLAN1")

	for idx := 0; idx < len(device.Interface)-1; idx++ {
		name := device.GetInterface("Ethernet" + strconv.Itoa(idx*4)).Name
		vlan1.Member = append(vlan1.Member, &gostruct.NetworkInstance_Vlan_Member{
			Interface: name,
		})
	}
	//
	//vlan1.Member = make([]*gostruct.NetworkInstance_Vlan_Member, len(device.Interface))
	//for idx := range vlan1.Member {
	//	//name:=device.GetInterface("Ethernet" + strconv.Itoa(idx*4)).Name
	//	vlan1.Member[idx].Interface=ygot.String("Test")
	//}

	return nil
}

func InitGoStruct() (ygot.ValidatedGoStruct, error) {
	device := &gostruct.Device{}

	err := initLldp(device)
	if err != nil {
		return nil, err
	}

	err = initInterface(device)
	if err != nil {
		return nil, err
	}

	err = initVlan(device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func main() {
	notifications := make([]*pb.Notification, 1)

	device := &gostruct.Device{}

	err := initLldp(device)
	if err != nil {
		print(err.Error())
	}

	err = initInterface(device)
	if err != nil {
		print(err.Error())
	}

	err = initVlan(device)
	if err != nil {
		print(err.Error())
	}

	testJsonTree, err := ygot.ConstructIETFJSON(device, &ygot.RFC7951JSONConfig{AppendModuleName: true})
	if err != nil {
		panic(fmt.Sprintf("JSON demo error: %v", err))
	}

	jsonDump, err := json.Marshal(testJsonTree)
	fmt.Println(jsonDump)
	if err != nil {
		fmt.Sprintf("error in marshaling IETF JSON tree to bytes: %v", err)
	}

	update := &pb.Update{
		Path: &pb.Path{},
		Val: &pb.TypedValue{
			Value: &pb.TypedValue_JsonIetfVal{
				JsonIetfVal: jsonDump,
			},
		},
	}
	ts := time.Now().UnixNano()
	notifications[0] = &pb.Notification{
		Timestamp: ts,
		Prefix:    &pb.Path{},
		Update:    []*pb.Update{update},
	}
	PrintProto(&pb.GetResponse{Notification: notifications})

}
