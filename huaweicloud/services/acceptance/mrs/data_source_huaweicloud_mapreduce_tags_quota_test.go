package mrs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTagsQuota_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_mapreduce_tags_quota.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTagsQuota_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "total_quota", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr(all, "available_quota", regexp.MustCompile(`^\d+$`)),
				),
			},
		},
	})
}

func testAccDataTagsQuota_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_mapreduce_tags_quota" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_MRS_CLUSTER_ID)
}
