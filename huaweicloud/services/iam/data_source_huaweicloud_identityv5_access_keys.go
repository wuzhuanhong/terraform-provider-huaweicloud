package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/users/{user_id}/access-keys
func DataSourceV5AccessKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5AccessKeysRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the IAM user.`,
			},
			"access_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the permanent access key (AK).`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the IAM user.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the access key.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the access key.`,
						},
					},
				},
				Description: `The list of the permanent access keys.`,
			},
		},
	}
}

func listV5AccessKeys(client *golangsdk.ServiceClient, userId string) ([]interface{}, error) {
	var (
		httpUrl = "v5/users/{user_id}/access-keys"
		result  = make([]interface{}, 0)
		marker  = ""
		// The default limit is 200, maximum is 200.
		limit = 200
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{user_id}", userId)
	listPath += fmt.Sprintf("?limit=%d", limit)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		accessKeys := utils.PathSearch("access_keys", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, accessKeys...)
		if len(accessKeys) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceV5AccessKeysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	userId := d.Get("user_id").(string)
	accessKeys, err := listV5AccessKeys(client, userId)
	if err != nil {
		return diag.Errorf("error querying permanent access keys under user (%s): %s", userId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID)

	return diag.FromErr(d.Set("access_keys", flattenV5AccessKeys(accessKeys)))
}

func flattenV5AccessKeys(accessKeys []interface{}) []interface{} {
	if len(accessKeys) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(accessKeys))
	for _, accessKey := range accessKeys {
		result = append(result, map[string]interface{}{
			"access_key_id": utils.PathSearch("access_key_id", accessKey, nil),
			"user_id":       utils.PathSearch("user_id", accessKey, nil),
			"created_at":    utils.PathSearch("created_at", accessKey, nil),
			"status":        utils.PathSearch("status", accessKey, nil),
		})
	}
	return result
}
