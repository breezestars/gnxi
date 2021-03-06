/*
Package gostruct is a generated package which contains definitions
of structs which represent a YANG schema. The generated schema can be
compressed by a series of transformations (compression was false
in this case).

This package was generated by /Users/cmo/.gvm/pkgsets/go1.10.2/global/src/github.com/openconfig/ygot/ygen/commongen.go
using the following YANG input files:
	- github.com/openconfig/public/release/models/interfaces/openconfig-if-aggregate.yang
	- github.com/openconfig/public/release/models/interfaces/openconfig-if-ethernet.yang
	- github.com/openconfig/public/release/models/interfaces/openconfig-interfaces.yang
	- github.com/openconfig/public/release/models/lacp/openconfig-lacp.yang
	- github.com/openconfig/public/release/models/lldp/openconfig-lldp.yang
	- github.com/openconfig/public/release/models/network-instance/openconfig-network-instance.yang
	- github.com/openconfig/public/release/models/platform/openconfig-platform.yang
	- github.com/openconfig/public/release/models/system/openconfig-system.yang
	- github.com/openconfig/public/release/models/vlan/openconfig-vlan.yang
Imported modules were sourced from:
	- github.com/openconfig/public/...
	- github.com/YangModels/yang/...
*/
package gostruct

// Device represents the /device YANG schema element.
type Device struct {
	Acl              *OpenconfigAcl_Acl                          `path:"acl" module:"openconfig-acl"`
	Bgp              *OpenconfigBgp_Bgp                          `path:"bgp" module:"openconfig-bgp"`
	Components       *OpenconfigPlatform_Components              `path:"components" module:"openconfig-platform"`
	Interfaces       *OpenconfigInterfaces_Interfaces            `path:"interfaces" module:"openconfig-interfaces"`
	Lacp             *OpenconfigLacp_Lacp                        `path:"lacp" module:"openconfig-lacp"`
	Lldp             *OpenconfigLldp_Lldp                        `path:"lldp" module:"openconfig-lldp"`
	LocalRoutes      *OpenconfigLocalRouting_LocalRoutes         `path:"local-routes" module:"openconfig-local-routing"`
	Mpls             *OpenconfigMpls_Mpls                        `path:"mpls" module:"openconfig-mpls"`
	NetworkInstances *OpenconfigNetworkInstance_NetworkInstances `path:"network-instances" module:"openconfig-network-instance"`
	RoutingPolicy    *OpenconfigRoutingPolicy_RoutingPolicy      `path:"routing-policy" module:"openconfig-routing-policy"`
	System           *OpenconfigSystem_System                    `path:"system" module:"openconfig-system"`
}

// IsYANGGoStruct ensures that Device implements the yang.GoStruct
// interface. This allows functions that need to handle this struct to
// identify it as being generated by ygen.
func (*Device) IsYANGGoStruct() {}
