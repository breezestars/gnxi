package utils

import (
	"github.com/openconfig/ygot/ygot"
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/breezestars/gnxi/utils/components_funcs"
	"github.com/go-redis/redis"
	"fmt"
	"time"
)

func InitGoStruct() (ygot.ValidatedGoStruct, error) {

	configClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       4,  // use config DB
		PoolSize: 10,
	})

	_, err := configClient.Ping().Result()
	if err != nil {
		return nil, err
	} else {
		fmt.Println("ConfigDB: Connect to localhost:6379 success.")
	}

	device := &gostruct.Device{}

	err = components_funcs.InitPlatform(device)
	if err != nil {
		return nil, err
	}

	tLldp := time.Now()
	err = components_funcs.InitLldp(device)
	if err != nil {
		return nil, err
	}
	tLldpD := time.Since(tLldp)
	fmt.Printf("=== Init LLDP, took %s === \n", tLldpD)

	err = components_funcs.InitInterface(device)
	if err != nil {
		return nil, err
	}

	go components_funcs.SyncInterface(device)

	tVlan := time.Now()
	err = components_funcs.InitVlan(device, configClient)
	if err != nil {
		return nil, err
	}
	tVlanD := time.Since(tVlan)
	fmt.Printf("=== Init Vlan, took %s === \n", tVlanD)

	go components_funcs.SyncVlan(device, configClient)

	return device, nil
}
