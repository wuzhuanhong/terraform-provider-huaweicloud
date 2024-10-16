package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLogicalClusterRings_basic(t *testing.T) {
	rName := "data.huaweicloud_dws_logical_cluster_rings.test"

	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsLogicalModeClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceLogicalClusterRings_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "cluster_rings.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.is_available"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.host_name"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.back_ip"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.cpu_cores"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.memory"),
					resource.TestCheckResourceAttrSet(rName, "cluster_rings.0.ring_hosts.0.disk_size"),
				),
			},
		},
	})
}

func testAccDatasourceLogicalClusterRings_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_cluster_rings" "test" {
  cluster_id = "%s"
}
`, acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID)
}
