// Generated by PMS #140
package cfw

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

func DataSourceCfwDomainNameGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwDomainNameGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the protected object ID.`,
			},
			"key_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the key word.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the domain name group type.`,
			},
			"config_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the configuration status.`,
			},
			"ref_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the domain name group reference count.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the domain name group ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of a domain name group.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The domain name group list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name group ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name group name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name group description.`,
						},
						"ref_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The domain name group reference count.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name group type.`,
						},
						"config_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The configuration status.`,
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The configuration message.`,
						},
						"rules": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The used rule list.`,
							Elem:        dataRecordsRuleElem(),
						},
						"domain_names": {
							Type:        schema.TypeList,
							Elem:        domainGroupDomainNames(),
							Computed:    true,
							Description: `The list of domain names.`,
						},
					},
				},
			},
		},
	}
}

// dataRecordsRuleElem
// The Elem of "data.records.rules"
func dataRecordsRuleElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The rule ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The rule name.`,
			},
		},
	}
}

func domainGroupDomainNames() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description.`,
			},
			"domain_address_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain address ID.`,
			},
		},
	}
}

type DomainNameGroupsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newDomainNameGroupsDSWrapper(d *schema.ResourceData, meta interface{}) *DomainNameGroupsDSWrapper {
	return &DomainNameGroupsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCfwDomainNameGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newDomainNameGroupsDSWrapper(d, meta)
	lisDomSetRst, err := wrapper.ListDomainSets()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listDomainSetsToSchema(lisDomSetRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CFW GET /v1/{project_id}/domain-sets
func (w *DomainNameGroupsDSWrapper) ListDomainSets() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cfw")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/domain-sets"
	params := map[string]any{
		"enterprise_project_id": w.Get("enterprise_project_id"),
		"fw_instance_id":        w.Get("fw_instance_id"),
		"object_id":             w.Get("object_id"),
		"key_word":              w.Get("key_word"),
		"domain_set_type":       w.Get("type"),
		"config_status":         w.Get("config_status"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("data.records", "offset", "limit", 1024).
		Filter(
			filters.New().From("data.records").
				Where("set_id", "=", w.Get("group_id")).
				Where("name", "=", w.Get("name")).
				Where("ref_count", "=", w.Get("ref_count")),
		).
		Request().
		Result()
}

func (w *DomainNameGroupsDSWrapper) listDomainSetsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	domainGroups, err := w.flattenDomainSets(body)
	if err != nil {
		return err
	}
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("records", domainGroups),
	)
	return mErr.ErrorOrNil()
}

func (w *DomainNameGroupsDSWrapper) flattenDomainSets(body *gjson.Result) ([]map[string]any, error) {
	records := body.Get("data.records").Array()
	domainGroups := make([]map[string]any, 0, len(records))
	for _, record := range records {
		rawDomains, err := w.ListDomains(record.Get("set_id").String())
		if err != nil {
			return nil, err
		}

		domains := rawDomains.Get("data.records").Array()
		domainNames := make([]map[string]any, 0, len(domains))
		for _, domain := range domains {
			domainNames = append(domainNames, map[string]any{
				"domain_name":       domain.Get("domain_name").Value(),
				"description":       domain.Get("description").Value(),
				"domain_address_id": domain.Get("domain_address_id").Value(),
			})
		}

		domainGroups = append(domainGroups, map[string]any{
			"description":   record.Get("description").Value(),
			"ref_count":     record.Get("ref_count").Value(),
			"type":          record.Get("domain_set_type").String(),
			"config_status": record.Get("config_status").String(),
			"message":       record.Get("message").Value(),
			"group_id":      record.Get("set_id").Value(),
			"name":          record.Get("name").Value(),
			"rules": schemas.SliceToList(record.Get("rules"),
				func(rule gjson.Result) any {
					return map[string]any{
						"id":   rule.Get("id").Value(),
						"name": rule.Get("name").Value(),
					}
				},
			),
			"domain_names": domainNames,
		})
	}
	return domainGroups, nil
}

// @API CFW GET /v1/{project_id}/domain-set/domains/{domain_set_id}
func (w *DomainNameGroupsDSWrapper) ListDomains(setID string) (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cfw")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/domain-set/domains/{domain_set_id}"
	uri = strings.ReplaceAll(uri, "{domain_set_id}", setID)
	params := map[string]any{
		"fw_instance_id": w.Get("fw_instance_id"),
		"object_id":      w.Get("object_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("data.records", "offset", "limit", 1024).
		Request().
		Result()
}