package provider

import (
	"context"
	"fmt"

	"github.com/TheWhale01/terraform-provider-jellyfin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SystemInfoDataSource struct {
	client *client.Client
}

type SystemInfoDataSourceModel struct {
	Id types.String `tfsdk:"id"`
	Version types.String `tfsdk:"version"`
}

func NewSystemInfoDataSource() datasource.DataSource {
	return &SystemInfoDataSource {}
}

func (d *SystemInfoDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_system_info"
}

func (d *SystemInfoDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema {
		Description: "Retrieves system information from the Jellyfin server.",
		Attributes: map[string]schema.Attribute {
			"id": schema.StringAttribute {
				Description: "The unique server identifier.",
				Computed: true,
			},
			"version": schema.StringAttribute {
				Description: "The Jellyfin server version.",
				Computed: true,
			},
		},
	}
}

func (d *SystemInfoDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}
	d.client = apiClient
}

func (d *SystemInfoDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SystemInfoDataSourceModel

	info, err := d.client.GetSystemInfo()
	if err != nil {
		resp.Diagnostics.AddError("Failed to fetch system info", err.Error())
		return
	}

	state.Id = types.StringValue(info.Id)
	state.Version = types.StringValue(info.Version)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
