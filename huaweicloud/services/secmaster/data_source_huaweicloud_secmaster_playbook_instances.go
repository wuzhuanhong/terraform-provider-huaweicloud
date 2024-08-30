// Generated by PMS #315
package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceSecmasterPlaybookInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterPlaybookInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID.`,
			},
			"from_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the search start time.`,
			},
			"to_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the search end time.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the playbook instance status.`,
			},
			"data_class_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the data class name.`,
			},
			"data_object_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the data object name.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the triggering type.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The playbook instance list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The playbook instance ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The playbook instance name.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The playbook instance status.`,
						},
						"trigger_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The triggering type.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The project ID.`,
						},
						"playbook": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The playbook information of the instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook ID.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook name.`,
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook version.`,
									},
									"version_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook version ID.`,
									},
								},
							},
						},
						"data_object": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The data object of the instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The data object ID of the instance.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The data object name of the instance.`,
									},
									"created_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The creation time of the data object.`,
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The update time of the data object.`,
									},
									"project_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The project ID of the data object.`,
									},
									"data_class_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The data class ID of the data object.`,
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The data content of the data object.`,
									},
								},
							},
						},
						"data_class": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The data class of the instance.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The data class ID of the instance.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The data class name of the instance.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type PlaybookInstancesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newPlaybookInstancesDSWrapper(d *schema.ResourceData, meta interface{}) *PlaybookInstancesDSWrapper {
	return &PlaybookInstancesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceSecmasterPlaybookInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newPlaybookInstancesDSWrapper(d, meta)
	lisPlaInsRst, err := wrapper.ListPlaybookInstances()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listPlaybookInstancesToSchema(lisPlaInsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances
func (w *PlaybookInstancesDSWrapper) ListPlaybookInstances() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "secmaster")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/instances"
	uri = strings.ReplaceAll(uri, "{workspace_id}", w.Get("workspace_id").(string))
	params := map[string]any{
		"status":          w.Get("status"),
		"dataclass_name":  w.Get("data_class_name"),
		"dataobject_name": w.Get("data_object_name"),
		"trigger_type":    w.Get("trigger_type"),
		"from_date":       w.Get("from_date"),
		"to_date":         w.Get("to_date"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("instances", "offset", "limit", 100).
		Request().
		Result()
}

func (w *PlaybookInstancesDSWrapper) listPlaybookInstancesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("instances", schemas.SliceToList(body.Get("instances"),
			func(instances gjson.Result) any {
				return map[string]any{
					"id":           instances.Get("id").Value(),
					"name":         instances.Get("name").Value(),
					"status":       instances.Get("status").Value(),
					"trigger_type": instances.Get("trigger_type").Value(),
					"created_at":   instances.Get("start_time").Value(),
					"updated_at":   instances.Get("end_time").Value(),
					"project_id":   instances.Get("project_id").Value(),
					"playbook": schemas.SliceToList(instances.Get("playbook"),
						func(playbook gjson.Result) any {
							return map[string]any{
								"id":         playbook.Get("id").Value(),
								"name":       playbook.Get("name").Value(),
								"version":    playbook.Get("version").Value(),
								"version_id": playbook.Get("version_id").Value(),
							}
						},
					),
					"data_object": schemas.SliceToList(instances.Get("dataobject"),
						func(dataObject gjson.Result) any {
							return map[string]any{
								"id":            dataObject.Get("id").Value(),
								"name":          dataObject.Get("name").Value(),
								"created_at":    dataObject.Get("create_time").Value(),
								"updated_at":    dataObject.Get("update_time").Value(),
								"project_id":    dataObject.Get("project_id").Value(),
								"data_class_id": dataObject.Get("dataclass_id").Value(),
								"content":       dataObject.Get("content").Value(),
							}
						},
					),
					"data_class": schemas.SliceToList(instances.Get("dataclass"),
						func(dataClass gjson.Result) any {
							return map[string]any{
								"id":   dataClass.Get("id").Value(),
								"name": dataClass.Get("name").Value(),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}