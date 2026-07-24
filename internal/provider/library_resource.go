package provider

import (
	"fmt"
	"context"

	"github.com/TheWhale01/terraform-provider-jellyfin/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LibraryResource struct {
	client *client.Client
}

type LibraryResourceModel struct {
	Id types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	CollectionType types.String `tfsdk:"collection_type"`
	Paths []types.String `tfsdk:"paths"`
}

func NewLibraryResource() resource.Resource {
	return &LibraryResource{}
}

func (r *LibraryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_library"
}

func (r *LibraryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema {
		Description: "Retrieves system information from the Jellyfin server.",
		Attributes: map[string]schema.Attribute {
			"id": schema.StringAttribute {
				Description: "The internal Id of the library.",
				Computed: true,
				PlanModifiers: []planmodifier.String {
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute {
				Description: "The name of the library.",
				Required: true,
				PlanModifiers: []planmodifier.String {
					stringplanmodifier.RequiresReplace(),
				},
			},
			"collection_type": schema.StringAttribute {
				Description: "The type of the media (movies, tvshows, music, musicvideos, homevideos, boxsets, books, mixed)",
				Required: true,
			},
			"paths": schema.ListAttribute {
				Description: "List of the absolute paths of the library",
				ElementType: types.StringType,
				Required: true,
			},
		},
	}
}

func (r *LibraryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	r.client = apiClient
}

func (r *LibraryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan LibraryResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	var paths []string
	for _, p := range plan.Paths {
		paths = append(paths, p.ValueString())
	}
	options := client.VirtualFolder {
		Name: plan.Name.ValueString(),
		CollectionType: plan.CollectionType.ValueString(),
		Locations: paths,
	}
	if err := r.client.AddLibrary(options); err != nil {
		resp.Diagnostics.AddError("Error creating library.", err.Error())
		return
	}
	folders, err := r.client.GetLibraries()
	if err != nil {
		resp.Diagnostics.AddError("Error retriving libraries", err.Error())
		return
	}
	for _, folder := range folders {
		if folder.Name == options.Name {
			plan.Id = types.StringValue(folder.ItemId)
			break
		}
	}
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *LibraryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state LibraryResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	folders, err := r.client.GetLibraries()
	if err != nil {
		resp.Diagnostics.AddError("Error reading libraries.", err.Error())
		return
	}
	found := false
	for _, folder := range folders {
		if folder.ItemId == state.Id.ValueString() {
			found = true
			break
		}
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *LibraryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		"Update not fully implemented",
		"Work in progress...",
	)
	var plan LibraryResourceModel
	req.Plan.Get(ctx, &plan)
	resp.State.Set(ctx, &plan)
}

func (r *LibraryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state LibraryResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.RemoveLibrary(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting library", err.Error())
		return
	}
}
