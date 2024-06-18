package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// BIOSInfo struct to hold BIOS information
type BIOSInfo struct {
	Vendor      string `json:"vendor"`
	Version     string `json:"version"`
	ReleaseDate string `json:"release_date"`
}

// SysInfo struct to hold System information
type SysInfo struct {
	Manufacturer  string `json:"manufacturer"`
	ProductName   string `json:"product_name"`
	Version       string `json:"version"`
	SerialNumber  string `json:"serial_number"`
	UUID          string `json:"uuid"`
	Family        string `json:"family"`
}

// BaseboardInfo struct to hold Baseboard information
type BaseboardInfo struct {
	Manufacturer string `json:"manufacturer"`
	ProductName  string `json:"product_name"`
	Version      string `json:"version"`
	SerialNumber string `json:"serial_number"`
	AssetTag     string `json:"asset_tag"`
}

// ChassisInfo struct to hold Chassis information
type ChassisInfo struct {
	Manufacturer string `json:"manufacturer"`
	Type         string `json:"type"`
	Version      string `json:"version"`
	SerialNumber string `json:"serial_number"`
	AssetTag     string `json:"asset_tag"`
}

// ProcessorInfo struct to hold Processor information
type ProcessorInfo struct {
	Family       string `json:"family"`
	Manufacturer string `json:"manufacturer"`
	Version      string `json:"version"`
	Frequency    string `json:"frequency"`
	Cores        string `json:"cores"`
	Threads      string `json:"threads"`
}

// MemoryInfo struct to hold Memory information
type MemoryInfo struct {
	Model     string `json:"model"`
	Size      string `json:"size"`
	Speed     string `json:"speed"`
	Slot      string `json:"slot"`
}

// SystemSlotInfo struct to hold System Slot information
type SystemSlotInfo struct {
	Type        string `json:"type"`
	Usage       string `json:"usage"`
	Status      string `json:"status"`
	BusAddress  string `json:"bus_address"`
	Designation string `json:"designation"`
}

// SystemPowerSupplyInfo struct to hold System Power Supply information
type SystemPowerSupplyInfo struct {
	PowerUnitGroup     string `json:"power_unit_group"`
	Location           string `json:"location"`
	Name               string `json:"name"`
	Manufacturer       string `json:"manufacturer"`
	AssetTag           string `json:"asset_tag"`
	ModelPartNumber    string `json:"model_part_number"`
	MaxPowerCapacity   string `json:"max_power_capacity"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	Plugged            string `json:"plugged"`
	HotReplaceable     string `json:"hot_replaceable"`
	CoolingDeviceHandle string `json:"cooling_device_handle"`
}

// OnboardDeviceInfo struct to hold Onboard Device information
type OnboardDeviceInfo struct {
	ReferenceDesignation string `json:"reference_designation"`
	Type                 string `json:"type"`
	Status               string `json:"status"`
	TypeInstance         string `json:"type_instance"`
}

// RaidCardInfo struct to hold RAID card information
type RaidCardInfo struct {
	ControllerID string      `json:"controller_id"`
	Model        string      `json:"model"`
	Ports        string      `json:"ports"`
	PDs          string      `json:"pds"`
	DGs          string      `json:"dgs"`
	DNOpt        string      `json:"dnopt"`
	VDs          string      `json:"vds"`
	VNOpt        string      `json:"vnopt"`
	BBU          string      `json:"bbu"`
	SPR          string      `json:"spr"`
	DS           string      `json:"ds"`
	EHS          string      `json:"ehs"`
	ASOs         string      `json:"asos"`
	Hlth         string      `json:"hlth"`
	VirtualDrives []VirtualDriveInfo `json:"virtual_drives"`
	PhysicalDrives []PhysicalDriveInfo `json:"physical_drives"`
}

// VirtualDriveInfo struct to hold Virtual Drive information
type VirtualDriveInfo struct {
	DGVD     string `json:"dg_vd"`
	Type     string `json:"type"`
	State    string `json:"state"`
	Access   string `json:"access"`
	Consist  string `json:"consist"`
	Cache    string `json:"cache"`
	Cac      string `json:"cac"`
	SCC      string `json:"scc"`
	Size     string `json:"size"`
	Name     string `json:"name"`
}

// PhysicalDriveInfo struct to hold Physical Drive information
type PhysicalDriveInfo struct {
	EIDSlt  string `json:"eid_slt"`
	DID     string `json:"did"`
	State   string `json:"state"`
	DG      string `json:"dg"`
	Size    string `json:"size"`
	Intf    string `json:"intf"`
	Med     string `json:"med"`
	SED     string `json:"sed"`
	PI      string `json:"pi"`
	SeSz    string `json:"se_sz"`
	Model   string `json:"model"`
	Sp      string `json:"sp"`
	Type    string `json:"type"`
}

// FullSystemInfo struct to hold all information
type FullSystemInfo struct {
	BIOS                BIOSInfo               `json:"bios"`
	System              SysInfo                `json:"system"`
	Baseboard           BaseboardInfo          `json:"baseboard"`
	Chassis             ChassisInfo            `json:"chassis"`
	Processors          []ProcessorInfo        `json:"processors"`
	Memory              []MemoryInfo           `json:"memory"`
	SystemSlots         []SystemSlotInfo       `json:"system_slots"`
	SystemPowerSupplies []SystemPowerSupplyInfo `json:"system_power_supplies"`
	OnboardDevices      []OnboardDeviceInfo    `json:"onboard_devices"`
	RaidCardInfo        RaidCardInfo           `json:"raid_card_info"`
}

// runCommand runs the specified command with the given arguments and returns the output
func runCommand(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(output), "\r\n", "\n"), nil
}

func parseDmidecode(output string) map[string]map[string]string {
	result := make(map[string]map[string]string)
	var currentSection string
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			currentSection = ""
		} else if !strings.HasPrefix(line, "\t") {
			currentSection = strings.TrimSpace(line)
			result[currentSection] = make(map[string]string)
		} else if currentSection != "" {
			parts := strings.SplitN(strings.TrimSpace(line), ":", 2)
			if len(parts) == 2 {
				result[currentSection][strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}
	}
	return result
}

func main() {
	var dmidecodeCmd string
	if runtime.GOOS == "windows" {
		dmidecodeCmd = "dmidecode.exe"
	} else {
		dmidecodeCmd = "dmidecode"
	}

	output, err := runCommand(dmidecodeCmd)
	if err != nil {
		fmt.Printf("Unable to get dmidecode information. Error: %v\n", err)
		return
	}

	dmiData := parseDmidecode(output)
	var systemInfo FullSystemInfo

	if biosData, ok := dmiData["BIOS Information"]; ok {
		systemInfo.BIOS = BIOSInfo{
			Vendor:      biosData["Vendor"],
			Version:     biosData["Version"],
			ReleaseDate: biosData["Release Date"],
		}
	}

	if sysData, ok := dmiData["System Information"]; ok {
		systemInfo.System = SysInfo{
			Manufacturer: sysData["Manufacturer"],
			ProductName:  sysData["Product Name"],
			Version:      sysData["Version"],
			SerialNumber: sysData["Serial Number"],
			UUID:         sysData["UUID"],
			Family:       sysData["Family"],
		}
	}

	if baseboardData, ok := dmiData["Base Board Information"]; ok {
		systemInfo.Baseboard = BaseboardInfo{
			Manufacturer: baseboardData["Manufacturer"],
			ProductName:  baseboardData["Product Name"],
			Version:      baseboardData["Version"],
			SerialNumber: baseboardData["Serial Number"],
			AssetTag:     baseboardData["Asset Tag"],
		}
	}

	if chassisData, ok := dmiData["Chassis Information"]; ok {
		systemInfo.Chassis = ChassisInfo{
			Manufacturer: chassisData["Manufacturer"],
			Type:         chassisData["Type"],
			Version:      chassisData["Version"],
			SerialNumber: chassisData["Serial Number"],
			AssetTag:     chassisData["Asset Tag"],
		}
	}

	for key, value := range dmiData {
		if strings.HasPrefix(key, "Processor Information") {
			processor := ProcessorInfo{
				Family:       value["Family"],
				Manufacturer: value["Manufacturer"],
				Version:      value["Version"],
				Frequency:    value["Max Speed"],
				Cores:        value["Core Count"],
				Threads:      value["Thread Count"],
			}
			systemInfo.Processors = append(systemInfo.Processors, processor)
		}
		if strings.HasPrefix(key, "Memory Device") {
			memory := MemoryInfo{
				Model:     value["Part Number"],
				Size:      value["Size"],
				Speed:     value["Speed"],
				Slot:      value["Locator"],
			}
			systemInfo.Memory = append(systemInfo.Memory, memory)
		}
		if strings.HasPrefix(key, "System Slot") {
			slot := SystemSlotInfo{
				Type:        value["Type"],
				Usage:       value["Current Usage"],
				Status:      value["Slot Status"],
				BusAddress:  value["Bus Address"],
				Designation: value["Designation"],
			}
			systemInfo.SystemSlots = append(systemInfo.SystemSlots, slot)
		}
		if strings.HasPrefix(key, "System Power Supply") {
			powerSupply := SystemPowerSupplyInfo{
				PowerUnitGroup:     value["Power Unit Group"],
				Location:           value["Location"],
				Name:               value["Name"],
				Manufacturer:       value["Manufacturer"],
				AssetTag:           value["Asset Tag"],
				ModelPartNumber:    value["Model Part Number"],
				MaxPowerCapacity:   value["Max Power Capacity"],
				Status:             value["Status"],
				Type:               value["Type"],
				Plugged:            value["Plugged"],
				HotReplaceable:     value["Hot Replaceable"],
				CoolingDeviceHandle: value["Cooling Device Handle"],
			}
			systemInfo.SystemPowerSupplies = append(systemInfo.SystemPowerSupplies, powerSupply)
		}
		if strings.HasPrefix(key, "On Board Device") {
			onboardDevice := OnboardDeviceInfo{
				ReferenceDesignation: value["Reference Designation"],
				Type:                 value["Type"],
				Status:               value["Status"],
				TypeInstance:         value["Type Instance"],
			}
			systemInfo.OnboardDevices = append(systemInfo.OnboardDevices, onboardDevice)
		}
	}

	if err := getRaidInfo(&systemInfo); err != nil {
		fmt.Printf("Failed to get RAID information: %v\n", err)
	}

	jsonData, err := json.MarshalIndent(systemInfo, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert system information to JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}

// getRaidInfo uses storcli to gather RAID and disk information
func getRaidInfo(systemInfo *FullSystemInfo) error {
	storcliCmd := "storcli"
	if runtime.GOOS == "windows" {
		storcliCmd = "storcli.exe"
	}

	ctrlOutput, err := runCommand(storcliCmd, "show", "all", "J")
	if err != nil {
		return fmt.Errorf("failed to get controller info: %v", err)
	}

	var storcliOutput map[string]interface{}
	if err := json.Unmarshal([]byte(ctrlOutput), &storcliOutput); err != nil {
		return fmt.Errorf("failed to parse JSON output: %v", err)
	}

	controllers, ok := storcliOutput["Controllers"].([]interface{})
	if !ok || len(controllers) == 0 {
		return fmt.Errorf("no controllers found in output")
	}

	for _, ctrl := range controllers {
		controllerData, ok := ctrl.(map[string]interface{})
		if !ok {
			continue
		}

		responseData, ok := controllerData["Response Data"].(map[string]interface{})
		if !ok {
			continue
		}

		systemOverview, ok := responseData["System Overview"].([]interface{})
		if !ok || len(systemOverview) == 0 {
			continue
		}

		overview := systemOverview[0].(map[string]interface{})
		raidInfo := RaidCardInfo{
			ControllerID: fmt.Sprintf("%v", overview["Ctl"]),
			Model:        fmt.Sprintf("%v", overview["Model"]),
			Ports:        fmt.Sprintf("%v", overview["Ports"]),
			PDs:          fmt.Sprintf("%v", overview["PDs"]),
			DGs:          fmt.Sprintf("%v", overview["DGs"]),
			DNOpt:        fmt.Sprintf("%v", overview["DNOpt"]),
			VDs:          fmt.Sprintf("%v", overview["VDs"]),
			VNOpt:        fmt.Sprintf("%v", overview["VNOpt"]),
			BBU:          fmt.Sprintf("%v", overview["BBU"]),
			SPR:          fmt.Sprintf("%v", overview["sPR"]),
			DS:           fmt.Sprintf("%v", overview["DS"]),
			EHS:          fmt.Sprintf("%v", overview["EHS"]),
			ASOs:         fmt.Sprintf("%v", overview["ASOs"]),
			Hlth:         fmt.Sprintf("%v", overview["Hlth"]),
		}

		// Get detailed info for each controller
		ctrlDetails, err := runCommand(storcliCmd, fmt.Sprintf("/c%s", raidInfo.ControllerID), "show", "all", "J")
		if err != nil {
			fmt.Printf("Failed to get controller details for controller %s: %v\n", raidInfo.ControllerID, err)
			continue
		}

		var ctrlDetailOutput map[string]interface{}
		if err := json.Unmarshal([]byte(ctrlDetails), &ctrlDetailOutput); err != nil {
			fmt.Printf("Failed to parse JSON output for controller %s: %v\n", raidInfo.ControllerID, err)
			continue
		}

		ctrlResponseData, ok := ctrlDetailOutput["Controllers"].([]interface{})[0].(map[string]interface{})["Response Data"].(map[string]interface{})
		if !ok {
			continue
		}

		// Parse Virtual Drives (VD LIST)
		if vdList, ok := ctrlResponseData["VD LIST"].([]interface{}); ok {
			for _, vd := range vdList {
				vdInfo := vd.(map[string]interface{})
				virtualDrive := VirtualDriveInfo{
					DGVD:    fmt.Sprintf("%v", vdInfo["DG/VD"]),
					Type:    fmt.Sprintf("%v", vdInfo["TYPE"]),
					State:   fmt.Sprintf("%v", vdInfo["State"]),
					Access:  fmt.Sprintf("%v", vdInfo["Access"]),
					Consist: fmt.Sprintf("%v", vdInfo["Consist"]),
					Cache:   fmt.Sprintf("%v", vdInfo["Cache"]),
					Cac:     fmt.Sprintf("%v", vdInfo["Cac"]),
					SCC:     fmt.Sprintf("%v", vdInfo["sCC"]),
					Size:    fmt.Sprintf("%v", vdInfo["Size"]),
					Name:    fmt.Sprintf("%v", vdInfo["Name"]),
				}
				raidInfo.VirtualDrives = append(raidInfo.VirtualDrives, virtualDrive)
			}
		}

		// Parse Physical Drives (PD LIST)
		if pdList, ok := ctrlResponseData["PD LIST"].([]interface{}); ok {
			for _, pd := range pdList {
				pdInfo := pd.(map[string]interface{})
				physicalDrive := PhysicalDriveInfo{
					EIDSlt: fmt.Sprintf("%v", pdInfo["EID:Slt"]),
					DID:    fmt.Sprintf("%v", pdInfo["DID"]),
					State:  fmt.Sprintf("%v", pdInfo["State"]),
					DG:     fmt.Sprintf("%v", pdInfo["DG"]),
					Size:   fmt.Sprintf("%v", pdInfo["Size"]),
					Intf:   fmt.Sprintf("%v", pdInfo["Intf"]),
					Med:    fmt.Sprintf("%v", pdInfo["Med"]),
					SED:    fmt.Sprintf("%v", pdInfo["SED"]),
					PI:     fmt.Sprintf("%v", pdInfo["PI"]),
					SeSz:   fmt.Sprintf("%v", pdInfo["SeSz"]),
					Model:  fmt.Sprintf("%v", pdInfo["Model"]),
					Sp:     fmt.Sprintf("%v", pdInfo["Sp"]),
					Type:   fmt.Sprintf("%v", pdInfo["Type"]),
				}
				raidInfo.PhysicalDrives = append(raidInfo.PhysicalDrives, physicalDrive)
			}
		}

		systemInfo.RaidCardInfo = raidInfo
	}

	return nil
}
