package neuvector_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/theobori/go-neuvector/client"
	"github.com/theobori/go-neuvector/controller/admission"
	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestAccResourceAdmissionRule(t *testing.T) {
	var adm admission.AdmissionRule

	resource.Test(t, resource.TestCase{
		ProviderFactories: testutils.ProviderFactories,
		CheckDestroy:      testAccAdmissionRuleCheckDestroy(&adm),
		Steps: []resource.TestStep{
			{
				Config: testutils.TestAccExampleFile(t, "resources/neuvector_admission_rule/resource.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("neuvector_admission_rule.test", "criteria.#", "2"),
					testAccAdmissionRuleCheckExists("neuvector_admission_rule.test", &adm),
					resource.TestCheckResourceAttrSet("neuvector_admission_rule.test", "disable"),
					resource.TestCheckResourceAttrSet("neuvector_admission_rule.test", "cfg_type"),
				),
			},
			{
				ResourceName:            "neuvector_admission_rule.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rule_mode"},
			},
		},
	})
}

func testAccAdmissionRuleCheckExists(rn string, adm *admission.AdmissionRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]

		if !ok {
			return fmt.Errorf("resource not found: %s\n %#v", rn, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		APIClient := testutils.Provider.Meta().(*client.Client)
		id, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return err
		}

		rule, err := admission.GetAdmissionRule(APIClient, id)

		if err != nil {
			return err
		}

		*adm = rule.Rule

		return nil
	}
}

func testAccAdmissionRuleCheckDestroy(adm *admission.AdmissionRule) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		APIClient := testutils.Provider.Meta().(*client.Client)

		_, err := admission.GetAdmissionRule(APIClient, adm.ID)

		if err == nil {
			return fmt.Errorf("admission rule still exists")
		}

		return nil
	}
}
