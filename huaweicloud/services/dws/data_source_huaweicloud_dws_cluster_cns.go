// Generated by PMS #343
package dws

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
)

func DataSourceDwsClusterCns() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDwsClusterCnsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DWS cluster ID to which the CNs belong.`,
			},
			"min_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum number of CNs supported by the cluster.`,
			},
			"max_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum number of CNs supported by the cluster.`,
			},
			"cns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the CNs under specified DWS cluster.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the CN.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the CN.`,
						},
						"private_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The private IP address of the CN.`,
						},
					},
				},
			},
		},
	}
}

type ClusterCnsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newClusterCnsDSWrapper(d *schema.ResourceData, meta interface{}) *ClusterCnsDSWrapper {
	return &ClusterCnsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDwsClusterCnsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newClusterCnsDSWrapper(d, meta)
	listClusterCnRst, err := wrapper.ListClusterCn()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listClusterCnToSchema(listClusterCnRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}/cns
func (w *ClusterCnsDSWrapper) ListClusterCn() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dws")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/clusters/{cluster_id}/cns"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *ClusterCnsDSWrapper) listClusterCnToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("min_num", body.Get("min_num").Value()),
		d.Set("max_num", body.Get("max_num").Value()),
		d.Set("cns", schemas.SliceToList(body.Get("instances"),
			func(cns gjson.Result) any {
				return map[string]any{
					"id":         cns.Get("id").Value(),
					"name":       cns.Get("name").Value(),
					"private_ip": cns.Get("private_ip").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
