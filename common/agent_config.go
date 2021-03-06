package common

import (
	"aws-sagemaker-edge-quick-device-setup/cli"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type AgentConfig struct {
	DeviceName                   string `json:"sagemaker_edge_core_device_name"`
	DeviceFleetName              string `json:"sagemaker_edge_core_device_fleet_name"`
	IotThingName                 string `json:"sagemaker_edge_core_iot_thing_name"`
	CapturDataBatchSize          int    `json:"sagemaker_edge_core_capture_data_batch_size"`
	CaptureDataBufferSize        int    `json:"sagemaker_edge_core_capture_data_buffer_size"`
	CaptureDataPushPeriodSeconds int    `json:"sagemaker_edge_core_capture_data_push_period_seconds"`
	FolderPrefix                 string `json:"sagemaker_edge_core_folder_prefix"`
	Region                       string `json:"sagemaker_edge_core_region"`
	AwsRootCertsPath             string `json:"sagemaker_edge_core_root_certs_path"`
	AwsCaCertFile                string `json:"sagemaker_edge_provider_aws_ca_cert_file"`
	AwsCertFile                  string `json:"sagemaker_edge_provider_aws_cert_file"`
	AwsCertPKFile                string `json:"sagemaker_edge_provider_aws_cert_pk_file"`
	ProviderAwsIotCredEndpoint   string `json:"sagemaker_edge_provider_aws_iot_cred_endpoint"`
	ProviderProvider             string `json:"sagemaker_edge_provider_provider"`
	ProviderProviderPath         string `json:"sagemaker_edge_provider_provider_path"`
	S3BucketName                 string `json:"sagemaker_edge_provider_s3_bucket_name"`
	DataCaptureDestination       string `json:"sagemaker_edge_core_capture_data_destination"`
}

func (config *AgentConfig) FromCliArgs(cliArgs *cli.CliArgs) {
	config.DeviceName = cliArgs.DeviceName
	config.DeviceFleetName = cliArgs.DeviceFleet
	config.IotThingName = cliArgs.IotThingName
	config.CapturDataBatchSize = 1
	config.CaptureDataBufferSize = 2
	config.CaptureDataPushPeriodSeconds = 5
	config.FolderPrefix = cliArgs.S3FolderPrefix
	config.Region = cliArgs.Region
	config.AwsRootCertsPath = filepath.Join(cliArgs.AgentDirectory, "certificates")
	config.AwsCaCertFile = filepath.Join(cliArgs.AgentDirectory, "iot-credentials", "AmazonRootCA1.pem")
	config.AwsCertFile = filepath.Join(cliArgs.AgentDirectory, "iot-credentials", "device.pem.crt")
	config.AwsCertPKFile = filepath.Join(cliArgs.AgentDirectory, "iot-credentials", "private.pem.key")
	config.ProviderAwsIotCredEndpoint = "endpoint"
	config.ProviderProvider = "Aws"
	config.ProviderProviderPath = filepath.Join(cliArgs.AgentDirectory, "lib", "libprovider_aws.so")
	config.S3BucketName = cliArgs.DeviceFleetBucket
	config.DataCaptureDestination = "Cloud"
}

func (config *AgentConfig) WriteToJson(filepath *string) {
	conf, _ := json.MarshalIndent(config, "", " ")
	fmt.Println(string(conf))
	_ = ioutil.WriteFile(*filepath, conf, 0400)
}
