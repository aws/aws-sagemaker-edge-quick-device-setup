package common

import (
	"aws-sagemaker-edge-quick-device-setup/cli"
	"testing"
)

func TestFromCliArgs(t *testing.T) {
	var config AgentConfig
	tp := cli.TargetPlatform{"linux", "amd64", ""}
	cliArgs := cli.CliArgs{
		DeviceFleet:       "some-fleet",
		DeviceName:        "some-device",
		IotThingType:      "",
		IotThingName:      "",
		DeviceFleetRole:   "sagemaker-edge-device-fleet-role-321",
		DeviceFleetBucket: "sagemaker-edge-bucket-use",
		Account:           "012345679012",
		Region:            "us-west-2",
		AgentDirectory:    "/home/ubuntu/smedge_agent",
		S3FolderPrefix:    "random-prefix",
		TargetPlatform:    tp,
		EnableDB:          false,
		EnableDeployment:  false,
	}
	config.FromCliArgs(&cliArgs)

	if config.ProviderProviderPath != "/home/ubuntu/smedge_agent/lib/libprovider_aws.so" {
		t.Fatal("Mismatch in provider path")
	}
	if config.DBModulePath != "" {
		t.Fatal("When db is disabled module path should not be set")
	}
	if config.LocalDataRootPath != "" {
		t.Fatal("When db is disabled local data root path should not be used.")
	}
	if config.DeploymentLibPath != "" {
		t.Fatal("When deployment is disabled module path should not be set")
	}
	if config.DeploymentPollInterval != 0 {
		t.Fatal("When deployment  is disabled interval should be zero")
	}
	cliArgsEnabled := cli.CliArgs{
		DeviceFleet:       "some-fleet",
		DeviceName:        "some-device",
		IotThingType:      "",
		IotThingName:      "",
		DeviceFleetRole:   "sagemaker-edge-device-fleet-role-321",
		DeviceFleetBucket: "sagemaker-edge-bucket-use",
		Account:           "012345679012",
		Region:            "us-west-2",
		AgentDirectory:    "/home/ubuntu/smedge_agent",
		S3FolderPrefix:    "random-prefix",
		TargetPlatform:    tp,
		EnableDB:          true,
		EnableDeployment:  true,
	}
	config.FromCliArgs(&cliArgsEnabled)
	if config.DBModulePath != "/home/ubuntu/smedge_agent/lib/libsagemaker_edge_db_handler_library.so" {
		t.Fatal("When db is enabled module path should  be set")
	}
	if config.LocalDataRootPath != "/home/ubuntu/smedge_agent/local_data" {
		t.Fatal("When db is enabled local data root path should  be used.")
	}
	if config.DeploymentLibPath != "/home/ubuntu/smedge_agent/lib/libdeployment_smedge_library.so" {
		t.Fatal("When deployment is enabled module path should  be set")
	}
	if config.DeploymentPollInterval != 1440 {
		t.Fatal("When deployment  is enabled interval should be 1440.")
	}
}
