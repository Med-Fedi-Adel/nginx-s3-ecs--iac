package infra

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func SetupNetworking(ctx *pulumi.Context) (*NetworkingResources, error) {
	// Get the default VPC
	defaultVpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{
		Default: pulumi.BoolRef(true),
	})
	if err != nil {
		return nil, err
	}

	// Get public subnets from the default VPC
	subnets, err := ec2.GetSubnets(ctx, &ec2.GetSubnetsArgs{
		Filters: []ec2.GetSubnetsFilter{
			{
				Name:   "vpc-id",
				Values: []string{defaultVpc.Id},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a security group for the load balancer and ECS tasks
	webSg, err := ec2.NewSecurityGroup(ctx, "webSecurityGroup", &ec2.SecurityGroupArgs{
		Description: pulumi.String("Allow HTTP traffic"),
		VpcId:       pulumi.String(defaultVpc.Id),
		Ingress: ec2.SecurityGroupIngressArray{
			&ec2.SecurityGroupIngressArgs{
				Protocol:   pulumi.String("tcp"),
				FromPort:   pulumi.Int(80),
				ToPort:     pulumi.Int(80),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
		Egress: ec2.SecurityGroupEgressArray{
			&ec2.SecurityGroupEgressArgs{
				Protocol:   pulumi.String("-1"),
				FromPort:   pulumi.Int(0),
				ToPort:     pulumi.Int(0),
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &NetworkingResources{
		VPC:     defaultVpc,
		Subnets: subnets,
		WebSG:   webSg,
	}, nil
}
