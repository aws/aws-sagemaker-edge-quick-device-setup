package aws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
	"path/filepath"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iot/types"
)

type IotClient interface {
	DescribeThingType(ctx context.Context, params *iot.DescribeThingTypeInput, optFns ...func(*iot.Options)) (*iot.DescribeThingTypeOutput, error)
	CreateThingType(ctx context.Context, params *iot.CreateThingTypeInput, optFns ...func(*iot.Options)) (*iot.CreateThingTypeOutput, error)
	DescribeThing(ctx context.Context, params *iot.DescribeThingInput, optFns ...func(*iot.Options)) (*iot.DescribeThingOutput, error)
	CreateThing(ctx context.Context, params *iot.CreateThingInput, optFns ...func(*iot.Options)) (*iot.CreateThingOutput, error)
	CreateKeysAndCertificate(ctx context.Context, params *iot.CreateKeysAndCertificateInput, optFns ...func(*iot.Options)) (*iot.CreateKeysAndCertificateOutput, error)
	DescribeEndpoint(ctx context.Context, params *iot.DescribeEndpointInput, optFns ...func(*iot.Options)) (*iot.DescribeEndpointOutput, error)
	AttachThingPrincipal(ctx context.Context, params *iot.AttachThingPrincipalInput, optFns ...func(*iot.Options)) (*iot.AttachThingPrincipalOutput, error)
	CreatePolicy(ctx context.Context, params *iot.CreatePolicyInput, optFns ...func(*iot.Options)) (*iot.CreatePolicyOutput, error)
	AttachPolicy(ctx context.Context, params *iot.AttachPolicyInput, optFns ...func(*iot.Options)) (*iot.AttachPolicyOutput, error)
}

func GetIotThingType(client IotClient, iotThingType *string) *iot.DescribeThingTypeOutput {
	ret, err := client.DescribeThingType(context.TODO(), &iot.DescribeThingTypeInput{
		ThingTypeName: iotThingType,
	})

	if err != nil {
		var rnf *types.ResourceNotFoundException
		if errors.As(err, &rnf) {
			return nil
		}
		log.Fatalf("Failed to describe thing type %s, Encountered error %s\n", *iotThingType, err)
	}

	return ret
}

type CreateIotThingTypeOutput struct {
	ThingTypeArn  *string
	ThingTypeId   *string
	ThingTypeName *string
}

func CreateIotThingType(client IotClient, iotThingType *string) *CreateIotThingTypeOutput {

	describeThingTypeOutput := GetIotThingType(client, iotThingType)

	if describeThingTypeOutput != nil {
		return &CreateIotThingTypeOutput{
			ThingTypeName: describeThingTypeOutput.ThingTypeName,
			ThingTypeArn:  describeThingTypeOutput.ThingTypeArn,
			ThingTypeId:   describeThingTypeOutput.ThingTypeId,
		}
	}

	ret, err := client.CreateThingType(context.TODO(), &iot.CreateThingTypeInput{
		ThingTypeName: iotThingType,
	})

	if err != nil {
		log.Fatalf("Failed to create thing type %s.Encountered error %s\n", *iotThingType, err)
	}

	return &CreateIotThingTypeOutput{
		ThingTypeArn:  ret.ThingTypeArn,
		ThingTypeId:   ret.ThingTypeId,
		ThingTypeName: ret.ThingTypeName,
	}
}

func GetIotThing(client IotClient, iotThingName *string) *iot.DescribeThingOutput {
	ret, err := client.DescribeThing(context.TODO(), &iot.DescribeThingInput{
		ThingName: iotThingName,
	})

	if err != nil {
		var rne *types.ResourceNotFoundException

		if errors.As(err, &rne) {
			log.Println("Thing doesn't exist")
			return nil
		}
		log.Fatalf("Failed to describe thing %s. Encountered error %s\n", *iotThingName, err)
	}

	return ret
}

type CreateIotThingOutput struct {
	ThingName     *string
	ThingArn      *string
	ThingId       *string
	ThingTypeName *string
}

func CreateIotThing(client IotClient, iotThingType *string, iotThingName *string) *CreateIotThingOutput {

	describeThingOutput := GetIotThing(client, iotThingName)

	if describeThingOutput != nil {
		return &CreateIotThingOutput{
			ThingName: describeThingOutput.ThingName,
			ThingId:   describeThingOutput.ThingId,
			ThingArn:  describeThingOutput.ThingArn,
		}
	}

	ret, err := client.CreateThing(context.TODO(), &iot.CreateThingInput{
		ThingName:     iotThingName,
		ThingTypeName: iotThingType,
	})

	if err != nil {
		log.Fatalf("Failed to create thing %s. Encountered error %s\n", *iotThingName, err)
	}

	return &CreateIotThingOutput{
		ThingName: ret.ThingName,
		ThingId:   ret.ThingId,
		ThingArn:  ret.ThingArn,
	}
}

func CreateIOTCertificates(client IotClient) *iot.CreateKeysAndCertificateOutput {
	ret, err := client.CreateKeysAndCertificate(context.TODO(), &iot.CreateKeysAndCertificateInput{
		SetAsActive: true,
	})

	if err != nil {
		log.Fatalf("Failed to create keys and certificates. Encountered error %s\n", err)
	}
	return ret
}

func writeStringToFile(filePath *string, contents *string) {
	file, err := os.Create(*filePath)

	if err != nil {
		log.Fatalf("Failed to create file %s. Encountered error %s\n", *filePath, err)
	}

	defer file.Close()

	file.WriteString(*contents)

}

func WriteCertificatesToFile(certs *iot.CreateKeysAndCertificateOutput, fleetName *string, deviceName *string, certsDirectory *string) {
	os.MkdirAll(*certsDirectory, os.ModePerm)
	pemFilePath := filepath.Join(*certsDirectory, "device.pem.crt")
	privateKeyFilePath := filepath.Join(*certsDirectory, "private.pem.key")
	publicKeyFilePath := filepath.Join(*certsDirectory, "public.pem.key.pub")

	writeStringToFile(&pemFilePath, certs.CertificatePem)
	writeStringToFile(&privateKeyFilePath, certs.KeyPair.PrivateKey)
	writeStringToFile(&publicKeyFilePath, certs.KeyPair.PublicKey)
}

func GetIotCredentialProviderEndpoint(client IotClient, roleNameAlias *string) *string {
	endpointType := "iot:CredentialProvider"
	ret, err := client.DescribeEndpoint(context.TODO(), &iot.DescribeEndpointInput{
		EndpointType: &endpointType,
	})

	if err != nil {
		log.Fatalf("Failed to describe endpoint for %s. Encountered error %s\n", endpointType, err)
	}

	endpoint := fmt.Sprintf("https://%s/role-aliases/%s/credentials", *ret.EndpointAddress, *roleNameAlias)
	return &endpoint
}

func AttachThingToCertificate(client IotClient, certificateArn *string, iotThingName *string) {
	_, err := client.AttachThingPrincipal(context.TODO(), &iot.AttachThingPrincipalInput{
		Principal: certificateArn,
		ThingName: iotThingName,
	})

	if err != nil {
		log.Fatalf("Failed to attach thing principal for thing name %s and certificate %s. Encountered error %s\n", *iotThingName, *certificateArn, err)
	}
}

func CreateAndAttachRoleAliasPolicy(client IotClient, roleAliasArn *string, certArn *string, iotThingName *string) {
	policyDocument := `{		
		"Version": "2012-10-17",
		"Statement": {
		  "Effect": "Allow",
		  "Action": "iot:AssumeRoleWithCertificate",
		  "Resource": "%s"
		}
	}`

	policyDocument = fmt.Sprintf(policyDocument, *roleAliasArn)
	now := time.Now()
	policyName := fmt.Sprintf("aliaspolicy-%d", now.UTC().Unix())

	if _, err := client.CreatePolicy(context.TODO(), &iot.CreatePolicyInput{
		PolicyName:     &policyName,
		PolicyDocument: &policyDocument,
	}); err != nil {
		log.Fatalf("Failed to cerate iot policy %s. Encountered error %s\n", policyName, err)
	}

	if _, err := client.AttachPolicy(context.TODO(), &iot.AttachPolicyInput{
		PolicyName: &policyName,
		Target:     certArn,
	}); err != nil {
		log.Fatalf("Failed to attach iot policy %s. Encountered error %s\n", policyName, err)
	}
}
