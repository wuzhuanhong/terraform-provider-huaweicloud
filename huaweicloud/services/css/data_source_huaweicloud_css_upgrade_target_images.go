// Generated by PMS #168
package css

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceCssUpgradeTargetImages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCssUpgradeTargetImagesRead,

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
				Description: `Specifies the ID of the cluster to be upgraded.`,
			},
			"upgrade_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the upgrade type.`,
			},
			"image_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the target image.`,
			},
			"engine_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the datastore type of the target image.`,
			},
			"engine_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the datastore version of the target image.`,
			},
			"images": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the upgrade target images.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The target image ID that can be upgraded.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the target image that can be upgraded.`,
						},
						"engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image datastore type.`,
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image datastore version.`,
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The target image priority.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image description information.`,
						},
					},
				},
			},
		},
	}
}

type UpgradeTargetImagesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newUpgradeTargetImagesDSWrapper(d *schema.ResourceData, meta interface{}) *UpgradeTargetImagesDSWrapper {
	return &UpgradeTargetImagesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCssUpgradeTargetImagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newUpgradeTargetImagesDSWrapper(d, meta)
	listImagesRst, err := wrapper.ListImages()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listImagesToSchema(listImagesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/target/{upgrade_type}/images
func (w *UpgradeTargetImagesDSWrapper) ListImages() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "css")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/clusters/{cluster_id}/target/{upgrade_type}/images"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	uri = strings.ReplaceAll(uri, "{upgrade_type}", w.Get("upgrade_type").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Filter(
			filters.New().From("imageInfoList").
				Where("id", "=", w.Get("image_id")).
				Where("datastoreType", "=", w.Get("engine_type")).
				Where("datastoreVersion", "=", w.Get("engine_version")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *UpgradeTargetImagesDSWrapper) listImagesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("images", schemas.SliceToList(body.Get("imageInfoList"),
			func(images gjson.Result) any {
				return map[string]any{
					"id":             images.Get("id").Value(),
					"name":           images.Get("displayName").Value(),
					"engine_type":    images.Get("datastoreType").Value(),
					"engine_version": images.Get("datastoreVersion").Value(),
					"priority":       images.Get("priority").Value(),
					"description":    images.Get("imageDesc").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
