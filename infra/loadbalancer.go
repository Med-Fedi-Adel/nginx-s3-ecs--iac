package infra

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/lb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func SetupLoadBalancer(ctx *pulumi.Context, networking *NetworkingResources) (*LoadBalancerResources, error) {
	// Create an Application Load Balancer
	alb, err := lb.NewLoadBalancer(ctx, "nginxAlb", &lb.LoadBalancerArgs{
		Internal:         pulumi.Bool(false),
		SecurityGroups:   pulumi.StringArray{networking.WebSG.ID()},
		Subnets:          pulumi.ToStringArray(networking.Subnets.Ids),
		LoadBalancerType: pulumi.String("application"),
	})
	if err != nil {
		return nil, err
	}

	// Create a target group for the ALB
	targetGroup, err := lb.NewTargetGroup(ctx, "nginxTargetGroup", &lb.TargetGroupArgs{
		Port:       pulumi.Int(80),
		Protocol:   pulumi.String("HTTP"),
		TargetType: pulumi.String("ip"),
		VpcId:      pulumi.String(networking.VPC.Id),
		HealthCheck: &lb.TargetGroupHealthCheckArgs{
			Enabled:  pulumi.Bool(true),
			Path:     pulumi.String("/"),
			Protocol: pulumi.String("HTTP"),
			Port:     pulumi.String("80"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create a listener for the ALB
	listener, err := lb.NewListener(ctx, "nginxListener", &lb.ListenerArgs{
		LoadBalancerArn: alb.Arn,
		Port:            pulumi.Int(80),
		Protocol:        pulumi.String("HTTP"),
		DefaultActions: lb.ListenerDefaultActionArray{
			&lb.ListenerDefaultActionArgs{
				Type:           pulumi.String("forward"),
				TargetGroupArn: targetGroup.Arn,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return &LoadBalancerResources{
		ALB:         alb,
		TargetGroup: targetGroup,
		Listener:    listener,
	}, nil
}
