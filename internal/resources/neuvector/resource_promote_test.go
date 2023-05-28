package neuvector_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourcePromote(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "resources/neuvector_promote/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("neuvector_promote.test", "port"),
					resource.TestCheckResourceAttrSet("neuvector_promote.test", "server"),
					resource.TestCheckResourceAttrSet("neuvector_promote.test", "user"),
					resource.TestCheckResourceAttrSet("neuvector_promote.test", "name"),
				),
			},
		},
	})
}
