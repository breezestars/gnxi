package components_funcs

import (
	"fmt"
	"time"
	"strings"
	"github.com/openconfig/ygot/ygot"
	"github.com/breezestars/gnxi/gnmi/modeldata/gostruct"
	"os/exec"
)

func InitPlatform(device *gostruct.Device) error {

	t0 := time.Now()
	cmd := "show platform syseeprom | grep '0x' | awk -F'0x' '{print $2,$3}' | awk -F' ' '{print \"0x\"$1,$3}'"

	syseeprom, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmd)
	}
	t0d := time.Since(t0)
	t1 := time.Now()

	//Manufacture Date     0x25
	//Label Revision       0x27
	//Platform Name        0x28
	//ONIE Version         0x29
	//Manufacturer         0x2B
	//Manufacture Country  0x2C
	//Diag Version         0x2E
	//Base MAC Address     0x24
	//Serial Number        0x23
	//Part Number          0x22
	//Product Name         0x21
	//MAC Addresses        0x2A
	//Vendor Name          0x2D
	//CRC-32               0xFE
	//(checksum valid)`

	//syseeprom := `0x25 05/25/2017
	//0x27 R02B
	//0x28 x86_64-accton_as7712_32x-r0
	//0x29 2016.08.00.03
	//0x2B Accton
	//0x2C TW
	//0x2E 0.0.5.3
	//0x24 A8:2B:B5:38:01:1D
	//0x23 771232X1721087
	//0x22 FP3ZZ7632014A
	//0x21 7712-32X-O-AC-F
	//0x2A 131
	//0x2D Edgecore
	//0xFE CB07A5F3`

	cmd = "show version | grep 'Software Version'"
	versionOutput, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("Failed to execute command: %s", cmd)
	}
	t1d := time.Since(t1)
	t2 := time.Now()
	//versionOutput := `SONiC Software Version: SONiC.HEAD.638-a0bd656`

	consoleArray := strings.Split(string(syseeprom), "\n")

	product, err := search(consoleArray, "0x21")
	hwVersion, err := search(consoleArray, "0x28")
	serialNo, err := search(consoleArray, "0x23")
	partNo, err := search(consoleArray, "0x22")
	mfgName, err := search(consoleArray, "0x2B")
	mfgDate, err := search(consoleArray, "0x25")
	splitMfgDate := strings.Split(mfgDate, "/")
	swVersion := strings.Split(string(versionOutput), ": ")[1]

	device.Components = &gostruct.OpenconfigPlatform_Components{}

	compnent, err := device.Components.NewComponent(product)
	if err != nil {
		return err
	}

	compnent.Name = ygot.String(product)
	compnent.Config = &gostruct.OpenconfigPlatform_Components_Component_Config{
		Name: ygot.String(product),
	}
	compnent.State = &gostruct.OpenconfigPlatform_Components_Component_State{
		SerialNo:        ygot.String(serialNo),
		PartNo:          ygot.String(partNo),
		HardwareVersion: ygot.String(hwVersion),
		SoftwareVersion: ygot.String(swVersion),
		MfgName:         ygot.String(mfgName),
		MfgDate:         ygot.String(splitMfgDate[2] + "-" + splitMfgDate[0] + "-" + splitMfgDate[1]),
		Type: &gostruct.OpenconfigPlatform_Components_Component_State_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT{
			gostruct.OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT_CHASSIS,
		},
	}

	compnent.State.SerialNo = ygot.String(serialNo)

	t2d := time.Since(t2)

	fmt.Println("Parsed all components information, the timeing information is following")
	fmt.Printf("Get 1 command, took %s \n", t0d)
	fmt.Printf("Get 2 command, took %s \n", t1d)
	fmt.Printf("Parse all data, took %s \n", t2d)

	return nil
}

func search(s []string, tgt string) (string, error) {
	for _, c := range s {
		if strings.Contains(c, tgt) {
			return strings.Split(c, " ")[1], nil
		}
	}
	return "", fmt.Errorf("Not found")
}
