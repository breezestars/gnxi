package gostruct

//go:generate sh -c "go get -u github.com/openconfig/ygot; (cd $GOPATH/src/github.com/openconfig/ygot && go get -t -d ./...); go get -u github.com/openconfig/public; go get -u github.com/YangModels/yang; cd $GOPATH/src && go run github.com/openconfig/ygot/generator/generator.go -generate_fakeroot -output_file github.com/breezestars/gnxi/gnmi/modeldata/gostruct/generated.go -package_name gostruct -exclude_modules ietf-interfaces -path github.com/openconfig/public,github.com/YangModels/yang github.com/openconfig/public/release/models/interfaces/openconfig-interfaces.yang github.com/openconfig/public/release/models/openflow/openconfig-openflow.yang github.com/openconfig/public/release/models/platform/openconfig-platform.yang github.com/openconfig/public/release/models/system/openconfig-system.yang"
