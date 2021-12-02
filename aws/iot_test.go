package aws

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

type mockClient struct{}

var iotMockDescribeThingType func(ctx context.Context, params *iot.DescribeThingTypeInput, optFns ...func(*iot.Options)) (*iot.DescribeThingTypeOutput, error)
var iotMockCreateThingType func(ctx context.Context, params *iot.CreateThingTypeInput, optFns ...func(*iot.Options)) (*iot.CreateThingTypeOutput, error)
var iotMockDescribeThing func(ctx context.Context, params *iot.DescribeThingInput, optFns ...func(*iot.Options)) (*iot.DescribeThingOutput, error)
var iotMockCreateThing func(ctx context.Context, params *iot.CreateThingInput, optFns ...func(*iot.Options)) (*iot.CreateThingOutput, error)
var iotMockCreateKeysAndCertificate func(ctx context.Context, params *iot.CreateKeysAndCertificateInput, optFns ...func(*iot.Options)) (*iot.CreateKeysAndCertificateOutput, error)
var iotMockDescribeEndpoint func(ctx context.Context, params *iot.DescribeEndpointInput, optFns ...func(*iot.Options)) (*iot.DescribeEndpointOutput, error)
var iotMockAttachThingPrincipal func(ctx context.Context, params *iot.AttachThingPrincipalInput, optFns ...func(*iot.Options)) (*iot.AttachThingPrincipalOutput, error)
var iotMockCreatePolicy func(ctx context.Context, params *iot.CreatePolicyInput, optFns ...func(*iot.Options)) (*iot.CreatePolicyOutput, error)
var iotMockAttachPolicy func(ctx context.Context, params *iot.AttachPolicyInput, optFns ...func(*iot.Options)) (*iot.AttachPolicyOutput, error)

func (iot mockClient) DescribeThingType(ctx context.Context, params *iot.DescribeThingTypeInput, optFns ...func(*iot.Options)) (*iot.DescribeThingTypeOutput, error) {
	return iotMockDescribeThingType(ctx, params, optFns...)
}

func (iot mockClient) CreateThingType(ctx context.Context, params *iot.CreateThingTypeInput, optFns ...func(*iot.Options)) (*iot.CreateThingTypeOutput, error) {
	return iotMockCreateThingType(ctx, params, optFns...)
}

func (iot mockClient) DescribeThing(ctx context.Context, params *iot.DescribeThingInput, optFns ...func(*iot.Options)) (*iot.DescribeThingOutput, error) {
	return iotMockDescribeThing(ctx, params, optFns...)
}

func (iot mockClient) CreateThing(ctx context.Context, params *iot.CreateThingInput, optFns ...func(*iot.Options)) (*iot.CreateThingOutput, error) {
	return iotMockCreateThing(ctx, params, optFns...)
}

func (iot mockClient) CreateKeysAndCertificate(ctx context.Context, params *iot.CreateKeysAndCertificateInput, optFns ...func(*iot.Options)) (*iot.CreateKeysAndCertificateOutput, error) {
	return iotMockCreateKeysAndCertificate(ctx, params, optFns...)
}

func (iot mockClient) DescribeEndpoint(ctx context.Context, params *iot.DescribeEndpointInput, optFns ...func(*iot.Options)) (*iot.DescribeEndpointOutput, error) {
	return iotMockDescribeEndpoint(ctx, params, optFns...)
}

func (iot mockClient) AttachThingPrincipal(ctx context.Context, params *iot.AttachThingPrincipalInput, optFns ...func(*iot.Options)) (*iot.AttachThingPrincipalOutput, error) {
	return iotMockAttachThingPrincipal(ctx, params, optFns...)
}

func (iot mockClient) CreatePolicy(ctx context.Context, params *iot.CreatePolicyInput, optFns ...func(*iot.Options)) (*iot.CreatePolicyOutput, error) {
	return iotMockCreatePolicy(ctx, params, optFns...)
}

func (iot mockClient) AttachPolicy(ctx context.Context, params *iot.AttachPolicyInput, optFns ...func(*iot.Options)) (*iot.AttachPolicyOutput, error) {
	return iotMockAttachPolicy(ctx, params, optFns...)
}

func TestGetIotThingType(t *testing.T) {
	client := mockClient{}
	nonExistantThingType := "NonExistantThingType"
	dummyThingType := "DummyThingType"

	iotMockDescribeThingType = func(ctx context.Context, params *iot.DescribeThingTypeInput, optFns ...func(*iot.Options)) (*iot.DescribeThingTypeOutput, error) {
		if *params.ThingTypeName == nonExistantThingType {
			return nil, &types.ResourceNotFoundException{}
		} else {
			return &iot.DescribeThingTypeOutput{
				ThingTypeName: params.ThingTypeName,
			}, nil
		}
	}

	ret := GetIotThingType(client, &dummyThingType)

	if *ret.ThingTypeName != dummyThingType {
		t.Fatalf("Invalid thing type in response")
	}

	ret = GetIotThingType(client, &nonExistantThingType)

	if ret != nil {
		t.Fatalf(fmt.Sprintf("Should return nil for %s", nonExistantThingType))
	}
}

func TestCreateIotThingType(t *testing.T) {
	client := mockClient{}
	existingThingType := "existingThingType"
	dummyThingType := "DummyThingType"

	iotMockDescribeThingType = func(ctx context.Context, params *iot.DescribeThingTypeInput, optFns ...func(*iot.Options)) (*iot.DescribeThingTypeOutput, error) {
		if *params.ThingTypeName == existingThingType {
			return &iot.DescribeThingTypeOutput{
				ThingTypeName: params.ThingTypeName,
				ThingTypeArn:  params.ThingTypeName,
				ThingTypeId:   params.ThingTypeName,
			}, nil
		} else {
			return nil, &types.ResourceNotFoundException{}
		}
	}

	ret := CreateIotThingType(client, &existingThingType)

	if *ret.ThingTypeName != existingThingType {
		log.Fatalf("Should return existing thing type.")
	}

	iotMockCreateThingType = func(ctx context.Context, params *iot.CreateThingTypeInput, optFns ...func(*iot.Options)) (*iot.CreateThingTypeOutput, error) {
		return &iot.CreateThingTypeOutput{
			ThingTypeArn:  &dummyThingType,
			ThingTypeId:   &dummyThingType,
			ThingTypeName: &dummyThingType,
		}, nil
	}

	ret = CreateIotThingType(client, &dummyThingType)

	if *ret.ThingTypeName != dummyThingType {
		log.Fatalf(fmt.Sprintf("Should return %s!", dummyThingType))
	}
}

func TestGetIotThing(t *testing.T) {
	client := mockClient{}
	existingThingName := "ExisingThing"
	nonExistingThingName := "NonExisingThing"
	iotMockDescribeThing = func(ctx context.Context, params *iot.DescribeThingInput, optFns ...func(*iot.Options)) (*iot.DescribeThingOutput, error) {
		if *params.ThingName == existingThingName {
			return &iot.DescribeThingOutput{
				ThingName: &existingThingName,
			}, nil
		} else {
			return nil, &types.ResourceNotFoundException{}
		}
	}

	ret := GetIotThing(client, &existingThingName)

	if ret.ThingName != &existingThingName {
		t.Fatalf("Invalid thing name!")
	}

	ret = GetIotThing(client, &nonExistingThingName)

	if ret != nil {
		t.Fatalf("Should return nil for non existing thing")
	}
}

func TestCreateIotThing(t *testing.T) {
	client := mockClient{}
	existingThingName := "ExisingThing"
	nonExistingThingName := "NonExisingThing"
	thingType := "DummyThingType"
	iotMockDescribeThing = func(ctx context.Context, params *iot.DescribeThingInput, optFns ...func(*iot.Options)) (*iot.DescribeThingOutput, error) {
		if *params.ThingName == existingThingName {
			return &iot.DescribeThingOutput{
				ThingName: &existingThingName,
			}, nil
		} else {
			return nil, &types.ResourceNotFoundException{}
		}
	}

	iotMockCreateThing = func(ctx context.Context, params *iot.CreateThingInput, optFns ...func(*iot.Options)) (*iot.CreateThingOutput, error) {
		if *params.ThingName == existingThingName {
			t.Fatalf("This should not be called for existing thing name.")
		}
		return &iot.CreateThingOutput{
			ThingName: params.ThingName,
			ThingId:   params.ThingName,
			ThingArn:  params.ThingName,
		}, nil

	}

	ret := CreateIotThing(client, &thingType, &existingThingName)

	if *ret.ThingName != existingThingName {
		t.Fatalf("Invalid thing name")
	}

	ret = CreateIotThing(client, &thingType, &nonExistingThingName)

	if *ret.ThingName != nonExistingThingName {
		t.Fatalf("Invalid thing name")
	}
}
