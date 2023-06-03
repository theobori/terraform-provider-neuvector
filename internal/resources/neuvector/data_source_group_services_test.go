package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccDataSourceGroupServices(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "data-sources/neuvector_group_services/data-source.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.neuvector_group_services.test", "name"),
					resource.TestCheckResourceAttrSet("data.neuvector_group_services.test", "services.#"),
				),
			},
		},
	})
}
