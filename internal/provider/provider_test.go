package provider_test

import (
	"testing"

	"github.com/theobori/terraform-provider-neuvector/internal/testutils"
)

func TestProvider(t *testing.T) {
	if err := testutils.Provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
