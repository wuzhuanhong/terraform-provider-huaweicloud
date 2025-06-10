package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS POST /v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/content/query
func DataSourceLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceLogsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log group to which the log to be queried belongs.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log stream to which the log to be queried belongs.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the log to be queried, in milliseconds.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the log to be queried, in milliseconds.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The label list of the log to be queried.`,
			},
			"keywords": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The keywords of the log to be queried.`,
			},
			"highlight": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to highlight the keywords.`,
			},
			"is_desc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to sort logs in descending order.`,
			},
			"__time__": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom time of the log to be queried.`,
			},
			"is_iterative": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether the log query is iterative.`,
			},
			"logs": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The content of the log.`,
						},
						"line_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The line number of the log.`,
						},
						"labels": {
							Type:        schema.TypeMap,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: `The labels associated with the log.`,
						},
					},
				},
				Computed:    true,
				Description: `The list of logs.`,
			},
		},
	}
}

func resourceLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	logs, err := queryLogsWithMarker(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("logs", flattenLogs(logs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryLogsBodyParams(d *schema.ResourceData, lineNum string, limit int) map[string]interface{} {
	return map[string]interface{}{
		"start_time":   d.Get("start_time"),
		"end_time":     d.Get("end_time"),
		"line_num":     utils.ValueIgnoreEmpty(lineNum),
		"limit":        limit,
		"labels":       utils.ValueIgnoreEmpty(d.Get("labels")),
		"keywords":     utils.ValueIgnoreEmpty(d.Get("keywords")),
		"highlight":    d.Get("highlight"),
		"is_desc":      d.Get("is_desc"),
		"__time__":     utils.ValueIgnoreEmpty(d.Get("__time__")),
		"is_iterative": d.Get("is_iterative"),
	}
}

func queryLogsWithMarker(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/content/query"
		limit   = 100
		marker  = ""
		result  = make([]interface{}, 0)
		getOpt  = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		}
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{log_group_id}", d.Get("log_group_id").(string))
	getPath = strings.ReplaceAll(getPath, "{log_stream_id}", d.Get("log_stream_id").(string))

	for {
		getOpt.JSONBody = utils.RemoveNil(buildQueryLogsBodyParams(d, marker, limit))
		resp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return nil, fmt.Errorf("error querying LTS logs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		logs := utils.PathSearch("logs", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, logs...)
		if len(logs) < limit {
			break
		}

		marker = utils.PathSearch("[-1].line_num", logs, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func flattenLogs(logs []interface{}) []interface{} {
	if len(logs) == 0 {
		return nil
	}

	result := make([]interface{}, len(logs))
	for i, log := range logs {
		result[i] = map[string]interface{}{
			"content":  utils.PathSearch("content", log, nil),
			"line_num": utils.PathSearch("line_num", log, nil),
			"labels":   utils.PathSearch("labels", log, nil),
		}
	}

	return result
}
