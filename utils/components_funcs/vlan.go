package components_funcs

import "github.com/breezestars/gnxi/gnmi/modeldata/gostruct"

func InitVlan(device *gostruct.Device) error {

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
