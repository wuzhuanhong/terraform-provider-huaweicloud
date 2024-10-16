package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLogicalClusters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dws_logical_clusters.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsLogicalModeClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			// The cluster_id dose not exist (the cluster ID must be in the standard UUID format).
			{
				Config:      testDataSourceLogicalClusters_expectError(),
				ExpectError: regexp.MustCompile("DWS.0047"),
			},
			// The cluster_id dose not exist (in non-standard UUID format).
			{
				Config:      testDataSourceLogicalClusters_clusterIdNotExist,
				ExpectError: regexp.MustCompile("DWS.0001"),
			},
			{
				Config: testDataSourceLogicalClusters_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "logical_clusters.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "logical_clusters.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_clusters.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "logical_clusters.0.status"),
					resource.TestCheckOutput("is_existed_logical_cluster", "true"),
				),
			},
		},
	})
}

func testDataSourceLogicalClusters_expectError() string {
	randUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_clusters" "test" {
  cluster_id = "%s"
}
`, randUUID)
}

const testDataSourceLogicalClusters_clusterIdNotExist = `
data "huaweicloud_dws_logical_clusters" "test" {
  cluster_id = "not-found-cluster-id"
}
`

func testDataSourceLogicalClusters_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_logical_clusters" "test" {
  cluster_id = "%[1]s"
}

locals {
  logical_clusters      = data.huaweicloud_dws_logical_clusters.test.logical_clusters
  logical_cluster_names = local.logical_clusters[*].name
}

# Assert that the query results contain the current logical cluster ID and name.
# The "contains" method is an exact match.
output "is_existed_logical_cluster" {
  value = contains(local.logical_cluster_names, "%[2]s")
}
`, acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID, acceptance.HW_DWS_LOGICAL_CLUSTER_NAME)
}
