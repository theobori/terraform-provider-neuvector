package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourceService(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "resources/neuvector_service/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("neuvector_service.test", "name"),
					resource.TestCheckResourceAttrSet("neuvector_service.test", "comment"),
					resource.TestCheckResourceAttrSet("neuvector_service.test", "not_scored"),
				),
			},
			{
				ResourceName:            "neuvector_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"comment"},
			},
		},
	})
}
