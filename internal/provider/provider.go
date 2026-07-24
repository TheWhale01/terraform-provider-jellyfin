package provider

import (
	"context"
	"fmt"

	"github.com/TheWhale01/terraform-provider-jellyfin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type JellyfinProvider struct {
	version string
}

type JellyfinProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey types.String `tfsdk:"api_key"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JellyfinProvider {
			version: version,
		}
	}
}

func (p *JellyfinProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jellyfin"
	resp.Version = p.version
}

func (p *JellyfinProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema {
		Description: "The Jellyfin provider allows you to manage a Jellyfin instance.",
		Attributes: map[string]schema.Attribute {
			"endpoint": schema.StringAttribute {
				Description: "The URL of the the Jellyfin server.",
				Optional: true,
			},
			"api_key": schema.StringAttribute {
				Description: "The API key used to authenticate to the jellyfin server.",
				Optional: true,
				Sensitive: true,
			},
			"username": schema.StringAttribute {
				Description: "The username used to authenticate to the jellyfin server.",
				Optional: true,
				Sensitive: true,
			},
			"password": schema.StringAttribute {
				Description: "The password used to authenticate to the jellyfin server.",
				Optional: true,
				Sensitive: true,
			},
		},
	}
}

func (p *JellyfinProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config JellyfinProviderModel

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Endpoint.IsNull() || config.Endpoint.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create API client.",
			"Cannot use unknown values for endpoint or api_key.",
		)
		return
	}

	endpoint := config.Endpoint.ValueString()
	apiKey := config.APIKey.ValueString()
	apiClient := client.NewClient(endpoint, apiKey)
	if apiKey == "" {
		var accessToken, err = apiClient.AuthenticateByName(config.Username.ValueString(), config.Password.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("User auth failed", fmt.Sprintf("Jellyfin auth failed: %s", err.Error()))
			return
		}
		apiClient.APIKey = accessToken.AccessToken
	}
	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *JellyfinProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource {
		NewSystemInfoDataSource,
	}
}

func (p *JellyfinProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource {
		NewLibraryResource,
	}
}
