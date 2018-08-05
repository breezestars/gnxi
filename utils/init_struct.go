package utils

import (
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"github.com/breezestars/gnxi/utils/components_funcs"
	"github.com/go-redis/redis"
	"fmt"
	"time"
	"sync"
)

func InitGoStruct(mu *sync.RWMutex) (*gostruct.Device, error) {

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

	tComponent := time.Now()
	err = components_funcs.InitPlatform(device)
	if err != nil {
		return nil, err
	}
	tComponentD := time.Since(tComponent)
	fmt.Printf("=== Init Components, took %s === \n", tComponentD)

	tLldp := time.Now()
	err = components_funcs.InitLldp(device)
	if err != nil {
		return nil, err
	}
	tLldpD := time.Since(tLldp)
	fmt.Printf("=== Init LLDP, took %s === \n", tLldpD)

	tInterface := time.Now()
	err = components_funcs.InitInterface(device)
	if err != nil {
		return nil, err
	}
	tInterfaceD := time.Since(tInterface)
	fmt.Printf("=== Init Interface, took %s === \n", tInterfaceD)

	go components_funcs.SyncInterface(device, mu)

	tVlan := time.Now()
	err = components_funcs.InitVlan(device, configClient)
	if err != nil {
		return nil, err
	}
	tVlanD := time.Since(tVlan)
	fmt.Printf("=== Init Vlan, took %s === \n", tVlanD)

	go components_funcs.SyncVlan(device, configClient, mu)

	return device, nil
}
