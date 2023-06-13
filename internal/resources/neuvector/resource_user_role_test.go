package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourceUserRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "resources/neuvector_user_role/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("neuvector_user_role.test", "name"),
					resource.TestCheckResourceAttr("neuvector_user_role.test", "permission.#", "1"),
				),
			},
			{
				ResourceName:      "neuvector_user_role.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
