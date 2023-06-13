package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccDataSourceGroupMetadata(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "data-sources/neuvector_group_metadata/data-source.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.neuvector_group_metadata.test", "name"),
					resource.TestCheckResourceAttrSet("data.neuvector_group_metadata.test", "services.#"),
				),
			},
		},
	})
}
