package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type TargetPlatform struct {
	Os          string
	Arch        string
	Accelerator string
}

func (tp *TargetPlatform) Print() {
	fmt.Println("Target Platform")
	fmt.Printf("\tOs: %s\n", tp.Os)
	fmt.Printf("\tArchitecture: %s\n", tp.Arch)
	fmt.Printf("\tAccelerator: %s\n", tp.Accelerator)
}

func (tp *TargetPlatform) Validate() {
	if tp.Os != "linux" && tp.Os != "windows" {
		log.Fatal("Invalid Os!")
	}

	if tp.Os == "linux" {
		if tp.Arch != "armv8" && tp.Arch != "x64" {
			log.Fatal("Invalid architecture for Linux.")
		}
	}
	if tp.Os == "windows" {
		if tp.Arch != "x86" && tp.Arch != "x64" {
			log.Fatal("Invalid architecture for Windows.")
		}
	}
}

type CliArgs struct {
	DeviceFleet       string
	DeviceName        string
	IotThingType      string
	IotThingName      string
	DeviceFleetRole   string
	DeviceFleetBucket string
	Account           string
	Region            string
	AgentDirectory    string
	S3FolderPrefix    string
	TargetPlatform    TargetPlatform
}

func (cliArgs *CliArgs) Print() {
	fmt.Printf("Account: %s\n", cliArgs.Account)
	fmt.Printf("Region: %s\n", cliArgs.Region)
	fmt.Printf("DeviceFleet: %s\n", cliArgs.DeviceFleet)
	fmt.Printf("DeviceName: %s\n", cliArgs.DeviceName)
	fmt.Printf("IOT Thing Type: %s\n", cliArgs.IotThingType)
	fmt.Printf("IOT Thing Name: %s\n", cliArgs.IotThingName)
	fmt.Printf("Device Fleet Role: %s\n", cliArgs.DeviceFleetRole)
	fmt.Printf("Device Fleet Bucket: %s\n", cliArgs.DeviceFleetBucket)
	fmt.Printf("Agent Directory: %s\n", cliArgs.AgentDirectory)
	cliArgs.TargetPlatform.Print()
}

func ParseArgs(cliArgs *CliArgs) {
	accountId := flag.String("account", "", "AWS AccountId")
	region := flag.String("region", "us-west-2", "AWS Region")
	deviceFleet := flag.String("deviceFleet", "", "Name of the device fleet.")
	deviceName := flag.String("deviceName", "", "Name of the device.")

	targetOs := flag.String("os", "", "Name of Os")
	targetArch := flag.String("arch", "", "Name of device architecture.")
	targetAccelerator := flag.String("accelerator", "", "Name of accelerator.")

	iotThingType := flag.String("iotThingType", "", "Iot thing type for the device.")
	iotThingName := flag.String("iotThingName", "", "IOT thing name for the device.")
	deviceFleetRole := flag.String("deviceFleetRole", "", "Role for the device fleet.")
	deviceFleetBucket := flag.String("deviceFleetBucket", "", "Bucket to store device related data.")
	s3FolderPrefix := flag.String("s3FolderPrefix", "", "S3 prefix to store captured data.")

	cwd, err := os.Getwd()

	if err != nil {
		log.Fatal("Error", err)
	}

	defaultAgentDirectory := fmt.Sprintf("%s/.agent", cwd)
	agentDirectory := flag.String("agentDirectory", defaultAgentDirectory, "Local path to store agent")

	flag.Parse()

	if *deviceFleet == "" || *deviceName == "" || *accountId == "" {
		log.Fatal("Missing DeviceFleet or DeviceName or Account")
	}

	cliArgs.DeviceFleet = *deviceFleet
	cliArgs.DeviceName = *deviceName
	cliArgs.TargetPlatform = TargetPlatform{Os: strings.ToLower(*targetOs), Arch: strings.ToLower(*targetArch), Accelerator: strings.ToLower(*targetAccelerator)}
	cliArgs.TargetPlatform.Validate()

	cliArgs.Account = *accountId
	cliArgs.Region = *region
	if !strings.HasPrefix(*agentDirectory, "/") {
		*agentDirectory = fmt.Sprintf("%s/%s", cwd, *agentDirectory)
	}
	cliArgs.AgentDirectory = *agentDirectory

	if *iotThingType == "" {
		*iotThingType = fmt.Sprintf("Sagemaker_%s", cliArgs.DeviceFleet)
	}
	if *iotThingName == "" {
		*iotThingType = fmt.Sprintf("Sagemaker_%s", cliArgs.DeviceName)
	}
	if *deviceFleetRole == "" {
		*deviceFleetRole = fmt.Sprintf("Sagemaker_%s_role", cliArgs.DeviceFleet)
	}

	if *s3FolderPrefix == "" {
		*s3FolderPrefix = "demo"
	}

	cliArgs.IotThingType = *iotThingType
	cliArgs.IotThingName = *iotThingType
	cliArgs.DeviceFleetRole = *deviceFleetRole
	cliArgs.DeviceFleetBucket = *deviceFleetBucket
	cliArgs.S3FolderPrefix = *s3FolderPrefix
}
