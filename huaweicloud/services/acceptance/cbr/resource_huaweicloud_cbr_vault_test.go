package cbr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
)

func getVaultResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cbr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CBR client: %s", err)
	}
	return cbr.GetVaultById(client, state.Primary.ID)
}

func TestAccVault_backupServer(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		basicConfig  = testAccVault_backupServer_base(name)

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupServer_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cloud_type", "public"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_bind", "false"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "bind_rules.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "policy.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "backup_name_prefix", "test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_az", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVault_backupServer_step2(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.server_id",
						"huaweicloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.0",
						"huaweicloud_evs_volume.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.1",
						"huaweicloud_evs_volume.test.2", "id"),
				),
			},
			{
				Config: testAccVault_backupServer_step3(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "backup_name_prefix", "test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "is_multi_az", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.server_id", "huaweicloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.0",
						"huaweicloud_evs_volume.test.2", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.1.server_id", "huaweicloud_compute_instance.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.1.excludes.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.1.excludes.0",
						"huaweicloud_evs_volume.test.1", "id"),
				),
			},
			{
				Config: testAccVault_backupServer_step4(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.server_id",
						"huaweicloud_compute_instance.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.0",
						"huaweicloud_evs_volume.test.1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.1",
						"huaweicloud_evs_volume.test.3", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resources_origin",
				},
			},
		},
	})
}

func testAccVault_backupServer_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  visibility   = "public"
  architecture = "x86"
  os           = "Ubuntu"
}

resource "huaweicloud_vpc" "test" {
  name                  = "%[2]s"
  cidr                  = "192.168.0.0/16"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[2]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                  = "%[2]s"
  delete_default_rules  = true
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name                  = format("%[2]s_%%d", count.index)
  image_id              = try(data.huaweicloud_images_images.test.images[0].id, null)
  flavor_id             = try(data.huaweicloud_compute_flavors.test.flavors[0].id, null)
  security_group_ids    = [huaweicloud_networking_secgroup.test.id]
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  system_disk_type = "SSD"
  system_disk_size = 40

  network { 
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_evs_volume" "test" {
  count = 4

  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  volume_type           = "SSD"
  name                  = format("%[2]s_data_%%d", count.index)
  size                  = 20
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_compute_volume_attach" "test" {
  count = 4

  instance_id = huaweicloud_compute_instance.test[count.index %% 2].id
  volume_id   = huaweicloud_evs_volume.test[count.index].id
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccVault_backupServer_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, basicConfig, name)
}

func testAccVault_backupServer_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true

  resources {
    server_id = huaweicloud_compute_instance.test[0].id
    excludes  = [
      huaweicloud_evs_volume.test[0].id,
      huaweicloud_evs_volume.test[2].id,
    ]
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, basicConfig, name)
}

func testAccVault_backupServer_step3(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true

  resources {
    server_id = huaweicloud_compute_instance.test[0].id
    excludes  = [
      huaweicloud_evs_volume.test[2].id,
    ]
  }
  resources {
    server_id = huaweicloud_compute_instance.test[1].id
    excludes  = [
      huaweicloud_evs_volume.test[1].id,
    ]
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, basicConfig, name)
}

func testAccVault_backupServer_step4(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  depends_on = [huaweicloud_compute_volume_attach.test]

  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
  backup_name_prefix    = "test-prefix-"
  is_multi_az           = true

  resources {
    server_id = huaweicloud_compute_instance.test[1].id
    excludes  = [
      huaweicloud_evs_volume.test[1].id,
      huaweicloud_evs_volume.test[3].id,
    ]
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, basicConfig, name)
}

func TestAccVault_replicationServer(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_replicationServer_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resources_origin",
				},
			},
		},
	})
}

func testAccVault_replicationServer_step1(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "replication"
  size                  = 200
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func TestAccVault_prePaidServer(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		basicConfig  = testAccVault_backupServer_base(name)

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_prePaidBackupServer_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVault_prePaidBackupServer_step2(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.server_id",
						"huaweicloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.0",
						"huaweicloud_compute_volume_attach.test.0", "volume_id"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.1",
						"huaweicloud_compute_volume_attach.test.2", "volume_id"),
				),
			},
			{
				Config: testAccVault_prePaidBackupServer_step3(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.server_id", "huaweicloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.0",
						"huaweicloud_compute_volume_attach.test.2", "volume_id"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.1.server_id", "huaweicloud_compute_instance.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.1.excludes.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.1.excludes.0",
						"huaweicloud_compute_volume_attach.test.1", "volume_id"),
				),
			},
			{
				Config: testAccVault_prePaidBackupServer_step4(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.server_id",
						"huaweicloud_compute_instance.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "resources.0.excludes.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.0",
						"huaweicloud_compute_volume_attach.test.1", "volume_id"),
					resource.TestCheckResourceAttrPair(resourceName, "resources.0.excludes.1",
						"huaweicloud_compute_volume_attach.test.3", "volume_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resources_origin",
					"period_unit",
					"period",
					"auto_renew",
					"auto_pay",
				},
			},
		},
	})
}

func testAccVault_prePaidBackupServer_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, basicConfig, name)
}

func testAccVault_prePaidBackupServer_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"

  resources {
    server_id = huaweicloud_compute_instance.test[0].id
    excludes  = [
      huaweicloud_compute_volume_attach.test[0].volume_id,
      huaweicloud_compute_volume_attach.test[2].volume_id,
    ]
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, basicConfig, name)
}

func testAccVault_prePaidBackupServer_step3(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  resources {
    server_id = huaweicloud_compute_instance.test[0].id
    excludes  = [
      huaweicloud_evs_volume.test[2].id,
    ]
  }
  resources {
    server_id = huaweicloud_compute_instance.test[1].id
    excludes  = [
      huaweicloud_evs_volume.test[1].id,
    ]
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, basicConfig, name)
}

func testAccVault_prePaidBackupServer_step4(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "true"

  resources {
    server_id = huaweicloud_compute_instance.test[1].id
    excludes  = [
      huaweicloud_evs_volume.test[1].id,
      huaweicloud_evs_volume.test[3].id,
    ]
  }

  tags = {
    foo1 = "bar"
    key  = "value_update"
  }
}
`, basicConfig, name)
}

func TestAccVault_volume(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		basicConfig  = testAccVault_volume_base(name)

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_volume_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
				),
			},
			{
				Config: testAccVault_volume_step2(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "true"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
				),
			},
			{
				Config: testAccVault_volume_step3(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"resources_origin",
				},
			},
		},
	})
}

func testAccVault_volume_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  count = 4

  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  volume_type           = "SSD"
  name                  = format("%[2]s_data_%%d", count.index)
  size                  = 20
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccVault_volume_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 50
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  resources {
    includes = slice(huaweicloud_evs_volume.test[*].id, 0, 2)
  }
}
`, basicConfig, name)
}

func testAccVault_volume_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 100
  auto_expand           = true
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  resources {
    includes = slice(huaweicloud_evs_volume.test[*].id, 1, 3)
  }
}
`, basicConfig, name)
}

func testAccVault_volume_step3(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 100
  auto_expand           = true
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  resources {
    includes = slice(huaweicloud_evs_volume.test[*].id, 1, 2)
  }
}
`, basicConfig, name)
}

func TestAccVault_backupTurbo(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupTurbo_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccVault_backupTurbo_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Vaults of type 'turbo'
func testAccVault_backupTurbo_base(name string) string {
	return fmt.Sprintf(`
%[1]s

variable "enterprise_project_id" {
  type    = string
  default = "%[2]s"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  count = 3

  name                  = format("%[3]s_%%d", count.index)
  size                  = 500
  share_proto           = "NFS"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}`, common.TestBaseNetwork(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccVault_backupTurbo_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 800
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  resources {
    includes = slice(huaweicloud_sfs_turbo.test[*].id, 0, 2)
  }
}
`, testAccVault_backupTurbo_base(rName), rName)
}

func testAccVault_backupTurbo_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 1000
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  resources {
    includes = slice(huaweicloud_sfs_turbo.test[*].id, 1, 3)
  }
}
`, testAccVault_backupTurbo_base(rName), rName)
}

func TestAccVault_replicationTurbo(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_replicationTurbo_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "replication"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVault_replicationTurbo_step1(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "turbo"
  protection_type       = "replication"
  size                  = 1000
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func TestAccVault_AutoBind(t *testing.T) {
	var (
		vault interface{}

		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_cbr_vault.test"

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_autoBind(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "auto_bind", "true"),
					resource.TestCheckResourceAttr(resourceName, "bind_rules.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "bind_rules.key", "value"),
				),
			},
			{
				Config: testAccVault_autoBindUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "auto_bind", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVault_autoBind(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  protection_type       = "backup"
  size                  = 800
  auto_bind             = true
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null

  bind_rules = {
    foo = "bar"
    key = "value"
  }
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccVault_autoBindUpdate(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type    = string
  default = "%[1]s"
}

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "server"
  protection_type       = "backup"
  size                  = 800
  auto_bind             = false
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func TestAccVault_bindPolicies(t *testing.T) {
	var (
		vault interface{}

		randName     = acceptance.RandomAccResourceName()
		mainRcName   = "huaweicloud_cbr_vault.test"
		legacyRcName = "huaweicloud_cbr_vault.legacy"

		mainRc   = acceptance.InitResourceCheck(mainRcName, &vault, getVaultResourceFunc)
		legacyRc = acceptance.InitResourceCheck(legacyRcName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckReplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      mainRc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_bindPolicies_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "policy_id",
						"huaweicloud_cbr_policy.backup.0", "id"),
				),
			},
			{
				Config: testAccVault_bindPolicies_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					mainRc.CheckResourceExists(),
					resource.TestCheckResourceAttr(mainRcName, "name", randName),
					resource.TestCheckResourceAttr(mainRcName, "policy.#", "2"),
					legacyRc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(legacyRcName, "policy_id",
						"huaweicloud_cbr_policy.backup.1", "id"),
				),
			},
			{
				ResourceName:      mainRcName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      legacyRcName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"policy_id",
				},
			},
		},
	})
}

func testAccVault_bindPolicies_base(name string) string {
	return fmt.Sprintf(`
variable "backup_configuration" {
  type = list(object({
    interval        = number
    days            = string
    execution_times = list(string)
  }))
  default = [{
    interval        = 5
    days            = null
    execution_times = ["06:00", "18:00"]
  },
  {
    interval        = null
    days            = "SA,SU"
    execution_times = ["08:00", "20:00"]
  }]
}

resource "huaweicloud_cbr_policy" "backup" {
  count = 2

  name            = format("%[1]s_%%d", count.index)
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    interval        = var.backup_configuration[count.index].interval
    days            = var.backup_configuration[count.index].days
    execution_times = var.backup_configuration[count.index].execution_times
  }
}

resource "huaweicloud_cbr_policy" "replication" {
  count = 2

  name                   = "%[1]s"
  type                   = "replication"
  destination_region     = "%[2]s"
  destination_project_id = "%[3]s"
  time_period            = 20

  backup_cycle {
    interval        = var.backup_configuration[count.index].interval
    days            = var.backup_configuration[count.index].days
    execution_times = var.backup_configuration[count.index].execution_times
  }
}

resource "huaweicloud_cbr_vault" "destination" {
  region          = "%[2]s"
  name            = "%[1]s"
  type            = "server"
  protection_type = "replication"
  size            = 200
}
`, name, acceptance.HW_DEST_REGION, acceptance.HW_DEST_PROJECT_ID)
}

func testAccVault_bindPolicies_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "legacy" {
  name            = "%[2]s_legacy"
  type            = "server"
  protection_type = "backup"
  size            = 200
  policy_id       = huaweicloud_cbr_policy.backup[0].id
}

resource "huaweicloud_cbr_vault" "test" {
  name            = "%[2]s"
  type            = "server"
  protection_type = "backup"
  size            = 200

  policy {
    id = huaweicloud_cbr_policy.backup[0].id
  }
  policy {
    id                   = huaweicloud_cbr_policy.replication[0].id
    destination_vault_id = huaweicloud_cbr_vault.destination.id
  }
}
`, testAccVault_bindPolicies_base(name), name)
}

func testAccVault_bindPolicies_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "legacy" {
  name            = "%[2]s_legacy"
  type            = "server"
  protection_type = "backup"
  size            = 200
  policy_id       = huaweicloud_cbr_policy.backup[1].id
}

resource "huaweicloud_cbr_vault" "test" {
  name            = "%[2]s"
  type            = "server"
  protection_type = "backup"
  size            = 200

  policy {
    id = huaweicloud_cbr_policy.backup[1].id
  }
  policy {
    id                   = huaweicloud_cbr_policy.replication[1].id
    destination_vault_id = huaweicloud_cbr_vault.destination.id
  }
}
`, testAccVault_bindPolicies_base(name), name)
}

func TestAccVault_backupWorkspace(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		basicConfig  = testAccVault_backupWorkspace_base(name)

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupWorkspace_step1(basicConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeWorkspace),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				Config: testAccVault_backupWorkspace_step2(basicConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeWorkspace),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVault_backupWorkspace_base(name string) string {
	wsDesktopName := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

%[1]s

resource "huaweicloud_workspace_service" "test" {
  access_mode = "INTERNET"
  vpc_id      = huaweicloud_vpc.test.id
  network_ids = [
    huaweicloud_vpc_subnet.test.id,
  ]
}

resource "huaweicloud_workspace_desktop" "test" {
  count = 2

  flavor_id         = "workspace.x86.ultimate.large2"
  image_type        = "market"
  image_id          = "63aa8670-27ad-4747-8c44-6d8919e785a7"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.test.id
  }

  name       = format("%[2]s-%%d", count.index)
  user_name  = format("user-%[3]s-%%d", count.index)
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }

  delete_user = true
}
`, common.TestBaseNetwork(name), wsDesktopName, name)
}

func testAccVault_backupWorkspace_step1(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "workspace"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"

  resources {
    server_id = huaweicloud_workspace_desktop.test[0].id
  }
}
`, basicConfig, name)
}

func testAccVault_backupWorkspace_step2(basicConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test" {
  name                  = "%[2]s"
  type                  = "workspace"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = "0"

  resources {
    server_id = huaweicloud_workspace_desktop.test[1].id
  }
}
`, basicConfig, name)
}

func TestAccVault_backupVMware(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupVMware_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cloud_type", "hybrid"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeVMware),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
		},
	})
}

func testAccVault_backupVMware_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  cloud_type       = "hybrid"
  name             = "%[1]s"
  type             = "vmware"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 100
}
`, name)
}

func TestAccVault_backupFile(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_backupFile_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "cloud_type", "hybrid"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeFile),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
				),
			},
		},
	})
}

func testAccVault_backupFile_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  cloud_type       = "hybrid"
  name             = "%[1]s"
  type             = "file"
  consistent_level = "crash_consistent"
  protection_type  = "backup"
  size             = 100
}
`, name)
}

func TestAccVault_locked(t *testing.T) {
	var (
		vault interface{}

		resourceName = "huaweicloud_cbr_vault.test"
		name         = acceptance.RandomAccResourceName()

		rc = acceptance.InitResourceCheck(resourceName, &vault, getVaultResourceFunc)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVault_locked_step(name, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "disk"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "locked", "true"),
				),
			},
			{
				Config: testAccVault_locked_step(name, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "disk"),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "locked", "false"),
				),
				ExpectError: regexp.MustCompile("vault not support to modify locked attribute from true to false."),
			},
		},
	})
}

func testAccVault_locked_step(name string, locked bool) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  name            = "%[1]s"
  type            = "disk"
  protection_type = "backup"
  size            = 100
  locked          = %v
}
`, name, locked)
}
