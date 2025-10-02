package infra

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func SetupECS(ctx *pulumi.Context, networking *NetworkingResources, lb *LoadBalancerResources, iam *IAMResources) (*ECSResources, error) {
	// Create an ECS cluster
	cluster, err := ecs.NewCluster(ctx, "nginxCluster", nil)
	if err != nil {
		return nil, err
	}

	// Create the ECS task definition
	taskDefinition, err := ecs.NewTaskDefinition(ctx, "nginxTask", &ecs.TaskDefinitionArgs{
		Family:                  pulumi.String("nginx"),
		Cpu:                     pulumi.String("256"),
		Memory:                  pulumi.String("512"),
		NetworkMode:             pulumi.String("awsvpc"),
		RequiresCompatibilities: pulumi.StringArray{pulumi.String("FARGATE")},
		ExecutionRoleArn:        iam.TaskExecRole.Arn,
		TaskRoleArn:             iam.TaskRole.Arn,
		ContainerDefinitions: pulumi.String(`[{
			"name": "nginx",
			"image": "nginx:latest",
			"portMappings": [{
				"containerPort": 80,
				"hostPort": 80,
				"protocol": "tcp"
			}],
			"essential": true
		}]`),
	})
	if err != nil {
		return nil, err
	}

	// Create the ECS service
	service, err := ecs.NewService(ctx, "nginxService", &ecs.ServiceArgs{
		Cluster:        cluster.Arn,
		DesiredCount:   pulumi.Int(1),
		LaunchType:     pulumi.String("FARGATE"),
		TaskDefinition: taskDefinition.Arn,
		NetworkConfiguration: &ecs.ServiceNetworkConfigurationArgs{
			AssignPublicIp: pulumi.Bool(true),
			Subnets:        pulumi.ToStringArray(networking.Subnets.Ids),
			SecurityGroups: pulumi.StringArray{networking.WebSG.ID()},
		},
		LoadBalancers: ecs.ServiceLoadBalancerArray{
			&ecs.ServiceLoadBalancerArgs{
				TargetGroupArn: lb.TargetGroup.Arn,
				ContainerName:  pulumi.String("nginx"),
				ContainerPort:  pulumi.Int(80),
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{lb.ALB}))
	if err != nil {
		return nil, err
	}

	return &ECSResources{
		Cluster:        cluster,
		TaskDefinition: taskDefinition,
		Service:        service,
	}, nil
}
