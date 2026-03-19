package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCertificateCancelRequest_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testResourceResourceCertificateCancelRequest_basic(),
			},
		},
	})
}

func testResourceResourceCertificateCancelRequest_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_certificate_cancel_request" "test" {
  certificate_id  = "%s"
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID)
}
