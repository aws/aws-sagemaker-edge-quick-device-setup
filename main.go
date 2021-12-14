package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"quick-device-setup/aws"
	"quick-device-setup/cli"
	"quick-device-setup/common"
	"quick-device-setup/distinfo"
	"strings"
	"time"

	awsStd "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
)

func main() {
	cliArgs := cli.CliArgs{}
	cli.ParseArgs(&cliArgs)
	cliArgs.Print()

	// return retry.AddWithErrorCodes(retry.NewStandard(), (*smTypes.Mal)(nil).ErrorCode())
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(cliArgs.Region), config.WithRetryer(func() awsStd.Retryer {
		return retry.AddWithErrorCodes(retry.AddWithMaxBackoffDelay(retry.AddWithMaxAttempts(retry.NewStandard(), 5), 1*time.Second), "ValidationException", "ThrottlingException")
	}))

	if err != nil {
		log.Fatal("Failed to load default aws config. Encountered Error ", err)
	}

	iamClient := iam.NewFromConfig(cfg)
	smClient := sagemaker.NewFromConfig(cfg)
	iotClient := iot.NewFromConfig(cfg)
	s3Client := s3.NewFromConfig(cfg)

	log.Println("Distribution Info:")
	log.Printf("Version: %s\n", distinfo.VERSION)
	log.Printf("Os: %s\n", distinfo.OS)
	log.Printf("Arch: %s\n", distinfo.ARCH)

	log.Println("Step-1 Creating S3 bucket for storing device fleet data...")
	s3OutputLocation := aws.CreateS3Bucket(s3Client, &cliArgs.DeviceFleetBucket, &cliArgs.Account)
	if s3OutputLocation == nil {
		return
	}
	log.Println("Step-1 Completed.")

	cliArgs.DeviceFleetBucket = *s3OutputLocation

	log.Println("Step-2a Creating device fleet policy...")
	fleetPolicy := aws.CreateDeviceFleetPolicy(iamClient, &cliArgs)
	log.Println("Step-2b Creating device fleet bucket policy...")
	bucketPolicy := aws.CreateDeviceFleetBucketPolicy(iamClient, &cliArgs)
	log.Println("Step-2 Completed.")

	log.Println("Step-3 Creating device fleet role...")
	role := aws.CreateDeviceFleetRoleIfNotExists(iamClient, &cliArgs.DeviceFleet, &cliArgs.DeviceFleetRole, fleetPolicy, bucketPolicy)
	log.Println("Step-3 Completed.")

	log.Println("Step-4 Creating iot thing type...")
	aws.CreateIotThingType(iotClient, &cliArgs.IotThingType)
	log.Println("Step-4 Completed.")

	log.Println("Step-5 Creating iot thing...")
	aws.CreateIotThing(iotClient, &cliArgs.IotThingType, &cliArgs.IotThingName)
	log.Println("Step-5 Completed.")

	log.Println("Step-6 Creating device fleet...")
	aws.CreateDeviceFleet(smClient, &cliArgs.DeviceFleet, role, s3OutputLocation)
	log.Println("Step-6 Completed.")

	log.Println("Step-7 Registering device...")
	aws.RegisterDevice(smClient, &cliArgs.DeviceFleet, &cliArgs.DeviceName, &cliArgs.IotThingName, &cliArgs.TargetPlatform)
	log.Println("Step-7 Completed.")

	log.Println("Step-8 Downloading Agent...")
	common.DownloadAgent(s3Client, &cliArgs)
	log.Println("Step-8 Complted.")

	log.Println("Step-9 Downloading code signing root certificate...")
	common.DownloadSigningRootCert(s3Client, &cliArgs)
	log.Println("Step-9 Completed.")

	log.Println("Step-10 Creating iot certificates...")
	certs := aws.CreateIOTCertificates(iotClient)
	log.Println("Step-10 Completed.")

	log.Println("Step-10 Attaching certificate to thing...")
	aws.AttachThingToCertificate(iotClient, certs.CertificateArn, &cliArgs.IotThingName)
	log.Println("Step-10 Completed.")

	log.Println("Step-11 Configuring Agent...")
	certsDirectory := fmt.Sprintf("%s/iot-credentials", cliArgs.AgentDirectory)
	aws.WriteCertificatesToFile(certs, &cliArgs.DeviceFleet, &cliArgs.DeviceName, &certsDirectory)
	rootCAPath := fmt.Sprintf("%s/AmazonRootCA1.pem", certsDirectory)
	common.DownloadFile(rootCAPath, "https://www.amazontrust.com/repository/AmazonRootCA1.pem")
	config := common.AgentConfig{}
	configPath := fmt.Sprintf("%s/sagemaker_edge_config.json", cliArgs.AgentDirectory)
	config.FromCliArgs(&cliArgs)
	roleAliasArn := aws.GetRoleAliasArn(smClient, &cliArgs.DeviceFleet)
	aws.CreateAndAttachRoleAliasPolicy(iotClient, roleAliasArn, certs.CertificateArn, &cliArgs.IotThingName)
	roleAliasSplits := strings.Split(*roleAliasArn, "/")
	config.ProviderAwsIotCredEndpoint = *aws.GetIotCredentialProviderEndpoint(iotClient, &roleAliasSplits[1])
	config.WriteToJson(&configPath)
	agentBinaryPath := fmt.Sprintf("%s/bin/sagemaker_edge_agent_binary", cliArgs.AgentDirectory)
	agentClientPath := fmt.Sprintf("%s/bin/sagemaker_edge_agent_client_example", cliArgs.AgentDirectory)
	os.Chmod(agentBinaryPath, 0700)
	os.Chmod(agentClientPath, 0700)
	log.Println("Step-11 Completed.")
}
