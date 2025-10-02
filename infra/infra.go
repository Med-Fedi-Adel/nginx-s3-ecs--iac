package infra

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lb"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
)

// Shared types used across all infrastructure components
type Infrastructure struct {
	Networking   *NetworkingResources
	LoadBalancer *LoadBalancerResources
	IAM          *IAMResources
	ECS          *ECSResources
	S3           *S3Resources
}

type NetworkingResources struct {
	VPC     *ec2.LookupVpcResult
	Subnets *ec2.GetSubnetsResult
	WebSG   *ec2.SecurityGroup
}

type LoadBalancerResources struct {
	ALB         *lb.LoadBalancer
	TargetGroup *lb.TargetGroup
	Listener    *lb.Listener
}

type IAMResources struct {
	TaskRole     *iam.Role
	TaskExecRole *iam.Role
}

type ECSResources struct {
	Cluster        *ecs.Cluster
	TaskDefinition *ecs.TaskDefinition
	Service        *ecs.Service
}

type S3Resources struct {
	Bucket            *s3.Bucket
	PublicAccessBlock *s3.BucketPublicAccessBlock
	BucketPolicy      *s3.BucketPolicy
}
