package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"smedge_installer/cli"

	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
)

type SagemakerClient interface {
	DescribeDeviceFleet(ctx context.Context, params *sagemaker.DescribeDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceFleetOutput, error)
	CreateDeviceFleet(ctx context.Context, params *sagemaker.CreateDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.CreateDeviceFleetOutput, error)
	DescribeDevice(ctx context.Context, params *sagemaker.DescribeDeviceInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceOutput, error)
	RegisterDevices(ctx context.Context, params *sagemaker.RegisterDevicesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.RegisterDevicesOutput, error)
}

func GetDeviceFleet(client SagemakerClient, fleetName *string) *sagemaker.DescribeDeviceFleetOutput {

	ret, err := client.DescribeDeviceFleet(context.TODO(), &sagemaker.DescribeDeviceFleetInput{
		DeviceFleetName: fleetName,
	})

	if err != nil {
		return nil
	}

	return ret

}

func CreateDeviceFleet(client SagemakerClient, fleetName *string, role *iamTypes.Role, s3Bucket *string) {
	s3OutputLocation := fmt.Sprintf("s3://%s/%s", *s3Bucket, *fleetName)

	describeDeviceFleetOutput := GetDeviceFleet(client, fleetName)

	if describeDeviceFleetOutput == nil {
		_, err := client.CreateDeviceFleet(context.TODO(), &sagemaker.CreateDeviceFleetInput{
			DeviceFleetName: fleetName,
			OutputConfig: &types.EdgeOutputConfig{
				S3OutputLocation: &s3OutputLocation,
			},
			RoleArn: role.Arn,
		})

		if err != nil {
			var oe *smithy.OperationError
			if errors.As(err, &oe) {
				log.Printf("failed to call service: %s, operation: %s, error: %v", oe.Service(), oe.Operation(), oe.Unwrap())
			}
			log.Fatal("Error: ", reflect.TypeOf(err))
		}

	}
}

func GetDevice(client SagemakerClient, fleetName *string, deviceName *string) *sagemaker.DescribeDeviceOutput {
	ret, err := client.DescribeDevice(context.TODO(), &sagemaker.DescribeDeviceInput{
		DeviceFleetName: fleetName,
		DeviceName:      deviceName,
	})

	if err != nil {
		return nil
	}

	return ret
}

func RegisterDevice(client SagemakerClient, fleetName *string, deviceName *string, iotThingName *string, targetPlatform *cli.TargetPlatform) {

	getDeviceOutput := GetDevice(client, fleetName, deviceName)

	targetOsKey := "os"
	targetArchKey := "arch"
	targetAccelerator := "accelerator"

	if getDeviceOutput == nil {
		_, err := client.RegisterDevices(context.TODO(), &sagemaker.RegisterDevicesInput{
			DeviceFleetName: fleetName,
			Devices: []types.Device{
				{
					DeviceName:   deviceName,
					IotThingName: iotThingName,
				},
			},
			Tags: []types.Tag{
				{
					Key:   &targetOsKey,
					Value: &targetPlatform.Os,
				},
				{
					Key:   &targetArchKey,
					Value: &targetPlatform.Arch,
				},
				{
					Key:   &targetAccelerator,
					Value: &targetPlatform.Accelerator,
				},
			},
		})

		if err != nil {
			log.Fatal("Error: ", err)
		}
	}
}

func GetRoleAliasArn(client SagemakerClient, deviceFleet *string) *string {
	ret, err := client.DescribeDeviceFleet(context.TODO(), &sagemaker.DescribeDeviceFleetInput{
		DeviceFleetName: deviceFleet,
	})

	if err != nil {
		log.Fatal("Error", err)
	}

	return ret.IotRoleAlias
}
