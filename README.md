| **WARNING**: Use `aws-sagemaker-egde-quick-device-setup` in the context of development/testing. We don't recommend the use of this tool in production! |
| --- |

# aws-sagemaker-egde-quick-device-setup

This package provides a command line interface to easily onboard device with [SageMaker Edge](https://aws.amazon.com/sagemaker/edge/)

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

Before using the cli, you need to configure your AWS credentials.
You can do this in several ways:

-  Configuration command
-  Environment variables
-  Shared credentials file
-  Config file
-  IAM Role


The quickest way to get started is to run the ``aws configure`` command:

```
   $ aws configure
   AWS Access Key ID: MYACCESSKEY
   AWS Secret Access Key: MYSECRETKEY
   Default region name [us-west-2]: us-west-2
   Default output format [None]: json
```

To use environment variables, do the following:

```
   $ export AWS_ACCESS_KEY_ID=<access_key>
   $ export AWS_SECRET_ACCESS_KEY=<secret_key>
```

To use the shared credentials file, create an INI formatted file like
this:

```
   [default]
   aws_access_key_id=MYACCESSKEY
   aws_secret_access_key=MYSECRETKEY

   [testing]
   aws_access_key_id=MYACCESKEY
   aws_secret_access_key=MYSECRETKEY
```

and place it in ``~/.aws/credentials`` (or in
``%UserProfile%\.aws/credentials`` on Windows). If you wish to place the
shared credentials file in a different location than the one specified
above, you need to tell aws-cli where to find it. Do this by setting the
appropriate environment variable:

```
   $ export AWS_SHARED_CREDENTIALS_FILE=/path/to/shared_credentials_file
```

To use a config file, create an INI formatted file like this:

```
   [default]
   aws_access_key_id=<default access key>
   aws_secret_access_key=<default secret key>
   # Optional, to define default region for this profile.
   region=us-west-1

   [profile testing]
   aws_access_key_id=<testing access key>
   aws_secret_access_key=<testing secret key>
   region=us-west-2
```

and place it in ``~/.aws/config`` (or in ``%UserProfile%\.aws\config``
on Windows). If you wish to place the config file in a different
location than the one specified above, you need to tell the cli
where to find it. Do this by setting the appropriate environment
variable:

```
   $ export AWS_CONFIG_FILE=/path/to/config_file
```

As you can see, you can have multiple ``profiles`` defined in both the
shared credentials file and the configuration file. You can then specify
which profile to use by using the ``--profile`` option. If no profile is
specified the ``default`` profile is used.

In the config file, except for the default profile, you **must** prefix
each config section of a profile group with ``profile``. For example, if
you have a profile named "testing" the section header would be
``[profile testing]``.

The final option for credentials is highly recommended if you are using
the AWS CLI on an EC2 instance. [IAM
Roles](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-roles-for-amazon-ec2.html)
are a great way to have credentials installed automatically on your
instance. If you are using IAM Roles, the AWS CLI will find and use them
automatically.

In addition to credentials, a number of other variables can be
configured either with environment variables, configuration file
entries, or both. See the [AWS Tools and SDKs Shared Configuration and
Credentials Reference
Guide](https://docs.aws.amazon.com/credref/latest/refdocs/overview.html)
for more information.

For more information about configuration options, please refer to the
AWS CLI Configuration Variables
[topic](http://docs.aws.amazon.com/cli/latest/topic/config-vars.html#cli-aws-help-config-vars).
You can access this topic from the AWS CLI as well by running
``aws help config-vars``.


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

