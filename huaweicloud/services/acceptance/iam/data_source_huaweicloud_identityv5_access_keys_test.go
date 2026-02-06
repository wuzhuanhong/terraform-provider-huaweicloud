package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccDataV5AccessKeys_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_access_keys.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV5AccessKeys_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(all, "user_id", "huaweicloud_identityv5_user.test", "id"),
					resource.TestCheckResourceAttr(all, "access_keys.#", "1"),
					resource.TestCheckResourceAttrPair(all, "access_keys.0.access_key_id",
						"huaweicloud_identityv5_access_key.test", "access_key_id"),
					resource.TestCheckResourceAttrPair(all, "access_keys.0.user_id",
						"huaweicloud_identityv5_user.test", "id"),
					resource.TestCheckResourceAttr(all, "access_keys.0.status", "inactive"),
					resource.TestCheckResourceAttrSet(all, "access_keys.0.created_at"),
				),
			},
		},
	})
}

func testAccDataV5AccessKeys_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_user.test.id
  status  = "inactive"
}

data "huaweicloud_identityv5_access_keys" "test" {
  user_id = huaweicloud_identityv5_user.test.id

  depends_on = [huaweicloud_identityv5_access_key.test]
}
`, name)
}
