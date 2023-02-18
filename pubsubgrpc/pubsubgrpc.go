// Package pubsubgrpc makes it easier to use raw gRPC with Google Cloud Pub/Sub.
package pubsubgrpc

import (
	"context"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

// https://cloud.google.com/pubsub/docs/reference/service_apis_overview
const endpoint = "pubsub.googleapis.com:443"
const scope = "https://www.googleapis.com/auth/pubsub"

func Dial(ctx context.Context) (*grpc.ClientConn, error) {
	credentials, err := oauth.NewApplicationDefault(ctx, scope)
	if err != nil {
		panic(err)
	}

	grpcOpts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(credentials),
		//grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
		grpc.WithBlock(),
	}
	return grpc.DialContext(ctx, endpoint, grpcOpts...)
}

// TopicPath returns the fully-qualified Pub/Sub topic name, for use in gRPC messages.
func TopicPath(projectID string, topicID string) string {
	return "projects/" + projectID + "/topics/" + topicID
}

// TopicRoutingCtx adds "x-goog-request-params" gRPC metadata, the same as the official client.
func TopicRoutingCtx(ctx context.Context, topicPath string) context.Context {
	// from Google's *publisherGRPCClient.Publish:
	// https://github.com/googleapis/google-cloud-go/blob/main/pubsub/apiv1/publisher_client.go#L652
	// The values must be URL query parameter encoded:
	// See: https://google.aip.dev/client-libraries/4222
	md := metadata.Pairs("x-goog-request-params", "topic="+url.QueryEscape(topicPath))
	return metadata.NewOutgoingContext(ctx, md)
}
