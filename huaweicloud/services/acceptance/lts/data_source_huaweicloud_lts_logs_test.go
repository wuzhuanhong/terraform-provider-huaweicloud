package lts

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceLogs_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_lts_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLabels   = "data.huaweicloud_lts_logs.filter_by_labels"
		dcByLabels = acceptance.InitDataSourceCheck(byLabels)

		byKeywordsAndHighlight   = "data.huaweicloud_lts_logs.filter_by_keywords_and_highlight"
		dcByKeywordsAndHighlight = acceptance.InitDataSourceCheck(byKeywordsAndHighlight)

		byKeywordsAndNotHighlight   = "data.huaweicloud_lts_logs.filter_by_keywords_and_not_highlight"
		dcByKeywordsAndNotHighlight = acceptance.InitDataSourceCheck(byKeywordsAndNotHighlight)

		byDesc   = "data.huaweicloud_lts_logs.filter_by_desc"
		dcByDesc = acceptance.InitDataSourceCheck(byDesc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLogs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "logs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByLabels.CheckResourceExists(),
					resource.TestCheckOutput("is_labels_useful", "true"),
					dcByKeywordsAndHighlight.CheckResourceExists(),
					resource.TestCheckOutput("is_keywords_and_highlight_filter_useful", "true"),
					dcByKeywordsAndNotHighlight.CheckResourceExists(),
					resource.TestCheckOutput("is_keywords_and_not_highlight_filter_useful", "true"),
					dcByDesc.CheckResourceExists(),
					resource.TestCheckOutput("is_desc_useful", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.labels.%"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.line_num"),
				),
			},
		},
	})
}

func testAccDataSourceLogs_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_group" "test" {
  group_name  = "%[2]s"
  ttl_in_days = 1
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
}

resource "huaweicloud_lts_stream_index_configuration" "test" {
  group_id  = huaweicloud_lts_group.test.id
  stream_id = huaweicloud_lts_stream.test.id

  full_text_index {
    enable          = true
    case_sensitive  = true
    include_chinese = true
    tokenizer       = ", '\";=()[]{}@&<>/:\\n\\t\\r"
  }

  fields {
    field_type      = "string"
    field_name      = "action"
    case_sensitive  = false
    include_chinese = false
    tokenizer       = ", '\";=()[]{}@&<>/:\\n\\t\\r"
    quick_analysis  = true
  }
  fields {
    field_type      = "string"
    field_name      = "project_id"
    case_sensitive  = false
    include_chinese = false
    tokenizer       = ", '\";=()[]{}@&<>/:\\n\\t\\r"
    quick_analysis  = true
  }
  fields {
    field_type      = "string"
    field_name      = "log_status"
    case_sensitive  = false
    include_chinese = false
    tokenizer       = ", '\";=()[]{}@&<>/:\\n\\t\\r"
    quick_analysis  = true
  }
  fields {
    field_type      = "string"
    field_name      = "hostIP"
    case_sensitive  = false
    include_chinese = false
    tokenizer       = ", '\";=()[]{}@&<>/:\\n\\t\\r"
    quick_analysis  = true
  }
  fields {
    field_type      = "string"
    field_name      = "hostName"
    case_sensitive  = false
    include_chinese = false
    tokenizer       = ", '\";=()[]{}@&<>/:\\n\\t\\r"
    quick_analysis  = true
  }
}

resource "huaweicloud_vpc_flow_log" "test" {
  name          = "%[2]s"
  resource_type = "vpc"
  resource_id   = huaweicloud_vpc.test.id
  traffic_type  = "all"
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  visibility = "public"
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = try(data.huaweicloud_images_images.test.images[0].id, "")
  flavor_id          = try(data.huaweicloud_compute_flavors.test.ids[0], "")
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = try(data.huaweicloud_availability_zones.test.names[0], "")
  power_action       = "REBOOT"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

# Waiting for logs to be generated.
resource "null_resource" "test" {
  depends_on = [
    huaweicloud_compute_instance.test,
    huaweicloud_lts_stream_index_configuration.test
  ]

  provisioner "local-exec" {
    command = "sleep 600"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataSourceLogs_basic(name string) string {
	currentTime := time.Now()
	startTime := currentTime.UnixMilli()
	endTime := currentTime.Add(1 * time.Hour).UnixMilli()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_logs" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = tostring("%[2]v")
  end_time      = tostring("%[3]v")

  depends_on    = [null_resource.test]
}

# Filter by labels.
data "huaweicloud_lts_logs" "filter_by_labels" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = tostring("%[2]v")
  end_time      = tostring("%[3]v")

  labels = {
    action = "ACCEPT"
  }

  depends_on    = [null_resource.test]
}

locals {
  labels_filter_result = [
    for v in data.huaweicloud_lts_logs.filter_by_labels.logs[*].labels : strcontains(v.action, "ACCEPT")
  ]
}

output "is_labels_useful" {
  value = length(local.labels_filter_result) > 0 && alltrue(local.labels_filter_result)
}

# Filter by keywords and highlight (the value is true).
data "huaweicloud_lts_logs" "filter_by_keywords_and_highlight" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = tostring("%[2]v")
  end_time      = tostring("%[3]v")
  keywords      = "action : ACCEPT"

  depends_on    = [null_resource.test]
}

locals {
  keywords_and_highlight_filter_result = [
    for v in data.huaweicloud_lts_logs.filter_by_keywords_and_highlight.logs[*].labels :
    v.action == "<HighLightTag>ACCEPT</HighLightTag>"
  ]
}

output "is_keywords_and_highlight_filter_useful" {
  value = length(local.keywords_and_highlight_filter_result) > 0 && alltrue(local.keywords_and_highlight_filter_result)
}

# Filter by keywords and highlight (the value is false).
data "huaweicloud_lts_logs" "filter_by_keywords_and_not_highlight" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = tostring("%[2]v")
  end_time      = tostring("%[3]v")
  keywords      = "action : ACCEPT"
  highlight     = false

  depends_on    = [null_resource.test]
}

locals {
  keywords_and_not_highlight_filter_result = [
    for v in data.huaweicloud_lts_logs.filter_by_keywords_and_not_highlight.logs[*].labels : v.action == "ACCEPT"
  ]
}

output "is_keywords_and_not_highlight_filter_useful" {
  value = length(local.keywords_and_not_highlight_filter_result) > 0 && alltrue(local.keywords_and_not_highlight_filter_result)
}

# Filter by is_desc.
data "huaweicloud_lts_logs" "filter_by_desc" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  start_time    = tostring("%[2]v")
  end_time      = tostring("%[3]v")
  is_desc       = true

  depends_on    = [null_resource.test]
}

locals {
  asc_line_num       = try(data.huaweicloud_lts_logs.test.logs[0].line_num, "")
  desc_filter_result = data.huaweicloud_lts_logs.filter_by_desc.logs
  desc_line_num      = try(local.desc_filter_result[length(local.desc_filter_result) - 1].line_num, "")
}

output "is_desc_useful" {
  value = local.asc_line_num == local.desc_line_num
}
`, testAccDataSourceLogs_base(name), startTime, endTime)
}
