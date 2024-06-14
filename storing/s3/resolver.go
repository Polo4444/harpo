package storing_s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	transport "github.com/aws/smithy-go/endpoints"
)

type resolverV2 struct {
	bucket   string
	region   string
	endpoint string

	// When true, force a path-style endpoint to be used where the bucket name is part
	// of the path.
	//
	// Defaults to false if no value is
	// provided.
	//
	// AWS::S3::ForcePathStyle
	forcePath bool
}

func NewResolverV2(bucket, region, endpoint string, forcePath bool) *resolverV2 {
	return &resolverV2{bucket: bucket, region: region, endpoint: endpoint}
}

func (r *resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (transport.Endpoint, error) {

	params.Region = aws.String(r.region)
	params.Endpoint = aws.String(r.endpoint)
	params.Bucket = aws.String(r.bucket)
	params.ForcePathStyle = aws.Bool(r.forcePath)

	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}
