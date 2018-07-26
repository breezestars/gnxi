package utils

import (
	"github.com/openconfig/ygot/ygot"
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/breezestars/gnxi/utils/components_funcs"
)

func InitGoStruct() (ygot.ValidatedGoStruct, error) {
	device := &gostruct.Device{}

	err := components_funcs.InitPlatform(device)
	if err != nil {
		return nil, err
	}

	err = components_funcs.InitLldp(device)
	if err != nil {
		return nil, err
	}

	err = components_funcs.InitInterface(device)
	if err != nil {
		return nil, err
	}

	go components_funcs.SyncInterface(device)

	err = components_funcs.InitVlan(device)
	if err != nil {
		return nil, err
	}
	return device, nil
}



