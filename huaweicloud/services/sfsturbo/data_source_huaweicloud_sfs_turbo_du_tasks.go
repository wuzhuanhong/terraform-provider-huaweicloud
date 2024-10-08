// Generated by PMS #237
package sfsturbo

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

func DataSourceSfsTurboDuTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSfsTurboDuTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"share_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the SFS Turbo file system to which the DU tasks belong.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of DU tasks.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the DU task.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the DU task.`,
						},
						"dir_usage": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The resource usages of a directory (including subdirectories).`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The full path to a legal directory in the file system.`,
									},
									"used_capacity": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The used capacity, in byte.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The error message.`,
									},
									"file_count": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The total number of files in the directory.`,
										Elem:        tasksDirUsageFileCountElem(),
									},
								},
							},
						},
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the DU task, in RFC3339 format.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time of the DU task, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

// tasksDirUsageFileCountElem
// The Elem of "tasks.dir_usage.file_count"
func tasksDirUsageFileCountElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dir": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of directories.`,
			},
			"regular": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of common files.`,
			},
			"char": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of character devices.`,
			},
			"block": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of block devices.`,
			},
			"pipe": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of pipe files.`,
			},
			"socket": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of sockets.`,
			},
			"symlink": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of symbolic links.`,
			},
		},
	}
}

type TurboDuTasksDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newTurboDuTasksDSWrapper(d *schema.ResourceData, meta interface{}) *TurboDuTasksDSWrapper {
	return &TurboDuTasksDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceSfsTurboDuTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newTurboDuTasksDSWrapper(d, meta)
	listFsTasksRst, err := wrapper.ListFsTasks()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listFsTasksToSchema(listFsTasksRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API SFSTURBO GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/{feature}/tasks
func (w *TurboDuTasksDSWrapper) ListFsTasks() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "sfs-turbo")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/sfs-turbo/shares/{share_id}/fs/{feature}/tasks"
	uri = strings.ReplaceAll(uri, "{share_id}", w.Get("share_id").(string))
	uri = strings.ReplaceAll(uri, "{feature}", "dir-usage")
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		MarkerPager("tasks", "tasks[-1].task_id", "marker").
		Request().
		Result()
}

func (w *TurboDuTasksDSWrapper) listFsTasksToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("tasks", schemas.SliceToList(body.Get("tasks"),
			func(tasks gjson.Result) any {
				return map[string]any{
					"id":     tasks.Get("task_id").Value(),
					"status": tasks.Get("status").Value(),
					"dir_usage": schemas.SliceToList(tasks.Get("dir_usage"),
						func(dirUsage gjson.Result) any {
							return map[string]any{
								"path":          dirUsage.Get("path").Value(),
								"used_capacity": dirUsage.Get("used_capacity").Value(),
								"message":       dirUsage.Get("message").Value(),
								"file_count":    w.setTasDirUsaFilCount(dirUsage),
							}
						},
					),
					"begin_time": w.setTasksBeginTime(tasks),
					"end_time":   w.setTasksEndTime(tasks),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*TurboDuTasksDSWrapper) setTasDirUsaFilCount(dirUsage gjson.Result) any {
	return schemas.SliceToList(dirUsage.Get("file_count"), func(fileCount gjson.Result) any {
		return map[string]any{
			"dir":     fileCount.Get("dir").Value(),
			"regular": fileCount.Get("regular").Value(),
			"char":    fileCount.Get("char").Value(),
			"block":   fileCount.Get("block").Value(),
			"pipe":    fileCount.Get("pipe").Value(),
			"socket":  fileCount.Get("socket").Value(),
			"symlink": fileCount.Get("symlink").Value(),
		}
	})
}

func (*TurboDuTasksDSWrapper) setTasksBeginTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("begin_time").String(), "2006-01-02 15:04:05")/1000, false)
}

func (*TurboDuTasksDSWrapper) setTasksEndTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("end_time").String(), "2006-01-02 15:04:05")/1000, false)
}
