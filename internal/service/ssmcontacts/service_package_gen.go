// Code generated by internal/generate/servicepackage/main.go; DO NOT EDIT.

package ssmcontacts

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssmcontacts"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

type servicePackage struct{}

func (p *servicePackage) FrameworkDataSources(ctx context.Context) []*types.ServicePackageFrameworkDataSource {
	return []*types.ServicePackageFrameworkDataSource{
		{
			Factory:  newDataSourceRotation,
			TypeName: "aws_ssmcontacts_rotation",
			Name:     "Rotation",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) FrameworkResources(ctx context.Context) []*types.ServicePackageFrameworkResource {
	return []*types.ServicePackageFrameworkResource{
		{
			Factory:  newResourceRotation,
			TypeName: "aws_ssmcontacts_rotation",
			Name:     "Rotation",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
	}
}

func (p *servicePackage) SDKDataSources(ctx context.Context) []*types.ServicePackageSDKDataSource {
	return []*types.ServicePackageSDKDataSource{
		{
			Factory:  DataSourceContact,
			TypeName: "aws_ssmcontacts_contact",
			Name:     "Contact",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  DataSourceContactChannel,
			TypeName: "aws_ssmcontacts_contact_channel",
			Name:     "Contact Channel",
		},
		{
			Factory:  DataSourcePlan,
			TypeName: "aws_ssmcontacts_plan",
			Name:     "Plan",
		},
	}
}

func (p *servicePackage) SDKResources(ctx context.Context) []*types.ServicePackageSDKResource {
	return []*types.ServicePackageSDKResource{
		{
			Factory:  ResourceContact,
			TypeName: "aws_ssmcontacts_contact",
			Name:     "Contact",
			Tags: &types.ServicePackageResourceTags{
				IdentifierAttribute: names.AttrARN,
			},
		},
		{
			Factory:  ResourceContactChannel,
			TypeName: "aws_ssmcontacts_contact_channel",
			Name:     "Contact Channel",
		},
		{
			Factory:  ResourcePlan,
			TypeName: "aws_ssmcontacts_plan",
			Name:     "Plan",
		},
	}
}

func (p *servicePackage) ServicePackageName() string {
	return names.SSMContacts
}

// NewClient returns a new AWS SDK for Go v2 client for this service package's AWS API.
func (p *servicePackage) NewClient(ctx context.Context, config map[string]any) (*ssmcontacts.Client, error) {
	cfg := *(config["aws_sdkv2_config"].(*aws.Config))

	return ssmcontacts.NewFromConfig(cfg,
		ssmcontacts.WithEndpointResolverV2(newEndpointResolverV2()),
		withBaseEndpoint(config[names.AttrEndpoint].(string)),
	), nil
}

func ServicePackage(ctx context.Context) conns.ServicePackage {
	return &servicePackage{}
}
