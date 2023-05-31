package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourceEULA(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "resources/neuvector_eula/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("neuvector_eula.test", "accepted"),
					resource.TestCheckResourceAttr("neuvector_eula.test", "accepted", "true"),
				),
			},
			{
				ResourceName:            "neuvector_eula.test",
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}
