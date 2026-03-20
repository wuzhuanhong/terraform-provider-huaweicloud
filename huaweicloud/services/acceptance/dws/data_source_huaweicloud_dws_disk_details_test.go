package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDiskDetails_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dws_disk_details.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byInstanceId   = "data.huaweicloud_dws_disk_details.filter_by_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byInstanceName   = "data.huaweicloud_dws_disk_details.filter_by_instance_name"
		dcByInstanceName = acceptance.InitDataSourceCheck(byInstanceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDiskDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "disk_details.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "disk_details.0.instance_name"),
					resource.TestCheckResourceAttrSet(all, "disk_details.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "disk_details.0.host_name"),
					resource.TestCheckResourceAttrSet(all, "disk_details.0.disk_name"),
					resource.TestCheckResourceAttrSet(all, "disk_details.0.disk_type"),
					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_id_filter_useful", "true"),
					dcByInstanceName.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDiskDetails_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_disk_details" "test" {
  cluster_id = "%[1]s"
}

# Filter by instance id.
locals {
  first_instance_id = data.huaweicloud_dws_disk_details.test.disk_details[0].instance_id
}

data "huaweicloud_dws_disk_details" "filter_by_instance_id" {
  cluster_id  = "%[1]s"
  instance_id = local.first_instance_id
}

locals {
  instance_id_filter_result = [
    for v in data.huaweicloud_dws_disk_details.filter_by_instance_id.disk_details[*].instance_id
    : v == local.first_instance_id
  ]
}

output "is_instance_id_filter_useful" {
  value = length(local.instance_id_filter_result) > 0 && alltrue(local.instance_id_filter_result)
}

# Filter by instance name.
locals {
  first_instance_name = data.huaweicloud_dws_disk_details.test.disk_details[0].instance_name
}

data "huaweicloud_dws_disk_details" "filter_by_instance_name" {
  cluster_id     = "%[1]s"
  instance_name  = local.first_instance_name
}

locals {
  instance_name_filter_result = [
    for v in data.huaweicloud_dws_disk_details.filter_by_instance_name.disk_details[*].instance_name
    : v == local.first_instance_name
  ]
}

output "is_instance_name_filter_useful" {
  value = length(local.instance_name_filter_result) > 0 && alltrue(local.instance_name_filter_result)
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
