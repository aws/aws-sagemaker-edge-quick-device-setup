| &#9888; **WARNING**: This tool is meant for development/testing use only. We don't recommend the use of this tool for production! &#9888; |
| --- |

# aws-sagemaker-egde-quick-device-setup

This package provides a command line interface to easily onboard device with [SageMaker Edge](https://aws.amazon.com/sagemaker/edge/). Run the cli on the device you would like to provision as it will create all the necessary artifacts on the device.

Jump to:

- [Getting Started ](#getting-started)
- [Getting Help](#getting-help)
- [More Resources](#more-resource)


Getting Started
---------------

This README is for aws-sagemaker-edge-quick-device-setup version 0.0.1

Installation
------------

`aws-sagemaker-edge-quick-device-setup` is written in golang. You can also geneate the binary directly from the source using

`go build`

We support out of the the box distributions for know os and architectures. Check out `Releases <#releases>`__ for latest distributions.

Configuration
-------------

Before using the cli, you need to configure your AWS credentials. Go to https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html to learn about how to configure aws credentials.


Permissions
-----------

In order to invoke the CLI to create required resources in cloud the user/role must have required permission. You can create/attach a policy containng the following permissions.

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:GetRole",
                "iam:AttachRolePolicy",
                "iam:CreatePolicy",
                "iam:PassRole",
                "iam:GetPolicy",
                "iam:CreateRole",
                "iam:ListAttachedRolePolicies",
                "iot:GetPolicy",
                "iot:CreateThing",
                "iot:AttachPolicy",
                "iot:AttachThingPrincipal",
                "iot:DescribeThing",
                "iot:CreatePolicy",
                "iot:CreateThingType",
                "iot:CreateKeysAndCertificate",
                "iot:DescribeThingType",
                "s3:CreateBucket",
                "sagemaker:DescribeDeviceFleet",
                "sagemaker:RegisterDevices",
                "sagemaker:UpdateDevices",
                "sagemaker:CreateDeviceFleet",
                "sagemaker:DescribeDevice"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": "s3:GetObject",
            "Resource": "arn:aws:s3:::sagemaker-edge-release-store-*/*"
        }
    ]
}
```

Basic Commands
--------------

The CLI command has the following structure:

```
   $ quick-setup --[options]
```

Following are all the opetions supported by the cli
```
  -accelerator string
        Name of accelerator.
  -account string
        AWS AccountId
  -agentDirectory string
        Local path to store agent (default "/home/ubuntu/edge_manager/aws-sagemaker-edge-quick-device-setup/src/demo-agent")
  -arch string
        Name of device architecture.
  -deviceFleet string
        Name of the device fleet.
  -deviceFleetBucket string
        Bucket to store device related data.
  -deviceFleetRole string
        Role for the device fleet.
  -deviceName string
        Name of the device.
  -iotThingName string
        IOT thing name for the device.
  -iotThingType string
        Iot thing type for the device.
  -os string
        Name of Os
  -region string
        AWS Region (default "us-west-2")
  -s3FolderPrefix string
        S3 prefix to store captured data.
  -version
        Prints the version of aws-sagemaker-edge-quick-device-setup
```

To view help documentation, use one of the following:

```
   $ quick-setup --help
```

To get the version of the cli:

```
   $ quick-setup --version
```

Getting Help
------------

The best way to interact with our team is through GitHub. You can [open
an issue](https://github.com/aws/aws-sagemaker-edge-quick-device-setup/issues/new/choose) and
choose from one of our templates for guidance, bug reports, or feature
requests.


Please check for open similar
[issues](https://github.com/aws/aws-sagemaker-edge-quick-device-setup/issues/)before opening
another one.

More Resources
--------------

-  [Changelog](https://github.com/aws/aws-cli/blob/develop/CHANGELOG.rst)
   [Reference](https://docs.aws.amazon.com/cli/latest/reference/)
-  [Amazon Web Services Discussion
   Forums](https://forums.aws.amazon.com/)
-  [AWS Support](https://console.aws.amazon.com/support/home#/)

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This project is licensed under the Apache-2.0 License.
