package components_funcs

import (
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/go-redis/redis"
	"strconv"
	"github.com/openconfig/ygot/ygot"
	"strings"
	"time"
	"fmt"
	"sync"
	"os/exec"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

const (
	DEFAULT_VLAN = "1000"
)

func InitVlan(device *gostruct.Device, client *redis.Client) error {

	device.NetworkInstances = &gostruct.OpenconfigNetworkInstance_NetworkInstances{}
	// TODO: Here is use static 'sonic' to be name, need to change to dynamic
	ni, err := device.NetworkInstances.NewNetworkInstance("sonic")
	if err != nil {
		return err
	}
	ni.Config = &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Config{
		Name: ygot.String("sonic"),
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
		ni.Vlans.Vlan[uint16(vlanId)].Config = &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Config{
			VlanId: ygot.Uint16(uint16(vlanId)),
		}
		vlanMbs := vlan.GetOrCreateMembers()

		intfs := client.Keys("VLAN_MEMBER|Vlan" + vid + "|*")

		for _, intf := range intfs.Val() {
			vlanMbs.Member = append(vlanMbs.Member, &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member{
				State: &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member_State{
					Interface: ygot.String(strings.Split(intf, "Vlan"+vid+"|")[1]),
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

func SyncVlan(device *gostruct.Device, client *redis.Client, mu *sync.RWMutex) error {
	for {
		// TODO: Here is use static 'sonic' to be name, need to change to dynamic
		device.NetworkInstances.NetworkInstance["sonic"].Vlans = &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans{}
		vlans := device.NetworkInstances.NetworkInstance["sonic"].Vlans

		keys := client.Keys("VLAN|*")

		for _, v := range keys.Val() {
			//fmt.Println("key is: ", k, " and value is: ", string(v[9:]))
			vid := string(v[9:])
			vlanId, err := strconv.Atoi(vid)
			if err != nil {
				fmt.Println("Error 1")
				return err
			}

			vlan, err := vlans.NewVlan(uint16(vlanId))
			if err != nil {
				fmt.Println("Error 2")
				return err
			}
			vlan.Config = &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Config{
				VlanId: ygot.Uint16(uint16(vlanId)),
			}

			vlan.Members = &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members{
				Member: nil,
			}
			vlanMbs := vlan.Members

			intfs := client.Keys("VLAN_MEMBER|Vlan" + vid + "|*")

			memberSlice := make([]*gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member, 0)
			for _, intf := range intfs.Val() {

				memberSlice = append(memberSlice, &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member{
					State: &gostruct.OpenconfigNetworkInstance_NetworkInstances_NetworkInstance_Vlans_Vlan_Members_Member_State{
						Interface: ygot.String(strings.Split(intf, "Vlan"+vid+"|")[1]),
					},
				})
			}
			vlanMbs.Member = memberSlice
		}
		fmt.Println("Syncing Vlan...")
		time.Sleep(5 * time.Second)
	}
	return nil
}

func SetVlan(key []string, value []string, str string, b bool) error {
	cmd := "sudo config vlan add " + str

	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}
	return nil
}

func DelVlan(key []string, value []string, str string, b bool) error {
	dbAddr := "localhost:6379"
	dbPass := ""

	applClient := redis.NewClient(&redis.Options{
		Addr:     dbAddr,
		Password: dbPass, // no password set
		DB:       0,      // use appl DB
		PoolSize: 10,
	})

	keys := applClient.Keys("VLAN_MEMBER_TABLE:Vlan" + value[1] + ":*").Val()
	for _, member := range keys {
		cmd := "sudo config vlan member del " + value[1] + " " + strings.Split(member, ":")[2]

		_, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			return status.Error(codes.Internal, "Failed to execute command: "+cmd)
		}

		cmd = "sudo config vlan member add " + DEFAULT_VLAN + " " + strings.Split(member, ":")[2]

		_, err = exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			return status.Error(codes.Internal, "Failed to execute command: "+cmd)
		}
	}

	cmd := "sudo config vlan del " + value[1]

	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}
	return nil
}

//Set interface as trunk port. The function will delete all other untagged vlan.
func SetVlanMember(key []string, value []string, str string, b bool) error {
	dbAddr := "localhost:6379"
	dbPass := ""

	applClient := redis.NewClient(&redis.Options{
		Addr:     dbAddr,
		Password: dbPass, // no password set
		DB:       0,      // use appl DB
		PoolSize: 10,
	})

	keys := applClient.Keys("VLAN_MEMBER_TABLE:*:" + str).Val()
	if len(keys) == 0 {

	} else {
		for _, vlan := range keys {
			if applClient.HGet(vlan, "tagging_mode").Val() == "untagged" {
				vlanId := strings.Split(strings.Split(vlan, "Vlan")[1], ":")[0]
				cmd := "sudo config vlan member del " + vlanId + " " + str
				fmt.Println(cmd)
				_, err := exec.Command("bash", "-c", cmd).Output()
				if err != nil {
					return status.Error(codes.Internal, "Failed to execute command: "+cmd)
				}
			}
		}
	}

	cmd := "sudo config vlan member add " + value[1] + " " + str
	fmt.Println(cmd)

	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}
	return nil
}

func SetInterfaceAccessVlan(key []string, value []string, str string, b bool) error {
	dbAddr := "localhost:6379"
	dbPass := ""

	applClient := redis.NewClient(&redis.Options{
		Addr:     dbAddr,
		Password: dbPass, // no password set
		DB:       0,      // use appl DB
		PoolSize: 10,
	})

	keys := applClient.Keys("VLAN_MEMBER_TABLE:*:" + value[0]).Val()
	if len(keys) == 0 {

	} else {
		for _, vlan := range keys {
			vlanId := strings.Split(strings.Split(vlan, "Vlan")[1], ":")[0]
			cmd := "sudo config vlan member del " + vlanId + " " + value[0]
			fmt.Println(cmd)
			_, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil {
				return status.Error(codes.Internal, "Failed to execute command: "+cmd)
			}
		}
	}

	cmd := "sudo config vlan member add -u " + str + " " + value[0]
	fmt.Println(cmd)
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}
	return nil
}

func DelVlanMember(key []string, value []string, str string, b bool) error {
	cmd := "sudo config vlan member del " + value[1] + " " + value[2]

	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}

	cmd = "sudo config vlan member add -u " + DEFAULT_VLAN + " " + value[2]

	_, err = exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return status.Error(codes.Internal, "Failed to execute command: "+cmd)
	}
	return nil
}
