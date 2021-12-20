package aws

import (
	"context"
	"testing"

	iamTypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

type mockSagemakerClient struct{}

var mockDescribeDeviceFleet func(ctx context.Context, params *sagemaker.DescribeDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceFleetOutput, error)
var mockCreateDeviceFleet func(ctx context.Context, params *sagemaker.CreateDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.CreateDeviceFleetOutput, error)
var mockDescribeDevice func(ctx context.Context, params *sagemaker.DescribeDeviceInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceOutput, error)
var mockRegisterDevices func(ctx context.Context, params *sagemaker.RegisterDevicesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.RegisterDevicesOutput, error)

func (sm mockSagemakerClient) DescribeDeviceFleet(ctx context.Context, params *sagemaker.DescribeDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceFleetOutput, error) {
	return mockDescribeDeviceFleet(ctx, params, optFns...)
}

func (sm mockSagemakerClient) CreateDeviceFleet(ctx context.Context, params *sagemaker.CreateDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.CreateDeviceFleetOutput, error) {
	return mockCreateDeviceFleet(ctx, params, optFns...)
}

func (sm mockSagemakerClient) DescribeDevice(ctx context.Context, params *sagemaker.DescribeDeviceInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceOutput, error) {
	return mockDescribeDevice(ctx, params, optFns...)
}

func (sm mockSagemakerClient) RegisterDevices(ctx context.Context, params *sagemaker.RegisterDevicesInput, optFns ...func(*sagemaker.Options)) (*sagemaker.RegisterDevicesOutput, error) {
	return mockRegisterDevices(ctx, params, optFns...)
}

func TestGetDeviceFleet(t *testing.T) {
	client := mockSagemakerClient{}
	nonExistantDeviceFleet := "NonExistantDeviceFleet"
	dummyFleet := "DummyFleet"

	mockDescribeDeviceFleet = func(ctx context.Context, params *sagemaker.DescribeDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceFleetOutput, error) {
		if *params.DeviceFleetName == nonExistantDeviceFleet {
			return nil, &types.ResourceNotFound{}
		}
		return &sagemaker.DescribeDeviceFleetOutput{
			DeviceFleetArn:  params.DeviceFleetName,
			DeviceFleetName: params.DeviceFleetName,
		}, nil
	}

	ret := GetDeviceFleet(client, &nonExistantDeviceFleet)

	if ret != nil {
		t.Fatalf("Should return nil for non existant device fleet")
	}

	ret = GetDeviceFleet(client, &dummyFleet)

	if *ret.DeviceFleetName != dummyFleet {
		t.Fatalf("Invalid device fleet name.")
	}

}

func TestCreateDeviceFleet(t *testing.T) {
	client := mockSagemakerClient{}
	nonExistantDeviceFleet := "NonExistantDeviceFleet"
	existingDeviceFleet := "ExistingDeviceFleet"
	roleArn := "DummyRole"
	role := iamTypes.Role{
		Arn: &roleArn,
	}
	s3Bucket := "DummyBucket"

	mockDescribeDeviceFleet = func(ctx context.Context, params *sagemaker.DescribeDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceFleetOutput, error) {
		if *params.DeviceFleetName == existingDeviceFleet {
			return &sagemaker.DescribeDeviceFleetOutput{
				DeviceFleetArn:  params.DeviceFleetName,
				DeviceFleetName: params.DeviceFleetName,
			}, nil
		}
		return nil, &types.ResourceNotFound{}
	}

	mockCreateDeviceFleet = func(ctx context.Context, params *sagemaker.CreateDeviceFleetInput, optFns ...func(*sagemaker.Options)) (*sagemaker.CreateDeviceFleetOutput, error) {
		if *params.DeviceFleetName == existingDeviceFleet {
			t.Fatalf("Shouldn't be called for existing device fleet")
		}
		return &sagemaker.CreateDeviceFleetOutput{}, nil
	}

	CreateDeviceFleet(client, &existingDeviceFleet, &role, &s3Bucket)
	CreateDeviceFleet(client, &nonExistantDeviceFleet, &role, &s3Bucket)
}

func TestGetDevice(t *testing.T) {
	client := mockSagemakerClient{}
	nonExistantDevice := "NonExistantDevice"
	existingDevice := "ExistantDevice"
	existingFleet := "ExisingFleet"
	mockDescribeDevice = func(ctx context.Context, params *sagemaker.DescribeDeviceInput, optFns ...func(*sagemaker.Options)) (*sagemaker.DescribeDeviceOutput, error) {
		if *params.DeviceName == nonExistantDevice {
			return nil, &types.ResourceNotFound{}
		}

		return &sagemaker.DescribeDeviceOutput{
			DeviceName:      params.DeviceName,
			DeviceFleetName: params.DeviceFleetName,
		}, nil
	}

	ret := GetDevice(client, &existingFleet, &nonExistantDevice)
	if ret != nil {
		t.Fatalf("Should return nil for non existing device")
	}

	ret = GetDevice(client, &existingFleet, &existingDevice)

	if *ret.DeviceName != existingDevice {
		t.Fatalf("Invalid device")
	}
}
