package gostruct

//go:generate sh -c "go get -u github.com/openconfig/ygot; (cd $GOPATH/src/github.com/openconfig/ygot && go get -t -d ./...); go get -u github.com/openconfig/public; go get -u github.com/YangModels/yang; cd $GOPATH/src && go run github.com/openconfig/ygot/generator/generator.go -generate_fakeroot -compress_paths=false -output_dir github.com/breezestars/gnxi/gnmi/modeldata/gostruct/ -package_name gostruct -exclude_modules ietf-interfaces -generate_delete -generate_getters -path github.com/openconfig/public,github.com/YangModels/yang github.com/openconfig/public/release/models/interfaces/openconfig-if-aggregate.yang github.com/openconfig/public/release/models/interfaces/openconfig-if-ethernet.yang github.com/openconfig/public/release/models/interfaces/openconfig-interfaces.yang github.com/openconfig/public/release/models/lacp/openconfig-lacp.yang github.com/openconfig/public/release/models/lldp/openconfig-lldp.yang github.com/openconfig/public/release/models/network-instance/openconfig-network-instance.yang github.com/openconfig/public/release/models/platform/openconfig-platform.yang github.com/openconfig/public/release/models/system/openconfig-system.yang github.com/openconfig/public/release/models/vlan/openconfig-vlan.yang"
