package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "resources/neuvector_user/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("neuvector_user.test", "fullname"),
					resource.TestCheckResourceAttrSet("neuvector_user.test", "username"),
					resource.TestCheckResourceAttrSet("neuvector_user.test", "email"),
					resource.TestCheckResourceAttrSet("neuvector_user.test", "role"),
				),
			},
			{
				ResourceName:            "neuvector_user.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}
