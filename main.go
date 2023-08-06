package main

import (
	"context"
	"fmt"
	"github.com/aserto-dev/go-aserto/authorizer/grpc"
	"github.com/aserto-dev/go-aserto/client"
	authz "github.com/aserto-dev/go-authorizer/aserto/authorizer/v2"
	"github.com/aserto-dev/go-authorizer/aserto/authorizer/v2/api"
)

func main() {
	ctx := context.Background()

	authorizer, err := grpc.New(
		ctx,
		//client.WithAPIKeyAuth("<API Key>"),
		client.WithAddr("localhost:8282"),
		client.WithInsecure(true),
	)

	if err != nil {
		panic(err)
	}
	resp, err := authorizer.Is(ctx, &authz.IsRequest{
		PolicyContext: &api.PolicyContext{
			Path:      "peoplefinder.GET.users.__id",
			Decisions: []string{"allowed"},
		},
		IdentityContext: &api.IdentityContext{
			Identity: "<user name>",
			Type:     api.IdentityType_IDENTITY_TYPE_SUB,
		},
	})

	fmt.Println(resp)
}
