package infra

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func SetupS3(ctx *pulumi.Context) (*S3Resources, error) {
	// Creating an S3 bucket to host static files
	bucket, err := s3.NewBucket(ctx, "myStaticWebsiteBucket", &s3.BucketArgs{
		Website: &s3.BucketWebsiteArgs{
			IndexDocument: pulumi.String("index.html"),
			ErrorDocument: pulumi.String("error.html"),
		},
	})
	if err != nil {
		return nil, err
	}

	// Disable the block on public policies for the bucket
	publicAccessBlock, err := s3.NewBucketPublicAccessBlock(ctx, "myBucketPublicAccessBlock", &s3.BucketPublicAccessBlockArgs{
		Bucket:                bucket.ID(),
		BlockPublicAcls:       pulumi.Bool(false),
		BlockPublicPolicy:     pulumi.Bool(false),
		IgnorePublicAcls:      pulumi.Bool(false),
		RestrictPublicBuckets: pulumi.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	// Creating and applying a bucket policy to make its content public
	policy := pulumi.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": "*",
				"Action": "s3:GetObject",
				"Resource": "arn:aws:s3:::%s/*"
			}
		]
	}`, bucket.ID())

	bucketPolicy, err := s3.NewBucketPolicy(ctx, "bucketPolicy", &s3.BucketPolicyArgs{
		Bucket: bucket.ID(),
		Policy: policy,
	})
	if err != nil {
		return nil, err
	}

	return &S3Resources{
		Bucket:            bucket,
		PublicAccessBlock: publicAccessBlock,
		BucketPolicy:      bucketPolicy,
	}, nil
}
