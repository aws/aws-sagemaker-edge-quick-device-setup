package aws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"smedge_installer/cli"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IamClient interface {
	CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
	ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
	AttachRolePolicy(ctx context.Context, params *iam.AttachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)
	GetPolicy(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error)
	CreatePolicy(ctx context.Context, params *iam.CreatePolicyInput, optFns ...func(*iam.Options)) (*iam.CreatePolicyOutput, error)
}

func CreateDeviceFleetRole(client IamClient, fleetName *string, roleName *string) *types.Role {
	assumeRolePolicyDocument := `{
		"Version": "2012-10-17",
		"Statement": [
			{
			  "Effect": "Allow",
			  "Principal": {"Service": "credentials.iot.amazonaws.com"},
			  "Action": ["sts:AssumeRole"]
			},
			{
			  "Effect": "Allow",
			  "Principal": {"Service": "sagemaker.amazonaws.com"},
			  "Action": ["sts:AssumeRole"]
			}
		]
	}`

	result, err := client.CreateRole(context.TODO(), &iam.CreateRoleInput{
		AssumeRolePolicyDocument: &assumeRolePolicyDocument,
		RoleName:                 roleName,
	})

	if err != nil {
		log.Fatal("Error", err)
	}

	return result.Role
}

func GetDeviceFleetRole(client IamClient, fleetName *string, roleName *string) *types.Role {
	result, err := client.GetRole(context.TODO(), &iam.GetRoleInput{
		RoleName: roleName,
	})

	if err != nil {
		var nse *types.NoSuchEntityException
		if errors.As(err, &nse) {
			log.Println("Role doesn't exist.")
			return nil
		}
		log.Fatal("Error", err)
	}

	return result.Role
}

func CheckIfPolicyIsAlreadyAttachedToTheRole(client IamClient, roleName *string, policyName *string) *types.AttachedPolicy {
	maxItems := int32(100)
	var marker *string

	for {
		ret, err := client.ListAttachedRolePolicies(context.TODO(), &iam.ListAttachedRolePoliciesInput{
			RoleName: roleName,
			MaxItems: &maxItems,
			Marker:   marker,
		})

		if err != nil {
			log.Fatal("Error", err)
		}

		for _, policy := range ret.AttachedPolicies {
			if *policy.PolicyName == *policyName {
				return &policy
			}
		}

		if ret.IsTruncated {
			marker = ret.Marker
		} else {
			break
		}
	}

	return nil
}

func AttachAmazonSageMakerEdgeDeviceFleetPolicy(client IamClient, role *types.Role, policyArn *string) {
	_, err := client.AttachRolePolicy(context.TODO(), &iam.AttachRolePolicyInput{
		PolicyArn: policyArn,
		RoleName:  role.RoleName,
	})

	if err != nil {
		log.Fatal("Error", err)
	}
}

type Principal struct {
	Service string
}

type StatementEntry struct {
	Sid       string `json:",omitempty"`
	Effect    string
	Action    []string
	Resource  []string
	Condition map[string]interface{} `json:",omitempty"`
	Principal Principal              `json:",omitempty"`
}

type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

func CreateDeviceFleetPolicy(client IamClient, cliArgs *cli.CliArgs) *types.Policy {
	var condition map[string]interface{}
	conditionByt := []byte(` {
		"StringEqualsIfExists": {
			"iam:PassedToService": [
				"iot.amazonaws.com",
				"credentials.iot.amazonaws.com"
			]
		}
	}`)

	if err := json.Unmarshal(conditionByt, &condition); err != nil {
		log.Fatal("Error", err)
	}

	policyDocument := &PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Sid:    "DeviceS3Access",
				Effect: "Allow",
				Action: []string{
					"s3:PutObject",
					"s3:GetBucketLocation",
				},
				Resource: []string{
					fmt.Sprintf("arn:aws:s3:::%s/*", cliArgs.DeviceFleetBucket),
					fmt.Sprintf("arn:aws:s3:::%s", cliArgs.DeviceFleetBucket),
				},
			},
			{
				Sid:    "SageMakerEdgeApis",
				Effect: "Allow",
				Action: []string{
					"sagemaker:SendHeartbeat",
					"sagemaker:GetDeviceRegistration",
				},
				Resource: []string{
					fmt.Sprintf("arn:aws:sagemaker:%s:%s:device-fleet/%s/device/%s", cliArgs.Region, cliArgs.Account, strings.ToLower(cliArgs.DeviceFleet), strings.ToLower(cliArgs.DeviceName)),
					fmt.Sprintf("arn:aws:sagemaker:%s:%s:device-fleet/%s", cliArgs.Region, cliArgs.Account, strings.ToLower(cliArgs.DeviceFleet)),
				},
			},
			{
				Sid:    "CreateIOTRoleAlias",
				Effect: "Allow",
				Action: []string{
					"iot:CreateRoleAlias",
					"iot:DescribeRoleAlias",
					"iot:UpdateRoleAlias",
					"iot:ListTagsForResource",
					"iot:TagResource",
				},
				Resource: []string{
					fmt.Sprintf("arn:aws:iot:%s:%s:rolealias/SageMakerEdge-%s", cliArgs.Region, cliArgs.Account, cliArgs.DeviceFleet),
				},
			},
			{
				Sid:    "CreateIoTRoleAliasIamPermissionsGetRole",
				Effect: "Allow",
				Action: []string{
					"iam:GetRole",
				},
				Resource: []string{
					fmt.Sprintf("arn:aws:iam::%s:role/%s", cliArgs.Account, cliArgs.DeviceFleetRole),
				},
			},
			{
				Sid:    "CreateIoTRoleAliasIamPermissionsPassRole",
				Effect: "Allow",
				Action: []string{
					"iam:PassRole",
				},
				Resource: []string{
					fmt.Sprintf("arn:aws:iam::%s:role/%s", cliArgs.Account, cliArgs.DeviceFleetRole),
				},
				Condition: condition,
			},
		},
	}
	policy, _ := json.MarshalIndent(policyDocument, "", " ")
	policyDoc := string(policy)

	policyDescription := fmt.Sprintf("SageMaker device fleet policy for %s", cliArgs.DeviceFleet)
	policyPath := "/"
	policyName := fmt.Sprintf("%s-policy", strings.ToLower(cliArgs.DeviceFleet))
	policyArn := fmt.Sprintf("arn:aws:iam::%s:policy/%s", cliArgs.Account, policyName)

	getPolicyOutput, err := client.GetPolicy(context.TODO(), &iam.GetPolicyInput{
		PolicyArn: &policyArn,
	})

	if err != nil {
		var nse *types.NoSuchEntityException
		if errors.As(err, &nse) {
			ret, err := client.CreatePolicy(context.TODO(), &iam.CreatePolicyInput{
				Description:    &policyDescription,
				Path:           &policyPath,
				PolicyDocument: &policyDoc,
				PolicyName:     &policyName,
			})

			if err != nil {
				log.Fatal("Error", err)
			}

			return ret.Policy
		}

		log.Fatal("Error", err)
	}

	return getPolicyOutput.Policy
}

func CreateDeviceFleetRoleIfNotExists(client IamClient, fleetName *string, roleName *string, policy *types.Policy) *types.Role {
	role := GetDeviceFleetRole(client, fleetName, roleName)
	if role == nil {
		role = CreateDeviceFleetRole(client, fleetName, roleName)
	}
	attachedPolicy := CheckIfPolicyIsAlreadyAttachedToTheRole(client, role.RoleName, policy.PolicyName)
	if attachedPolicy == nil {
		log.Println("Attaching device fleet policy to the")
		AttachAmazonSageMakerEdgeDeviceFleetPolicy(client, role, policy.Arn)
	}
	return role
}
