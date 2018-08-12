package components_funcs

import (
	"fmt"
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"sync"
	"github.com/go-redis/redis"
	"github.com/tidwall/gjson"
	"html/template"
	"bytes"
)

func InitInterface(device *gostruct.Device, appClient *redis.Client, configClient *redis.Client) error {

	device.Interfaces = &gostruct.OpenconfigInterfaces_Interfaces{}

	portstatCmd := "portstat | grep 'Ethernet' | awk -F' ' '{print $1,$3,$6,$7,$9,$12,$13}'"
	portstat, err := exec.Command("bash", "-c", portstatCmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", portstatCmd)
	}
	//		portstat := `Ethernet0 1595169299934 1824671749079 0 3739096352458 837308284053 0
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

	portstatArray := strings.Split(string(portstat), "\n")

	for i := 0; i < len(portstatArray)-1; i++ {
		intStatDetail := strings.Split(portstatArray[i], " ")

		intfName := intStatDetail[0]

		intf, err := device.Interfaces.NewInterface(intfName)
		if err != nil {
			return err
		}

		intfVlan := appClient.Keys("VLAN_MEMBER_TABLE:*:" + intfName).Val()
		//1) "VLAN_MEMBER_TABLE:Vlan1000:Ethernet12"
		//2) "VLAN_MEMBER_TABLE:Vlan12:Ethernet12"

		var accessVlan *uint16
		var nativeVlan *uint16
		var intfMode gostruct.E_OpenconfigVlan_VlanModeType
		trunkConfigVlans := make([]gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_Config_TrunkVlans_Union, 0)
		trunkStateVlans := make([]gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_State_TrunkVlans_Union, 0)

		if len(intfVlan) > 1 {
			intfMode = gostruct.OpenconfigVlan_VlanModeType_TRUNK
			for _, vlanMember := range intfVlan {
				if appClient.HGet(vlanMember, "tagging_mode").Val() == "untagged" {
					nativeVlanInt, err := strconv.Atoi(strings.Split(strings.Split(vlanMember, ":")[1], "Vlan")[1])
					if err != nil {
						return err
					}
					nativeVlan = ygot.Uint16(uint16(nativeVlanInt))
				} else {
					trunkVlanInt, err := strconv.Atoi(strings.Split(strings.Split(vlanMember, ":")[1], "Vlan")[1])
					if err != nil {
						return err
					}
					trunkConfigVlans = append(trunkConfigVlans, &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_Config_TrunkVlans_Union_Uint16{
						Uint16: uint16(trunkVlanInt),
					})
					trunkStateVlans = append(trunkStateVlans, &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_State_TrunkVlans_Union_Uint16{
						Uint16: uint16(trunkVlanInt),
					})
				}
			}
		} else {
			intfMode = gostruct.OpenconfigVlan_VlanModeType_ACCESS
			accessVlanInt, err := strconv.Atoi(strings.Split(strings.Split(intfVlan[0], ":")[1], "Vlan")[1])
			if err != nil {
				return err
			}
			accessVlan = ygot.Uint16(uint16(accessVlanInt))
			nativeVlan = accessVlan
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

		var enabled *bool
		var adminStatus gostruct.E_OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus
		var operStatus gostruct.E_OpenconfigInterfaces_Interfaces_Interface_State_OperStatus
		var mtu uint16

		if appClient.HGet("PORT_TABLE:"+intfName, "oper_status").Val() == "up" {
			operStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_UP
		} else {
			operStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_DOWN
		}

		if appClient.HGet("PORT_TABLE:"+intfName, "admin_status").Val() == "up" {
			enabled = ygot.Bool(true)
			adminStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_UP
		} else {
			enabled = ygot.Bool(false)
			adminStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_DOWN
		}

		mtu64, err := appClient.HGet("PORT_TABLE:"+intfName, "mtu").Uint64()
		if err != nil {
			return err
		}
		mtu = uint16(mtu64)

		intf.Config = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Config{
			Name: ygot.String(intfName),
		}

		intf.State = &gostruct.OpenconfigInterfaces_Interfaces_Interface_State{
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

		intf.Ethernet = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet{
			Config: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_Config{},
			State:  &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_State{},
			SwitchedVlan: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan{
				Config: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_Config{
					AccessVlan:    accessVlan,
					InterfaceMode: intfMode,
					NativeVlan:    nativeVlan,
					TrunkVlans:    trunkConfigVlans,
				},
				State: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_State{
					AccessVlan:    accessVlan,
					InterfaceMode: intfMode,
					NativeVlan:    nativeVlan,
					TrunkVlans:    trunkStateVlans,
				},
			},
		}
	}

	return nil
}

func SyncInterface(device *gostruct.Device, appClient *redis.Client, configClient *redis.Client, mu *sync.RWMutex) error {
	for {
		portstatCmd := "portstat | grep 'Ethernet' | awk -F' ' '{print $1,$3,$6,$7,$9,$12,$13}'"
		portstat, err := exec.Command("bash", "-c", portstatCmd).Output()
		if err != nil {
			return fmt.Errorf("Failed to execute command: %s", portstatCmd)
		}
		//		portstat := `Ethernet0 1595169299934 1824671749079 0 3739096352458 837308284053 0
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

		portstatArray := strings.Split(string(portstat), "\n")

		for i := 0; i < len(portstatArray)-1; i++ {
			intStatDetail := strings.Split(portstatArray[i], " ")

			intfName := intStatDetail[0]

			intf := device.Interfaces.GetOrCreateInterface(intfName)

			intfVlan := appClient.Keys("VLAN_MEMBER_TABLE:*:" + intfName).Val()
			//1) "VLAN_MEMBER_TABLE:Vlan1000:Ethernet12"
			//2) "VLAN_MEMBER_TABLE:Vlan12:Ethernet12"

			var accessVlan *uint16
			var nativeVlan *uint16
			var intfMode gostruct.E_OpenconfigVlan_VlanModeType
			trunkConfigVlans := make([]gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_Config_TrunkVlans_Union, 0)
			trunkStateVlans := make([]gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_State_TrunkVlans_Union, 0)

			if len(intfVlan) > 1 {
				intfMode = gostruct.OpenconfigVlan_VlanModeType_TRUNK
				for _, vlanMember := range intfVlan {
					if appClient.HGet(vlanMember, "tagging_mode").Val() == "untagged" {
						nativeVlanInt, err := strconv.Atoi(strings.Split(strings.Split(vlanMember, ":")[1], "Vlan")[1])
						if err != nil {
							return err
						}
						nativeVlan = ygot.Uint16(uint16(nativeVlanInt))
					} else {
						trunkVlanInt, err := strconv.Atoi(strings.Split(strings.Split(vlanMember, ":")[1], "Vlan")[1])
						if err != nil {
							return err
						}
						trunkConfigVlans = append(trunkConfigVlans, &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_Config_TrunkVlans_Union_Uint16{
							Uint16: uint16(trunkVlanInt),
						})
						trunkStateVlans = append(trunkStateVlans, &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_State_TrunkVlans_Union_Uint16{
							Uint16: uint16(trunkVlanInt),
						})
					}
				}
			} else {
				intfMode = gostruct.OpenconfigVlan_VlanModeType_ACCESS
				accessVlanInt, err := strconv.Atoi(strings.Split(strings.Split(intfVlan[0], ":")[1], "Vlan")[1])
				if err != nil {
					return err
				}
				accessVlan = ygot.Uint16(uint16(accessVlanInt))
				nativeVlan = accessVlan
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

			var enabled *bool
			var adminStatus gostruct.E_OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus
			var operStatus gostruct.E_OpenconfigInterfaces_Interfaces_Interface_State_OperStatus
			var mtu uint16

			if appClient.HGet("PORT_TABLE:"+intfName, "oper_status").Val() == "up" {
				operStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_UP
			} else {
				operStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_DOWN
			}

			if appClient.HGet("PORT_TABLE:"+intfName, "admin_status").Val() == "up" {
				enabled = ygot.Bool(true)
				adminStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_UP
			} else {
				enabled = ygot.Bool(false)
				adminStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_DOWN
			}

			mtu64, err := appClient.HGet("PORT_TABLE:"+intfName, "mtu").Uint64()
			if err != nil {
				return err
			}
			mtu = uint16(mtu64)

			intf.Config = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Config{
				Name: ygot.String(intfName),
			}

			intf.State = &gostruct.OpenconfigInterfaces_Interfaces_Interface_State{
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

			intf.Ethernet = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet{
				Config: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_Config{},
				State:  &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_State{},
				SwitchedVlan: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan{
					Config: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_Config{
						AccessVlan:    accessVlan,
						InterfaceMode: intfMode,
						NativeVlan:    nativeVlan,
						TrunkVlans:    trunkConfigVlans,
					},
					State: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Ethernet_SwitchedVlan_State{
						AccessVlan:    accessVlan,
						InterfaceMode: intfMode,
						NativeVlan:    nativeVlan,
						TrunkVlans:    trunkStateVlans,
					},
				},
			}
		}
		fmt.Println("Syncing Interfaces...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

func InitInterfaceAggregate(device *gostruct.Device, client *redis.Client) error {

	//device.Lacp = &gostruct.OpenconfigLacp_Lacp{}
	//intfs := device.Lacp.GetOrCreateInterfaces()

	intfs := device.GetOrCreateInterfaces()

	portChannels := client.Keys("PORTCHANNEL|*")

	//1) "PORTCHANNEL|ThisIsNotPC"
	//2) "PORTCHANNEL|PortChannel04"

	for _, value := range portChannels.Val() {
		intfName := strings.Split(value, "|")[1]
		intf, err := intfs.NewInterface(intfName)
		if err != nil {
			return err
		}

		intf.Config = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Config{
			Name: ygot.String(intfName),
		}

		intf.State = &gostruct.OpenconfigInterfaces_Interfaces_Interface_State{
			Enabled:     ygot.Bool(false),
			AdminStatus: gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_DOWN,
			OperStatus:  gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_DOWN,
		}

		cmd := "docker exec teamd cat /etc/teamd/" + intfName + ".conf"
		json, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			return fmt.Errorf("Failed to execute command: %s", cmd)
		}

		dbMembers := client.HGet(value, "members@")
		dbMemberIntfs := strings.Split(dbMembers.Val(), ",")
		member := make([]string, 0)
		if dbMembers.Val() != "" {
			for _, value := range dbMemberIntfs {
				member = append(member, value)
			}
		}

		var lagType gostruct.E_OpenconfigIfAggregate_AggregationType
		if gjson.Get(string(json), "runner.name").Value() == "lacp" {
			lagType = gostruct.OpenconfigIfAggregate_AggregationType_LACP
		} else {
			lagType = gostruct.OpenconfigIfAggregate_AggregationType_UNSET
		}

		intf.Aggregation = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Aggregation{
			Config: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Aggregation_Config{
				LagType:  lagType,
				MinLinks: ygot.Uint16(uint16(gjson.Get(string(json), "runner.min_ports").Uint())),
			},
			State: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Aggregation_State{
				LagType:  lagType,
				MinLinks: ygot.Uint16(uint16(gjson.Get(string(json), "runner.min_ports").Uint())),
				Member:   member,
			},
		}
	}
	return nil
}

func SyncInterfaceAggregate(device *gostruct.Device, client *redis.Client, mu *sync.RWMutex, ) error {

	//device.Lacp = &gostruct.OpenconfigLacp_Lacp{}
	//intfs := device.Lacp.GetOrCreateInterfaces()

	for {
		intfs := device.Interfaces

		portChannels := client.Keys("PORTCHANNEL|*")

		//1) "PORTCHANNEL|ThisIsNotPC"
		//2) "PORTCHANNEL|PortChannel04"

		for _, value := range portChannels.Val() {
			intfName := strings.Split(value, "|")[1]
			intf := intfs.GetOrCreateInterface(intfName)

			cmd := "sudo ip link show " + intfName + " up | grep state | awk -F' ' '{print $9}'"
			intfResult, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				return fmt.Errorf("Failed to execute command: %s", cmd)
			}

			enabled := ygot.Bool(false)
			adminStatus := gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_DOWN
			operStatus := gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_DOWN

			if string(intfResult) == "UP" {
				enabled = ygot.Bool(true)
				adminStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_AdminStatus_UP
				operStatus = gostruct.OpenconfigInterfaces_Interfaces_Interface_State_OperStatus_UP
			}

			intf.State = &gostruct.OpenconfigInterfaces_Interfaces_Interface_State{
				Enabled:     enabled,
				AdminStatus: adminStatus,
				OperStatus:  operStatus,
			}

			cmd = "docker exec teamd cat /etc/teamd/" + intfName + ".conf"
			json, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				return fmt.Errorf("Failed to execute command: %s", cmd)
			}

			dbMembers := client.HGet(value, "members@")
			dbMemberIntfs := strings.Split(dbMembers.Val(), ",")
			member := make([]string, 0)
			if dbMembers.Val() != "" {
				for _, value := range dbMemberIntfs {
					member = append(member, value)
				}
			}

			var lagType gostruct.E_OpenconfigIfAggregate_AggregationType
			if gjson.Get(string(json), "runner.name").Value() == "lacp" {
				lagType = gostruct.OpenconfigIfAggregate_AggregationType_LACP
			} else {
				lagType = gostruct.OpenconfigIfAggregate_AggregationType_UNSET
			}

			intf.Aggregation = &gostruct.OpenconfigInterfaces_Interfaces_Interface_Aggregation{
				Config: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Aggregation_Config{
					LagType:  lagType,
					MinLinks: ygot.Uint16(uint16(gjson.Get(string(json), "runner.min_ports").Uint())),
				},
				State: &gostruct.OpenconfigInterfaces_Interfaces_Interface_Aggregation_State{
					LagType:  lagType,
					MinLinks: ygot.Uint16(uint16(gjson.Get(string(json), "runner.min_ports").Uint())),
					Member:   member,
				},
			}
		}
		fmt.Println("Syncing Interfaces Aggregate ...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

func SetInterfacePortchannel(key []string, value []string, str string, b bool) error {

	//-update /interfaces/interface[name=test]/aggregation/config/lag-type:lacp

	type Teamd struct {
		Device  string
		HwAddr  string
		Runner  string
		MinPort int
	}

	filePath := "/etc/teamd/" + str + ".conf"

	TEAMD_CONF_TMPL := `
{
    "device": "{{.Device}}",
    "hwaddr": "{{.HwAddr}}",
    "runner": {
        "name": "{{.Runner}}",
        "active": true,
        "min_ports": {{.MinPort}},
        "tx_hash": ["eth", "ipv4", "ipv6"]
    },
    "link_watch": {
        "name": "ethtool"
    },
    "ports": {
    }
}`

	cmdGetMac := "ip link show eth0 | grep ether | awk -F' ' '{print $2}'"
	mac, err := exec.Command("bash", "-c", cmdGetMac).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmdGetMac)
	}

	teml, err := template.New("teamConfig").Parse(TEAMD_CONF_TMPL)
	if err != nil {
		return status.Error(codes.Internal, "Failed to generator team config file")
	}

	var config bytes.Buffer
	err = teml.Execute(&config, Teamd{
		Device:  value[0],
		HwAddr:  string(mac[:len(mac)-2]),
		Runner:  str,
		MinPort: 1,
	})
	if err != nil {
		return status.Error(codes.Internal, "Failed to generator team config file")
	}

	cmdWrConfig := "sonic-cfggen -a '{\"PORTCHANNEL\": {\"" + str + "\":{}}}' --write-to-db"
	_, err = exec.Command("bash", "-c", cmdWrConfig).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmdWrConfig)
	}

	cmdGenConfig := "echo '" + config.String() + "' | (docker exec -i teamd bash -c 'cat > " + filePath + "')"
	_, err = exec.Command("bash", "-c", cmdGenConfig).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmdGenConfig)
	}

	cmdReadConfig := "docker exec teamd teamd -d -f " + filePath
	_, err = exec.Command("bash", "-c", cmdReadConfig).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmdReadConfig)
	}

	return nil
}

func DelInterfacePortchannel(key []string, value []string, str string, b bool) error {
	return status.Error(codes.Unimplemented, "DelInterfacePortchannel is not implement yet")
}

func SetInterfacePortchannelMember(key []string, value []string, str string, b bool) error {
	//-update '/interfaces/interface[name=Ethernet2]/ethernet/config/aggregate-id:PortChannel2

	intfName := value[0]

	cmd := "teamdctl " + str + " port add " + intfName
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}

	return nil
}

func DelInterfacePortchannelMember(key []string, value []string, str string, b bool) error {
	intfName := value[0]

	cmd := "teamdctl " + str + " port remove " + intfName
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}

	return nil
}

func SetInterfaceConfigEnabled(key []string, value []string, str string, b bool) error {
	var cmd string

	enable, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}
	if enable {
		cmd = "sudo config interface startup " + value[len(value)-1]
	} else {
		cmd = "sudo config interface shutdown " + value[len(value)-1]
	}

	result, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}
	if strings.Contains(string(result), "Cannot find device") {
		return status.Error(codes.Internal, "Cannot find device: "+value[len(value)-1])
	}
	return nil
}

func initInterfaceEthernet() {
}

func initInterfaceIp() {

}
