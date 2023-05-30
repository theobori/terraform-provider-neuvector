package neuvector_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourceRegistry(t *testing.T) {
	var r goneuvector.Registry

	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		CheckDestroy:      testAccRegistryCheckDestroy(&r),
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: false,
				Config:             testutils.TestAccExampleFile(t, "resources/neuvector_registry/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("neuvector_registry.test", "filters.#", "1"),
					testAccRegistryCheckExists("neuvector_registry.test", &r),
					resource.TestCheckResourceAttrSet("neuvector_registry.test", "cfg_type"),
				),
			},
			{
				ResourceName:            "neuvector_registry.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cfg_type"},
			},
		},
	})
}

func testAccRegistryCheckExists(rn string, r *goneuvector.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]

		if !ok {
			return fmt.Errorf("resource not found: %s\n %#v", rn, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		APIClient := testutils.Provider.Meta().(*goneuvector.Client)

		registry, err := APIClient.GetRegistry(rs.Primary.ID)

		if err != nil {
			return err
		}

		*r = registry.Registry

		return nil
	}
}

func testAccRegistryCheckDestroy(r *goneuvector.Registry) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		APIClient := testutils.Provider.Meta().(*goneuvector.Client)

		_, err := APIClient.GetRegistry(r.Name)

		if err == nil {
			return fmt.Errorf("registry still exists")
		}

		return nil
	}
}
