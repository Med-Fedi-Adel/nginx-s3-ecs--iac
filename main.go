package main

import (
	"fmt"
	"pulumi-tutorial/infra" // This imports your local infra package

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Initialize infrastructure components
		networking, err := infra.SetupNetworking(ctx)
		if err != nil {
			return err
		}

		loadBalancer, err := infra.SetupLoadBalancer(ctx, networking)
		if err != nil {
			return err
		}

		iam, err := infra.SetupIAM(ctx)
		if err != nil {
			return err
		}

		ecs, err := infra.SetupECS(ctx, networking, loadBalancer, iam)
		if err != nil {
			return err
		}

		s3, err := infra.SetupS3(ctx)
		if err != nil {
			return err
		}

		// Export all outputs
		ctx.Export("nginxUrl", loadBalancer.ALB.DnsName.ApplyT(func(dnsName string) string {
			return fmt.Sprintf("http://%s", dnsName)
		}).(pulumi.StringOutput))
		ctx.Export("bucketName", s3.Bucket.ID())
		ctx.Export("websiteUrl", s3.Bucket.WebsiteEndpoint)
		ctx.Export("ecsClusterName", ecs.Cluster.Name)

		return nil
	})
}
