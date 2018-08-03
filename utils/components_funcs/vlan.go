package components_funcs

import (
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/go-redis/redis"
	"strconv"
	"github.com/openconfig/ygot/ygot"
)

func InitVlan(device *gostruct.Device, client *redis.Client) error {

	device.NetworkInstances = &gostruct.OpenconfigNetworkInstance_NetworkInstances{}
	// TODO: Here is use static 'sonic' to be name, need to change to dynamic
	ni, err := device.NetworkInstances.NewNetworkInstance("sonic")
	if err != nil {
		return err
	}
	ni.Vlans = &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans{}

	keys := client.Keys("VLAN|*")
	//vlanMembers := client.Keys("VLAN_MEMBER|*")

	for _, v := range keys.Val() {
		//fmt.Println("key is: ", k, " and value is: ", string(v[9:]))
		vid := string(v[9:])
		vlanId, err := strconv.Atoi(vid)
		if err != nil {
			return err
		}

		vlan := ni.Vlans.GetOrCreateVlan(uint16(vlanId))
		vlanMbs := vlan.GetOrCreateMembers()

		intfs := client.Keys("VLAN_MEMBER|Vlan" + vid + "|*")

		for _, intf := range intfs.Val() {
			vlanMbs.Member = append(vlanMbs.Member, &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member{
				State: &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member_State{
					Interface: ygot.String(string(intf[21:])),
				},
			})
		}

		// TODO: Need to Check the performance, which will be more faster? Redis search or process by for in for
		//for _, vlanMember := range vlanMembers.Val(){
		//	fmt.Println("range in "+vlanMember+" and will filter with: VLAN"+vid)
		//	if strings.Contains(vlanMember, "Vlan"+vid) {
		//		// VLAN_MEMBER|Vlan1000|
		//		intf :=  string(vlanMember[21:])
		//		fmt.Println("Got "+intf)
		//
		//		vlanMbs.Member = append(vlanMbs.Member, &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member{
		//			State: &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member_State{
		//				Interface: ygot.String(intf),
		//			},
		//		})
		//
		//	}
		//}
	}
	return nil
}

func SyncVlan(device *gostruct.Device, client *redis.Client) error {
	// TODO: Here is use static 'sonic' to be name, need to change to dynamic
	vlans := device.NetworkInstances.GetNetworkInstance("sonic").Vlans

	keys := client.Keys("VLAN|*")

	for _, v := range keys.Val() {
		//fmt.Println("key is: ", k, " and value is: ", string(v[9:]))
		vid := string(v[9:])
		vlanId, err := strconv.Atoi(vid)
		if err != nil {
			return err
		}

		vlan := vlans.GetOrCreateVlan(uint16(vlanId))
		vlan.Members = &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members{}
		vlanMbs := vlan.Members

		intfs := client.Keys("VLAN_MEMBER|Vlan" + vid + "|*")

		for _, intf := range intfs.Val() {
			vlanMbs.Member = append(vlanMbs.Member, &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member{
				State: &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member_State{
					Interface: ygot.String(string(intf[21:])),
				},
			})
		}
	}

	return nil
}
