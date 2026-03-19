package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV3PasswordPolicyResourceFunc(cfg *config.Config, _ *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV3PasswordPolicy(client, cfg.DomainID)
}

func TestAccV3PasswordPolicy_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_identity_password_policy.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getV3PasswordPolicyResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3PasswordPolicy_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "password_char_combination", "4"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_length", "12"),
					resource.TestCheckResourceAttr(resourceName, "number_of_recent_passwords_disallowed", "2"),
					resource.TestCheckResourceAttr(resourceName, "password_validity_period", "180"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_age", "0"),
					resource.TestCheckResourceAttr(resourceName, "maximum_consecutive_identical_chars", "0"),
					resource.TestCheckResourceAttr(resourceName, "password_not_username_or_invert", "true"),
				),
			},
			{
				Config: testAccV3PasswordPolicy_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "password_char_combination", "2"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_length", "12"),
					resource.TestCheckResourceAttr(resourceName, "number_of_recent_passwords_disallowed", "1"),
					resource.TestCheckResourceAttr(resourceName, "password_validity_period", "90"),
					resource.TestCheckResourceAttr(resourceName, "minimum_password_age", "60"),
					resource.TestCheckResourceAttr(resourceName, "maximum_consecutive_identical_chars", "4"),
					resource.TestCheckResourceAttr(resourceName, "password_not_username_or_invert", "true"),
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

const testAccV3PasswordPolicy_basic_step1 = `
resource "huaweicloud_identity_password_policy" "test" {
  password_char_combination             = 4
  minimum_password_length               = 12
  number_of_recent_passwords_disallowed = 2
  password_validity_period              = 180 
}
`

const testAccV3PasswordPolicy_basic_step2 = `
resource "huaweicloud_identity_password_policy" "test" {
  password_char_combination             = 2
  minimum_password_length               = 12
  number_of_recent_passwords_disallowed = 1
  maximum_consecutive_identical_chars   = 4
  minimum_password_age                  = 60
  password_validity_period              = 90  
}
`
