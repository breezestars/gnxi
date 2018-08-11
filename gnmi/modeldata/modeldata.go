/* Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package modeldata contains the following model data in gnmi proto struct:
//	openconfig-if-aggregate 2.3.1,
//	openconfig-if-ethernet 2.4.0,
//	openconfig-interfaces 2.3.1,
//	openconfig-lldp 0.1.0,
//	openconfig-platform 0.11.0,
//	openconfig-system 0.5.0,
//	openconfig-vlan 3.0.1,
package modeldata

import (
	pb "github.com/openconfig/gnmi/proto/gnmi"
)

const (
	// OpenconfigInterfaceAggregateModel is the openconfig YANG model for interface-if-aggregate.
	OpenconfigInterfaceAggregateModel = "openconfig-if-aggregate"
	// OpenconfigInterfaceEthernetModel is the openconfig YANG model for interface-if-ethernet.
	OpenconfigInterfaceEthernetModel = "openconfig-if-ethernet"
	// OpenconfigInterfacesModel is the openconfig YANG model for interfaces.
	OpenconfigInterfacesModel = "openconfig-interfaces"
	// OpenconfigLacpModel is the openconfig YANG model for lacp.
	OpenconfigLacpModel = "openconfig-lacp"
	// OpenconfigLldpModel is the openconfig YANG model for lldp.
	OpenconfigLldpModel = "openconfig-lldp"
	// OpenconfigNetworkInstanceModel is the openconfig YANG model for network-instance.
	OpenconfigNetworkInstanceModel = "openconfig-network-instance"
	// OpenconfigPlatformModel is the openconfig YANG model for platform.
	OpenconfigPlatformModel = "openconfig-platform"
	// OpenconfigSystemModel is the openconfig YANG model for system.
	OpenconfigSystemModel = "openconfig-system"
	// OpenconfigVlanModel is the openconfig YANG model for vlan.
	OpenconfigVlanModel = "openconfig-vlan"
)

var (
	// ModelData is a list of supported models.
	ModelData = []*pb.ModelData{{
		Name:         OpenconfigInterfaceAggregateModel,
		Organization: "OpenConfig working group",
		Version:      "2.3.1",
	}, {
		Name:         OpenconfigInterfaceEthernetModel,
		Organization: "OpenConfig working group",
		Version:      "2.4.0",
	}, {
		Name:         OpenconfigInterfacesModel,
		Organization: "OpenConfig working group",
		Version:      "2.3.1",
	}, {	Name:         OpenconfigLacpModel,
		Organization: "OpenConfig working group",
		Version:      "1.1.0",
	}, {
		Name:         OpenconfigLldpModel,
		Organization: "OpenConfig working group",
		Version:      "0.1.0",
	}, {
		Name:         OpenconfigNetworkInstanceModel,
		Organization: "OpenConfig working group",
		Version:      "0.10.2",
	}, {
		Name:         OpenconfigPlatformModel,
		Organization: "OpenConfig working group",
		Version:      "0.11.0",
	}, {
		Name:         OpenconfigSystemModel,
		Organization: "OpenConfig working group",
		Version:      "0.5.0",
	}, {
		Name:         OpenconfigVlanModel,
		Organization: "OpenConfig working group",
		Version:      "3.0.1",
	},}
)
