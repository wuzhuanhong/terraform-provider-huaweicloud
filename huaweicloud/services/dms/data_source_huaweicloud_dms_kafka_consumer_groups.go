// Generated by PMS #199
package dms

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDmsKafkaConsumerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsKafkaConsumerGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the consumer group name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the consumer group description.`,
			},
			"lag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the the number of accumulated messages.`,
			},
			"coordinator_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the coordinator ID.`,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the consumer group status.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the consumer groups.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the consumer group name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the consumer group description.`,
						},
						"lag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of accumulated messages.`,
						},
						"coordinator_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the coordinator ID.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the consumer group status.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
					},
				},
			},
		},
	}
}

type KafkaConsumerGroupsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newKafkaConsumerGroupsDSWrapper(d *schema.ResourceData, meta interface{}) *KafkaConsumerGroupsDSWrapper {
	return &KafkaConsumerGroupsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDmsKafkaConsumerGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newKafkaConsumerGroupsDSWrapper(d, meta)
	lisInsConGroRst, err := wrapper.ListInstanceConsumerGroups()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listInstanceConsumerGroupsToSchema(lisInsConGroRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/groups
func (w *KafkaConsumerGroupsDSWrapper) ListInstanceConsumerGroups() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dmsv2")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/instances/{instance_id}/groups"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"group": w.Get("name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Filter(
			filters.New().From("groups").
				Where("state", "=", w.Get("state")).
				Where("group_desc", "contains", w.Get("description")).
				Where("coordinator_id", "=", w.Get("coordinator_id")).
				Where("lag", "=", w.Get("lag")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *KafkaConsumerGroupsDSWrapper) listInstanceConsumerGroupsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("groups", schemas.SliceToList(body.Get("groups"),
			func(groups gjson.Result) any {
				return map[string]any{
					"name":           groups.Get("group_id").Value(),
					"description":    groups.Get("group_desc").Value(),
					"lag":            groups.Get("lag").Value(),
					"coordinator_id": groups.Get("coordinator_id").Value(),
					"state":          groups.Get("state").Value(),
					"created_at":     w.setGroupsCreatedAt(groups),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*KafkaConsumerGroupsDSWrapper) setGroupsCreatedAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339((data.Get("createdAt").Int())/1000, true)
}
