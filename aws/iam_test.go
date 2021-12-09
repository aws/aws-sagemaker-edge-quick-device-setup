package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"quick-device-setup/cli"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type mockIam struct{}

var mockCreateRole func(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error)
var mockGetRole func(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
var mockListAttachedRolePolicies func(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error)
var mockAttachRolePolicy func(ctx context.Context, params *iam.AttachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error)
var mockGetPolicy func(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error)
var mockCreatePolicy func(ctx context.Context, params *iam.CreatePolicyInput, optFns ...func(*iam.Options)) (*iam.CreatePolicyOutput, error)

func (iam mockIam) CreateRole(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error) {
	return mockCreateRole(ctx, params, optFns...)
}

func (iam mockIam) GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
	return mockGetRole(ctx, params, optFns...)
}

func (iam mockIam) ListAttachedRolePolicies(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error) {
	return mockListAttachedRolePolicies(ctx, params, optFns...)
}

func (iam mockIam) AttachRolePolicy(ctx context.Context, params *iam.AttachRolePolicyInput, optFns ...func(*iam.Options)) (*iam.AttachRolePolicyOutput, error) {
	return mockAttachRolePolicy(ctx, params, optFns...)
}

func (iam mockIam) GetPolicy(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error) {
	return mockGetPolicy(ctx, params, optFns...)
}

func (iam mockIam) CreatePolicy(ctx context.Context, params *iam.CreatePolicyInput, optFns ...func(*iam.Options)) (*iam.CreatePolicyOutput, error) {
	return mockCreatePolicy(ctx, params, optFns...)
}

func TestCreateDeviceFleetRole(t *testing.T) {
	client := mockIam{}
	testFleetName := "DummyFleet"
	roleName := "DummyFleetRole"
	mockCreateRole = func(ctx context.Context, params *iam.CreateRoleInput, optFns ...func(*iam.Options)) (*iam.CreateRoleOutput, error) {
		arn := fmt.Sprintf("%s-%s", *params.RoleName, testFleetName)
		roleId := fmt.Sprintf("%s-id", *params.RoleName)
		role := types.Role{
			Arn:      &arn,
			RoleId:   &roleId,
			RoleName: params.RoleName,
		}
		createRoleOutput := iam.CreateRoleOutput{
			Role: &role,
		}

		policyDocument := PolicyDocument{}
		json.Unmarshal([]byte(*params.AssumeRolePolicyDocument), &policyDocument)

		statements := policyDocument.Statement

		for _, statement := range statements {
			if statement.Principal.Service != "credentials.iot.amazonaws.com" && statement.Principal.Service != "sagemaker.amazonaws.com" {
				t.Fatalf("Invalid service principal in trust policy")
			}

			if len(statement.Action) == 0 || statement.Action[0] != "sts:AssumeRole" {
				t.Fatalf("Invalid action in trust policy")
			}
		}

		fmt.Println(statements[0].Principal.Service)

		return &createRoleOutput, nil
	}
	deviceFleetRole := CreateDeviceFleetRole(client, &testFleetName, &roleName)

	if deviceFleetRole.RoleName != &roleName {
		t.Fatalf("Invalid Role Name")
	}
}

func TestGetDeviceFleetRole(t *testing.T) {
	client := mockIam{}
	dummyFleet := "DummyFleet"
	dummyRoleName := "DummyRole"
	nonExistentRoleName := "NonExistent"
	mockGetRole = func(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error) {
		if *params.RoleName == nonExistentRoleName {
			message := fmt.Sprintf("Role \"%s\" does not exist", *params.RoleName)
			return nil, &types.NoSuchEntityException{Message: &message}
		}
		roleArn := fmt.Sprintf("%s-role", *params.RoleName)
		role := types.Role{
			Arn:      &roleArn,
			RoleName: params.RoleName,
		}
		getRoleOutput := iam.GetRoleOutput{
			Role: &role,
		}
		return &getRoleOutput, nil
	}

	role := GetDeviceFleetRole(client, &dummyFleet, &dummyRoleName)

	if role.RoleName != &dummyRoleName {
		t.Fatalf("Invalid Role Name")
	}
}

func TestCheckIfPolicyIsAlreadyAttachedToTheRole(t *testing.T) {
	client := mockIam{}
	dummyRoleName := "DummyRoleName"
	unAttachedPolicy := "UnattachedPolicy"
	attachedPolicy := "AttachedPolicyName"
	mockListAttachedRolePolicies = func(ctx context.Context, params *iam.ListAttachedRolePoliciesInput, optFns ...func(*iam.Options)) (*iam.ListAttachedRolePoliciesOutput, error) {
		return &iam.ListAttachedRolePoliciesOutput{
			AttachedPolicies: []types.AttachedPolicy{
				{
					PolicyName: &attachedPolicy,
					PolicyArn:  &attachedPolicy,
				},
			},
		}, nil
	}

	policy := CheckIfPolicyIsAlreadyAttachedToTheRole(client, &dummyRoleName, &unAttachedPolicy)

	if policy != nil {
		t.Fatalf("Policy should return nil!")
	}

	policy = CheckIfPolicyIsAlreadyAttachedToTheRole(client, &dummyRoleName, &attachedPolicy)

	if policy == nil {
		t.Fatalf("Policy should not return nil!")
	}
}

func TestCreateDeviceFleetPolicy(t *testing.T) {
	client := mockIam{}
	policyName1 := "DummyPolicy1"
	policyName2 := "DummyPolicy2"
	cliArgs := cli.CliArgs{
		DeviceFleet:     "DummyFleet",
		DeviceName:      "DummyDevice",
		Account:         "DummyAccount",
		DeviceFleetRole: "DummyRole",
		Region:          "DummyRegion",
	}
	mockGetPolicy = func(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error) {
		return nil, &types.NoSuchEntityException{}
	}

	mockCreatePolicy = func(ctx context.Context, params *iam.CreatePolicyInput, optFns ...func(*iam.Options)) (*iam.CreatePolicyOutput, error) {

		dummyPolicy := types.Policy{
			PolicyName: &policyName1,
		}
		return &iam.CreatePolicyOutput{
			Policy: &dummyPolicy,
		}, nil
	}

	policy := CreateDeviceFleetPolicy(client, &cliArgs)

	if *policy.PolicyName != policyName1 {
		t.Fatalf("Invalid response")
	}

	mockGetPolicy = func(ctx context.Context, params *iam.GetPolicyInput, optFns ...func(*iam.Options)) (*iam.GetPolicyOutput, error) {
		dummyPolicy := types.Policy{
			PolicyName: &policyName2,
		}
		return &iam.GetPolicyOutput{
			Policy: &dummyPolicy,
		}, nil
	}

	policy = CreateDeviceFleetPolicy(client, &cliArgs)

	if *policy.PolicyName != policyName2 {
		t.Fatalf("Invalid response")
	}
}
