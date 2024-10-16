package dws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDisasterRecoveryTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v2/{project_id}/disaster-recovery/{disaster_recovery_id}"
		product = "dws"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DWS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{disaster_recovery_id}", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DWS disaster recovery: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

// Before creating a disaster recovery task, ensure that the following conditions are consistent between the two clusters:
// 1. Availability zone.
// 2. Cluster type.
// 3. Node specification.
// 4. Number of nodes.
// 5. Virtual private cloud.
func TestAccResourceDisasterRecoveryTask_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_dws_disaster_recovery_task.test"
		name         = acceptance.RandomAccResourceName()
		password     = acceptance.RandomPassword() + "a"
	)
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDisasterRecoveryTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
			acceptance.TestAccPreCheckDwsLogicalModeClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAcDisasterRecoveryTask_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dr_sync_period", "2H"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "primary_cluster.0.id", acceptance.HW_DWS_CLUSTER_ID),
					resource.TestCheckResourceAttrSet(resourceName, "primary_cluster.0.name"),
					resource.TestCheckResourceAttr(resourceName, "standby_cluster.0.id", acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID),
					resource.TestCheckResourceAttrSet(resourceName, "standby_cluster.0.name"),
				),
			},
			{
				Config: testAcDisasterRecoveryTask_update(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "dr_sync_period", "3H"),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
				),
			},
			{
				Config: testAcDisasterRecoveryTask_switch(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "running"),
					resource.TestCheckResourceAttrPair(resourceName, "primary_cluster.0.id", resourceName, "standby_cluster_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action", "primary_cluster_id", "standby_cluster_id"},
			},
		},
	})
}

func testAcDisasterRecoveryTask_basic(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_disaster_recovery_task" "test" {
  name               = "%[1]s"
  dr_type            = "az"
  primary_cluster_id = "%[2]s"
  standby_cluster_id = "%[3]s"
  dr_sync_period     = "2H"
}
`, name, acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID)
}

func testAcDisasterRecoveryTask_update(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_disaster_recovery_task" "test" {
  name               = "%[1]s"
  dr_type            = "az"
  primary_cluster_id = "%[2]s"
  standby_cluster_id = "%[3]s"
  dr_sync_period     = "3H"
  action             = "start"
}
`, name, acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID)
}

func testAcDisasterRecoveryTask_switch(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dws_disaster_recovery_task" "test" {
  name               = "%[1]s"
  dr_type            = "az"
  primary_cluster_id = "%[2]s"
  standby_cluster_id = "%[3]s"
  dr_sync_period     = "3H"
  action             = "switchover"
}
`, name, acceptance.HW_DWS_CLUSTER_ID, acceptance.HW_DWS_LOGICAL_MODE_CLUSTER_ID)
}
