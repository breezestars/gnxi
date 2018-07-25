package utils

import (
	"github.com/openconfig/ygot/ygot"
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"fmt"
	"encoding/json"
	pb "github.com/openconfig/gnmi/proto/gnmi"
	"time"
	"strings"
	"strconv"
	"os/exec"
)

func initLldp(device *gostruct.Device) error {

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

func initInterface(device *gostruct.Device) error {

	//cmd := exec.Command("/bin/sh", "-c", `show interfaces status | grep Ethernet | awk -F' ' '{print $1" "$2" "$3" "$4" "$5" "$6" "$7}'`)
	//cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	//cmd.Start()
	//cmd.Run()
	//cmd.Wait()

	t0 := time.Now()
	cmd:="show interfaces status | grep Ethernet | awk -F' ' '{print $1\" \"$2\" \"$3\" \"$4\" \"$5\" \"$6\" \"$7}'"
	intfStatus, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmd)
	}
	t0d := time.Since(t0)

//	intfStatus := `Ethernet0 49,50,51,52 N/A 9100 hundredGigE1 down up
//Ethernet4 53,54,55,56 N/A 9100 hundredGigE2 down up
//Ethernet8 57,58,59,60 N/A 9100 hundredGigE3 down up
//Ethernet12 61,62,63,64 N/A 9100 hundredGigE4 down up
//Ethernet16 65,66,67,68 N/A 9100 hundredGigE5 down up
//Ethernet20 69,70,71,72 N/A 9100 hundredGigE6 down up
//Ethernet24 73,74,75,76 N/A 9100 hundredGigE7 down up
//Ethernet28 77,78,79,80 N/A 9100 hundredGigE8 down up
//Ethernet32 33,34,35,36 N/A 9100 hundredGigE9 down up
//Ethernet36 37,38,39,40 N/A 9100 hundredGigE10 down up
//Ethernet40 41,42,43,44 N/A 9100 hundredGigE11 down up
//Ethernet44 45,46,47,48 N/A 9100 hundredGigE12 down up
//Ethernet48 81,82,83,84 N/A 9100 hundredGigE13 down up
//Ethernet52 85,86,87,88 N/A 9100 hundredGigE14 down up
//Ethernet56 89,90,91,92 N/A 9100 hundredGigE15 down up
//Ethernet60 93,94,95,96 N/A 9100 hundredGigE16 down up
//Ethernet64 97,98,99,100 N/A 9100 hundredGigE17 down up
//Ethernet68 101,102,103,104 N/A 9100 hundredGigE18 down up
//Ethernet72 105,106,107,108 N/A 9100 hundredGigE19 down up
//Ethernet76 109,110,111,112 N/A 9100 hundredGigE20 down up
//Ethernet80 17,18,19,20 N/A 9100 hundredGigE21 down up
//Ethernet84 21,22,23,24 N/A 9100 hundredGigE22 down up
//Ethernet88 25,26,27,28 N/A 9100 hundredGigE23 down up
//Ethernet92 29,30,31,32 N/A 9100 hundredGigE24 down up
//Ethernet96 113,114,115,116 N/A 9100 hundredGigE25 down up
//Ethernet100 117,118,119,120 N/A 9100 hundredGigE26 down up
//Ethernet104 121,122,123,124 N/A 9100 hundredGigE27 down up
//Ethernet108 125,126,127,128 N/A 9100 hundredGigE28 down up
//Ethernet112 1,2,3,4 N/A 9100 hundredGigE29 down up
//Ethernet116 5,6,7,8 N/A 9100 hundredGigE30 down up
//Ethernet120 9,10,11,12 N/A 9100 hundredGigE31 down up
//Ethernet124 13,14,15,16 N/A 9100 hundredGigE32 down up`

	// 0:Interface
	// 1:Lanes
	// 2:Speed
	// 3:MTU
	// 4:Alias
	// 5:Oper
	// 6:Admin

	t1 := time.Now()

	portstatCmd:="portstat | grep 'Ethernet' | awk -F' ' '{print $1,$3,$6,$7,$9,$12,$13}'"
	portstat, err := exec.Command("bash","-c",portstatCmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmd)
	}
	t1d := time.Since(t1)
//	portstat := `Ethernet0 1595169299934 1824671749079 0 3739096352458 837308284053 0
//Ethernet4 2708494058496 1903492724623 0 2449799110453 787421450391 0
//Ethernet8 2199032828538 1860401179957 0 1883327787315 717481680718 0
//Ethernet12 3123175735980 2042899780017 0 2418081566910 914335203010 0
//Ethernet16 2592753524744 2111814837197 0 2567124117141 957743975636 0
//Ethernet20 2359947707839 2519429426229 0 2408779906784 968499673653 0
//Ethernet24 2976898626185 1587315888369 0 3720900543543 915817056263 0
//Ethernet28 3530619214310 1592361800267 0 2435798298488 1055980373566 0
//Ethernet32 1736687809478 1468562376259 0 2956777802728 998205405421 0
//Ethernet36 1177239757684 2378638162754 0 3281846160288 376319942357 0
//Ethernet40 1919963381701 1659665440142 0 2082884939332 1071753723587 0
//Ethernet44 3301425076496 1043799168311 0 2785942773405 234051655477 0
//Ethernet48 2861253648712 2607572907505 0 2305011130777 882745793901 0
//Ethernet52 2657621482082 2379791009054 0 2563218954334 547784979886 0
//Ethernet56 1579337757748 1318026597693 0 1481987331729 273225187315 0
//Ethernet60 2579880962972 1738780342206 0 3937968261527 300408567763 0
//Ethernet64 3098150711440 1982678989394 0 2718722642778 614889086585 0
//Ethernet68 2303723457227 2003216812516 0 3473275923942 1079108542357 0
//Ethernet72 2922526714740 2329990436128 0 2685253445396 488679336887 0
//Ethernet76 2000082104935 2008266939250 0 1864736644279 72442934988 0
//Ethernet80 1895668504939 1764958708863 0 2735664066143 159710846588 0
//Ethernet84 1547696085431 2236485893910 0 2067406122712 581349015388 0
//Ethernet88 2913148855672 640953222278 0 2860150441089 411691814934 0
//Ethernet92 3015098438525 2209788973943 0 2092503052332 499918289683 0
//Ethernet96 3729010717742 2035668903092 0 3587128038837 326702266862 0
//Ethernet100 3921732792118 849330454732 0 2688196061230 736452541765 0
//Ethernet104 3102950334553 1139837252765 0 3398934672308 234736768449 0
//Ethernet108 3473523684145 1774944030225 0 2990281101686 1094743647467 0
//Ethernet112 2948981687214 2803727841764 0 2142343083007 476568048283 0
//Ethernet116 2444901828803 3079140972446 0 2728375868725 1044599136395 0
//Ethernet120 3841189447065 2627771041383 0 2932649549811 669644450653 0
//Ethernet124 3481911419531 1484808049535 0 2635977290237 842710039415 0`

	// 0:IFACE
	// 1:RX_OK
	// 2:RX_ERR
	// 3:RX_DRP
	// 4:TX_OK
	// 5:TX_ERR
	// 6:TX_DRP

	t2 := time.Now()
	intfStatusArray := strings.Split(string(intfStatus), "\n")
	portstatArray := strings.Split(string(portstat), "\n")

	device.Interfaces = &gostruct.OpenconfigInterfaces_Interfaces{}

	//inf, err := device.Interfaces.NewInterface("eth0")
	//if err != nil {
	//	return err
	//}
	//
	//inf.Config.Name = ygot.String("eth0")
	//inf.State=&gostruct.OpenconfigInterfaces_Interfaces_Interface_State{
	//	AdminStatus:gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_UP,
	//	OperStatus:gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_UP,
	//
	//}
	//inf.Enabled = ygot.Bool(true)
	//inf.AdminStatus = gostruct.OpenconfigInterfaces_Interface_AdminStatus_UP
	//inf.OperStatus = gostruct.OpenconfigInterfaces_Interface_OperStatus_UP

	for j := 0; j < len(intfStatusArray)-1; j++ {
		intDetail := strings.Split(intfStatusArray[j], " ")
		//fmt.Println("Doing str: ", intDetail)
		intStatDetail := strings.Split(portstatArray[j], " ")

		intName := strings.Split(intfStatusArray[j], " ")[0]

		mtu, err := strconv.Atoi(intDetail[3])
		if err != nil {
			return err
		}

		InOctets, err := strconv.ParseUint(intStatDetail[1], 10, 64)
		if err != nil {
			panic(err)
		}

		InErrors, err := strconv.ParseUint(intStatDetail[2], 10, 64)
		if err != nil {
			panic(err)
		}

		InDiscards, err := strconv.ParseUint(intStatDetail[3], 10, 64)
		if err != nil {
			panic(err)
		}

		OutOctets, err := strconv.ParseUint(intStatDetail[4], 10, 64)
		if err != nil {
			panic(err)
		}

		OutErrors, err := strconv.ParseUint(intStatDetail[5], 10, 64)
		if err != nil {
			panic(err)
		}

		OutDiscards, err := strconv.ParseUint(intStatDetail[6], 10, 64)
		if err != nil {
			panic(err)
		}

		inf, err := device.Interfaces.NewInterface(intName)
		inf.Config = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Config{
			Name: ygot.String(intName),
		}

		var enabled *bool
		var adminStatus gostruct.E_OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus
		var operStatus gostruct.E_OpenconfigInterfaces_Interfaces_Interface_State_OperStatus

		if intDetail[6] == "up" {
			enabled = ygot.Bool(true)
			adminStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_UP
		} else {
			enabled = ygot.Bool(false)
			adminStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_DOWN
		}

		if intDetail[5] == "up" {
			operStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_UP
		} else {
			operStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_DOWN
		}

		inf.State = &gostruct.OpenconfigInterfaces_Interfaces_Interface_State{
			AdminStatus: adminStatus,
			OperStatus:  operStatus,
			Enabled:     enabled,
			Mtu:         ygot.Uint16(uint16(mtu)),
			Counters: &gostruct.OpenconfigInterfaces_Interfaces_Interface_State_Counters{
				InOctets:    ygot.Uint64(InOctets),
				InErrors:    ygot.Uint64(InErrors),
				InDiscards:  ygot.Uint64(InDiscards),
				OutOctets:   ygot.Uint64(OutOctets),
				OutErrors:   ygot.Uint64(OutErrors),
				OutDiscards: ygot.Uint64(OutDiscards),
			},
		}

	}
	t2d := time.Since(t2)

	fmt.Println("Parsed all Interface information, the timeing information is following")
	fmt.Printf("Get 1 command, took %s \n", t0d)
	fmt.Printf("Get 2 command, took %s \n", t1d)
	fmt.Printf("Parse all data, took %s \n", t2d)
	return nil
}

func initInterfaceEthernet() {

}

func initInterfaceAggregate() {

}

func initInterfaceIp() {

}

func initPlatform(device *gostruct.Device) error {

	t0 := time.Now()
	cmd:="show platform syseeprom | grep '0x' | awk -F'0x' '{print $2,$3}' | awk -F' ' '{print \"0x\"$1,$3}'"

	syseeprom, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmd)
	}
	t0d := time.Since(t0)
	t1 := time.Now()

	//Manufacture Date     0x25
	//Label Revision       0x27
	//Platform Name        0x28
	//ONIE Version         0x29
	//Manufacturer         0x2B
	//Manufacture Country  0x2C
	//Diag Version         0x2E
	//Base MAC Address     0x24
	//Serial Number        0x23
	//Part Number          0x22
	//Product Name         0x21
	//MAC Addresses        0x2A
	//Vendor Name          0x2D
	//CRC-32               0xFE
	//(checksum valid)`

//	syseeprom := `0x25 05/25/2017
//0x27 R02B
//0x28 x86_64-accton_as7712_32x-r0
//0x29 2016.08.00.03
//0x2B Accton
//0x2C TW
//0x2E 0.0.5.3
//0x24 A8:2B:B5:38:01:1D
//0x23 771232X1721087
//0x22 FP3ZZ7632014A
//0x21 7712-32X-O-AC-F
//0x2A 131
//0x2D Edgecore
//0xFE CB07A5F3`

	cmd="show version | grep 'Software Version'"
	versionOutput, err := exec.Command("bash","-c",cmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmd)
	}
	t1d := time.Since(t1)
	t2 := time.Now()
	//versionOutput := `SONiC Software Version: SONiC.HEAD.638-a0bd656`

	consoleArray := strings.Split(string(syseeprom), "\n")

	product, err := search(consoleArray, "0x21")
	hwVersion, err := search(consoleArray, "0x28")
	serialNo, err := search(consoleArray, "0x23")
	partNo, err := search(consoleArray, "0x22")
	mfgName, err := search(consoleArray, "0x2B")
	mfgDate, err := search(consoleArray, "0x25")
	swVersion := strings.Split(string(versionOutput), ": ")[1]

	device.Components = &gostruct.OpenconfigPlatform_Components{}

	compnent, err := device.Components.NewComponent(product)
	if err != nil {
		return err
	}

	compnent.Name = ygot.String(product)
	compnent.Config = &gostruct.OpenconfigPlatform_Components_Component_Config{
		Name: ygot.String(product),
	}
	compnent.State = &gostruct.OpenconfigPlatform_Components_Component_State{
		SerialNo:        ygot.String(serialNo),
		PartNo:          ygot.String(partNo),
		HardwareVersion: ygot.String(hwVersion),
		SoftwareVersion: ygot.String(swVersion),
		MfgName:         ygot.String(mfgName),
		MfgDate:         ygot.String(mfgDate),
		Type: &gostruct.OpenconfigPlatform_Components_Component_State_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT{
			gostruct.OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_CHASSIS,
		},
	}

	compnent.State.SerialNo = ygot.String(serialNo)

	t2d := time.Since(t2)

	fmt.Println("Parsed all components information, the timeing information is following")
	fmt.Printf("Get 1 command, took %s \n", t0d)
	fmt.Printf("Get 2 command, took %s \n", t1d)
	fmt.Printf("Parse all data, took %s \n", t2d)

	return nil
}

func initVlan(device *gostruct.Device) error {

	// TODO: Should gather information from switch, here is still use static data

	// I use the Switch Serial Number to this forwarding instance unique name
	// The Serial Number is come from following command
	// show platform syseeprom | grep 0x23 | awk -F' ' '{print $5}'
	//int, err := device.NewNetworkInstance("switch1")
	//if err != nil {
	//	return err
	//}
	//
	//vlan1, err := int.NewVlan(uint16(1))
	//if err != nil {
	//	return err
	//}
	//vlan1.Name = ygot.String("VLAN1")
	//
	//for idx := 0; idx < len(device.Interface)-1; idx++ {
	//	name := device.GetInterface("Ethernet" + strconv.Itoa(idx*4)).Name
	//	vlan1.Member = append(vlan1.Member, &gostruct.NetworkInstance_Vlan_Member{
	//		Interface: name,
	//	})
	//}
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

	err := initPlatform(device)
	if err != nil {
		return nil, err
	}

	err = initLldp(device)
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

func search(s []string, tgt string) (string, error) {
	for _, c := range s {
		if strings.Contains(c, tgt) {
			return strings.Split(c, " ")[1], nil
		}
	}
	return "", fmt.Errorf("Not found")
}
