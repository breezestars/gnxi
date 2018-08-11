package main

import (
	"github.com/tidwall/gjson"
	)

func main() {
	json := `{ "lldp": { "interface": { "eth0": { "age": "0 day, 08:41:34", "via": "LLDP", "rid": "1", "chassis": { "B2-Tor": { "id": { "type": "mac", "value": "cc:37:ab:60:cd:c1" }, "capability": [ { "type": "Bridge", "enabled": true }, { "type": "Router", "enabled": true } ], "mgmt-ip": "0.0.0.0", "descr": "PICA8 Inc., Model as4610_54, PicOS 2.6.2", "ttl": "120" } }, "port": { "id": { "type": "ifname", "value": "ge-1/1/28" }, "descr": "ge-1/1/28", "auto-negotiation": { "supported": true, "enabled": true, "advertised": [ { "type": "10Base-T", "hd": true, "fd": false }, { "type": "100Base-TX", "hd": true, "fd": true }, { "type": "1000Base-T", "hd": true, "fd": true } ], "current": "1000BaseTFD - Four-pair Category 5 UTP, full duplex mode" } }, "vlan": { "vlan-id": "61", "pvid": true } } } } }`

	gjson.Get(json, "lldp.interface.eth0.")
	//path:="lldp.interface.eth0.port.auto-negotiation.advertised"
	//fmt.Println(gjson.Get(json, path).Index)
	//fmt.Println(gjson.Get(json, path).Num)
	//fmt.Println(gjson.Get(json, path).Type)
	//fmt.Println(gjson.Get(json, path).IsArray())
	//fmt.Println(gjson.Get(json, path).IsObject())
	//fmt.Println(gjson.Get(json, path).Map())
	//fmt.Println(gjson.Get(json, path).Raw)
	//intfs := gjson.Get(json, path).Array()
	//for k, v := range intfs {
	//	fmt.Println(k, ", ", v)
	//}

}
