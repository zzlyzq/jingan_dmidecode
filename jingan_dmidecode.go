package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"os"
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

type IPMIInfo struct {
	IP          string       `json:"ip"`
	MAC         string       `json:"mac"`
	Users       []IPMIUser   `json:"users"`
}

type IPMIUser struct {
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	PrivilegeLevel string `json:"privilege_level"`
	Enable         string `json:"enable"`
}


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
	IPMI                IPMIInfo               `json:"ipmi"`
}


func runCommand(cmdName string, args ...string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	fmt.Printf("Command Output: %s\n", string(output)) // 添加调试信息
	return strings.ReplaceAll(string(output), "\r\n", "\n"), nil
}

func parseDmidecode(output string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func getStringValue(values map[string]interface{}, key string) string {
	if val, ok := values[key].(string); ok {
		return val
	}
	return ""
}

func getMapValue(entry map[string]interface{}, key string) (map[string]interface{}, bool) {
	val, ok := entry[key].(map[string]interface{})
	return val, ok
}

func main() {
	//var dmidecodeCmd string
	var output string
	var err error
	output = ""
	if runtime.GOOS == "windows" {
		output, err = runCommand("cmd.exe", "/c", "dmidecode.exe | jc.exe --dmidecode -q")
		if err != nil {
			fmt.Printf("Unable to get dmidecode information. Error: %v\n", err)
			return
		}
	} else {
		output, err = runCommand("sh", "-c", "./dmidecode | ./jc --dmidecode -q")
		if err != nil {
			fmt.Printf("Unable to get dmidecode information. Error: %v\n", err)
			return
		}
	}

	dmiData, err := parseDmidecode(output)
	if err != nil {
		fmt.Printf("Failed to parse dmidecode output. Error: %v\n", err)
		return
	}

	var systemInfo FullSystemInfo

	for _, entry := range dmiData {
		description := getStringValue(entry, "description")
		values, ok := getMapValue(entry, "values")
		if !ok {
			continue
		}

		switch description {
		case "BIOS Information":
			systemInfo.BIOS = BIOSInfo{
				Vendor:      getStringValue(values, "vendor"),
				Version:     getStringValue(values, "version"),
				ReleaseDate: getStringValue(values, "release_date"),
			}
		case "System Information":
			systemInfo.System = SysInfo{
				Manufacturer: getStringValue(values, "manufacturer"),
				ProductName:  getStringValue(values, "product_name"),
				Version:      getStringValue(values, "version"),
				SerialNumber: getStringValue(values, "serial_number"),
				UUID:         getStringValue(values, "uuid"),
				Family:       getStringValue(values, "family"),
			}
		case "Base Board Information":
			systemInfo.Baseboard = BaseboardInfo{
				Manufacturer: getStringValue(values, "manufacturer"),
				ProductName:  getStringValue(values, "product_name"),
				Version:      getStringValue(values, "version"),
				SerialNumber: getStringValue(values, "serial_number"),
				AssetTag:     getStringValue(values, "asset_tag"),
			}
		case "Chassis Information":
			systemInfo.Chassis = ChassisInfo{
				Manufacturer: getStringValue(values, "manufacturer"),
				Type:         getStringValue(values, "type"),
				Version:      getStringValue(values, "version"),
				SerialNumber: getStringValue(values, "serial_number"),
				AssetTag:     getStringValue(values, "asset_tag"),
			}
		case "Processor Information":
			processor := ProcessorInfo{
				Family:       getStringValue(values, "family"),
				Manufacturer: getStringValue(values, "manufacturer"),
				Version:      getStringValue(values, "version"),
				Frequency:    getStringValue(values, "current_speed"),
				Cores:        getStringValue(values, "core_count"),
				Threads:      getStringValue(values, "thread_count"),
			}
			systemInfo.Processors = append(systemInfo.Processors, processor)
		case "Memory Device":
			memory := MemoryInfo{
				Model:     getStringValue(values, "part_number"),
				Size:      getStringValue(values, "size"),
				Speed:     getStringValue(values, "speed"),
				Slot:      getStringValue(values, "locator"),
			}
			systemInfo.Memory = append(systemInfo.Memory, memory)
		case "System Slot Information":
			slot := SystemSlotInfo{
				Type:        getStringValue(values, "type"),
				Usage:       getStringValue(values, "current_usage"),
				Status:      getStringValue(values, "status"),
				BusAddress:  getStringValue(values, "bus_address"),
				Designation: getStringValue(values, "designation"),
			}
			systemInfo.SystemSlots = append(systemInfo.SystemSlots, slot)
		case "System Power Supply":
			powerSupply := SystemPowerSupplyInfo{
				PowerUnitGroup:     getStringValue(values, "power_unit_group"),
				Location:           getStringValue(values, "location"),
				Name:               getStringValue(values, "name"),
				Manufacturer:       getStringValue(values, "manufacturer"),
				AssetTag:           getStringValue(values, "asset_tag"),
				ModelPartNumber:    getStringValue(values, "model_part_number"),
				MaxPowerCapacity:   getStringValue(values, "max_power_capacity"),
				Status:             getStringValue(values, "status"),
				Type:               getStringValue(values, "type"),
				Plugged:            getStringValue(values, "plugged"),
				HotReplaceable:     getStringValue(values, "hot_replaceable"),
				CoolingDeviceHandle: getStringValue(values, "cooling_device_handle"),
			}
			systemInfo.SystemPowerSupplies = append(systemInfo.SystemPowerSupplies, powerSupply)
		case "Onboard Device":
			onboardDevice := OnboardDeviceInfo{
				ReferenceDesignation: getStringValue(values, "reference_designation"),
				Type:                 getStringValue(values, "type"),
				Status:               getStringValue(values, "status"),
				TypeInstance:         getStringValue(values, "type_instance"),
			}
			systemInfo.OnboardDevices = append(systemInfo.OnboardDevices, onboardDevice)
		}
	}

	if err := getRaidInfo(&systemInfo); err != nil {
		fmt.Printf("Failed to get RAID information: %v\n", err)
	}

	if err := getIPMIInfo(&systemInfo); err != nil {
		fmt.Printf("Failed to get IPMI information: %v\n", err)
	}

	jsonData, err := json.MarshalIndent(systemInfo, "", "  ")
	if err != nil {
		log.Fatalf("Failed to convert system information to JSON: %v", err)
	}


	fmt.Println(string(jsonData))

	// 创建文件名
	filename := fmt.Sprintf("%s_%s_%s.txt",
		strings.ReplaceAll(systemInfo.System.Manufacturer, " ", "_"),
		strings.ReplaceAll(systemInfo.System.ProductName, " ", "_"),
		strings.ReplaceAll(systemInfo.System.SerialNumber, " ", "_"))

	// 写入到文件
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write JSON data to file: %v", err)
	}

	fmt.Printf("System information has been written to %s\n", filename)
}


// getRaidInfo uses storcli to gather RAID and disk information
func getRaidInfo(systemInfo *FullSystemInfo) error {
	storcliCmd := "./storcli"
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

func getIPMIInfo(systemInfo *FullSystemInfo) error {
	ipmiCmd, err := getIPMICmdPath()
	if err != nil {
		return fmt.Errorf("failed to get IPMI command path: %v", err)
	}

	// Get IP and MAC address
	ipOutput, err := runCommand(ipmiCmd, "-m")
	if err != nil {
		return fmt.Errorf("failed to get IPMI IP and MAC address: %v", err)
	}
	fmt.Printf("IPMI IP and MAC output: %s\n", ipOutput) // 添加调试信息

	ipmiInfo := IPMIInfo{}
	for _, line := range strings.Split(ipOutput, "\n") {
		if strings.HasPrefix(line, "IP=") {
			ipmiInfo.IP = strings.TrimSpace(strings.Split(line, "=")[1])
		} else if strings.HasPrefix(line, "MAC=") {
			ipmiInfo.MAC = strings.TrimSpace(strings.Split(line, "=")[1])
		}
	}

	// Get user list
	userOutput, err := runCommand(ipmiCmd, "-user", "list")
	if err != nil {
		return fmt.Errorf("failed to get IPMI user list: %v", err)
	}
	fmt.Printf("IPMI User List output: %s\n", userOutput) // 添加调试信息

	users := []IPMIUser{}
	for _, line := range strings.Split(userOutput, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "User ID") || strings.HasPrefix(line, "-------") || strings.HasPrefix(line, "Maximum") || strings.HasPrefix(line, "Count") || line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 4)
		if len(parts) != 4 {
			continue
		}
		user := IPMIUser{
			UserID:         strings.TrimSpace(parts[0]),
			UserName:       strings.TrimSpace(parts[1]),
			PrivilegeLevel: strings.TrimSpace(parts[2]),
			Enable:         strings.TrimSpace(parts[3]),
		}
		users = append(users, user)
	}
	ipmiInfo.Users = users
	systemInfo.IPMI = ipmiInfo

	return nil
}

// 获取当前工作目录并构建ipmicfg.exe的绝对路径
func getIPMICmdPath() (string, error) {
	if runtime.GOOS == "windows" {
		wd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s\\ipmicfg.exe", wd), nil
	}
	return "./ipmicfg", nil
}

