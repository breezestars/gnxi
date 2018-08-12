package components_funcs

import (
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/go-redis/redis"
	"strings"
	"github.com/openconfig/ygot/ygot"
		"fmt"
	"time"
	"sync"
	)

func InitLldp(device *gostruct.Device, client *redis.Client) error {

	device.Lldp = &gostruct.OpenconfigLldp_Lldp{
		Config: &gostruct.OpenconfigLldp_Lldp_Config{
			Enabled: ygot.Bool(true),
		},
	}

	intfs := device.Lldp.GetOrCreateInterfaces()

	lldpKeys := client.Keys("LLDP_ENTRY_TABLE:*")
	for _, value := range lldpKeys.Val() {
		intfName := strings.Split(value, ":")[1]
		intf, err := intfs.NewInterface(intfName)
		if err != nil {
			return err
		}

		lldp_rem_chassis_id := client.HGet(value, "lldp_rem_chassis_id").Val()
		lldp_rem_sys_desc := client.HGet(value, "lldp_rem_sys_desc").Val()
		lldp_rem_port_desc := client.HGet(value, "lldp_rem_port_desc").Val()
		lldp_rem_sys_name := client.HGet(value, "lldp_rem_sys_name").Val()
		lldp_rem_port_id := client.HGet(value, "lldp_rem_port_id").Val()
		lldp_rem_chassis_id_subtype, err := client.HGet(value, "lldp_rem_chassis_id_subtype").Int64()
		if err != nil {
			return err
		}
		lldp_rem_time_mark, err := client.HGet(value, "lldp_rem_time_mark").Uint64()
		if err != nil {
			return err
		}
		lldp_rem_port_id_subtype, err := client.HGet(value, "lldp_rem_port_id_subtype").Int64()
		if err != nil {
			return err
		}

		intf.Config = &gostruct.OpenconfigLldp_Lldp_Interfaces_Interface_Config{
			Name: ygot.String(intfName),
		}
		nbrs := intf.GetOrCreateNeighbors()
		nbr, err := nbrs.NewNeighbor(lldp_rem_chassis_id)
		if err != nil {
			return err
		}

		nbr.State = &gostruct.OpenconfigLldp_Lldp_Interfaces_Interface_Neighbors_Neighbor_State{
			Id:                ygot.String(lldp_rem_chassis_id),
			ChassisIdType:     gostruct.E_OpenconfigLldp_ChassisIdType(lldp_rem_chassis_id_subtype),
			ChassisId:         ygot.String(lldp_rem_chassis_id),
			SystemDescription: ygot.String(lldp_rem_sys_desc),
			Age:               ygot.Uint64(lldp_rem_time_mark),
			PortDescription:   ygot.String(lldp_rem_port_desc),
			PortIdType:        gostruct.E_OpenconfigLldp_PortIdType(lldp_rem_port_id_subtype),
			PortId:            ygot.String(lldp_rem_port_id),
			SystemName:        ygot.String(lldp_rem_sys_name),
		}

		//1) "lldp_rem_chassis_id_subtype"
		//2) "4"
		//3) "lldp_rem_chassis_id"
		//4) "cc:37:ab:60:cd:c1"
		//5) "lldp_rem_sys_desc"
		//6) "PICA8 Inc., Model as4610_54, PicOS 2.6.2"
		//7) "lldp_rem_time_mark"
		//8) "3220"
		//9) "lldp_rem_port_desc"
		//10) "ge-1/1/28"
		//11) "lldp_rem_port_id_subtype"
		//12) "5"
		//13) "lldp_rem_sys_name"
		//14) "B2-Tor"
		//15) "lldp_rem_port_id"
		//16) "ge-1/1/28"

	}

	return nil
}

func SyncLldp(device *gostruct.Device, client *redis.Client, mu *sync.RWMutex) error {
	for {

		intfs := device.Lldp.GetOrCreateInterfaces()

		lldpKeys := client.Keys("LLDP_ENTRY_TABLE:*")
		for _, value := range lldpKeys.Val() {
			intfName := strings.Split(value, ":")[1]
			intf := intfs.GetOrCreateInterface(intfName)

			lldp_rem_chassis_id := client.HGet(value, "lldp_rem_chassis_id").Val()
			lldp_rem_sys_desc := client.HGet(value, "lldp_rem_sys_desc").Val()
			lldp_rem_port_desc := client.HGet(value, "lldp_rem_port_desc").Val()
			lldp_rem_sys_name := client.HGet(value, "lldp_rem_sys_name").Val()
			lldp_rem_port_id := client.HGet(value, "lldp_rem_port_id").Val()
			lldp_rem_chassis_id_subtype, err := client.HGet(value, "lldp_rem_chassis_id_subtype").Int64()
			if err != nil {
				return err
			}
			lldp_rem_time_mark, err := client.HGet(value, "lldp_rem_time_mark").Uint64()
			if err != nil {
				return err
			}
			lldp_rem_port_id_subtype, err := client.HGet(value, "lldp_rem_port_id_subtype").Int64()
			if err != nil {
				return err
			}

			intf.Neighbors=&gostruct.OpenconfigLldp_Lldp_Interfaces_Interface_Neighbors{}
			nbrs := intf.GetNeighbors()
			nbr, err := nbrs.NewNeighbor(lldp_rem_chassis_id)
			if err != nil {
				return err
			}

			nbr.State = &gostruct.OpenconfigLldp_Lldp_Interfaces_Interface_Neighbors_Neighbor_State{
				Id:                ygot.String(lldp_rem_chassis_id),
				ChassisIdType:     gostruct.E_OpenconfigLldp_ChassisIdType(lldp_rem_chassis_id_subtype),
				ChassisId:         ygot.String(lldp_rem_chassis_id),
				SystemDescription: ygot.String(lldp_rem_sys_desc),
				Age:               ygot.Uint64(lldp_rem_time_mark),
				PortDescription:   ygot.String(lldp_rem_port_desc),
				PortIdType:        gostruct.E_OpenconfigLldp_PortIdType(lldp_rem_port_id_subtype),
				PortId:            ygot.String(lldp_rem_port_id),
				SystemName:        ygot.String(lldp_rem_sys_name),
			}

			//1) "lldp_rem_chassis_id_subtype"
			//2) "4"
			//3) "lldp_rem_chassis_id"
			//4) "cc:37:ab:60:cd:c1"
			//5) "lldp_rem_sys_desc"
			//6) "PICA8 Inc., Model as4610_54, PicOS 2.6.2"
			//7) "lldp_rem_time_mark"
			//8) "3220"
			//9) "lldp_rem_port_desc"
			//10) "ge-1/1/28"
			//11) "lldp_rem_port_id_subtype"
			//12) "5"
			//13) "lldp_rem_sys_name"
			//14) "B2-Tor"
			//15) "lldp_rem_port_id"
			//16) "ge-1/1/28"

		}

		fmt.Println("Syncing Lldp...")
		time.Sleep(5 * time.Second)
	}
	return nil
}
