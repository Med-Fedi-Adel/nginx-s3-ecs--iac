package infra

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func SetupIAM(ctx *pulumi.Context) (*IAMResources, error) {
	// Create an IAM role for ECS tasks
	taskRole, err := iam.NewRole(ctx, "nginxTaskRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": "ecs-tasks.amazonaws.com"
				}
			}]
		}`),
	})
	if err != nil {
		return nil, err
	}

	// Create an IAM role for ECS task execution
	taskExecRole, err := iam.NewRole(ctx, "nginxTaskExecRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [{
				"Action": "sts:AssumeRole",
				"Effect": "Allow",
				"Principal": {
					"Service": "ecs-tasks.amazonaws.com"
				}
			}]
		}`),
	})
	if err != nil {
		return nil, err
	}

	// Attach the AmazonECSTaskExecutionRolePolicy to the execution role
	_, err = iam.NewRolePolicyAttachment(ctx, "nginxTaskExecRolePolicy", &iam.RolePolicyAttachmentArgs{
		Role:      taskExecRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"),
	})
	if err != nil {
		return nil, err
	}

	return &IAMResources{
		TaskRole:     taskRole,
		TaskExecRole: taskExecRole,
	}, nil
}
