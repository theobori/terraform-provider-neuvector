package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourceGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             testutils.TestAccExampleFile(t, "resources/neuvector_group/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("neuvector_group.test", "criteria.#", "1"),
					resource.TestCheckResourceAttr("neuvector_group.test", "name", "mytestgroup"),
					resource.TestCheckResourceAttrSet("neuvector_group.test", "cfg_type"),
				),
			},
			{
				ResourceName:      "neuvector_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
