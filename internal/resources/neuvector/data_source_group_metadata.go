// data_source_group_services.go
package neuvector

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	goneuvector "github.com/theobori/go-neuvector/neuvector"
	"github.com/theobori/go-neuvector/util"
)

var AllowedMetadata = []string{
	"services",
	"container_ids",
	"image_ids",
}

var dataGroupMetadataSchema = map[string]*schema.Schema{
	"services": {
		Type:        schema.TypeSet,
		Description: "List of every service name in the group.",
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"container_ids": {
		Type:        schema.TypeSet,
		Description: "List of every container id in the group.",
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"image_ids": {
		Type:        schema.TypeSet,
		Description: "List of every image id in the group.",
		Computed:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the gorup.",
		Required:    true,
	},
}

func DataSourceGroupMetadata() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupMetadataRead,
		Schema:      dataGroupMetadataSchema,
	}
}

func getMemberFieldFromString(m *goneuvector.WorkloadBrief, s string) (string, error) {
	var ret string

	switch s {
	case "services":
		ret = m.Service
	case "container_ids":
		ret = m.ID
	case "image_ids":
		ret = m.ImageID
	default:
		return ret, fmt.Errorf("invalid key")
	}

	return ret, nil
}

func addGroupInfo(
	infos *map[string][]string,
	m *goneuvector.WorkloadBrief,
	key string,
) error {
	value, err := getMemberFieldFromString(m, key)

	if err != nil {
		return err
	}

	exists, _ := util.ItemExists(
		(*infos)[key],
		value,
	)

	if exists {
		return fmt.Errorf("already exists")
	}

	(*infos)[key] = append((*infos)[key], value)

	return nil
}

func addGroupMetadata(
	infos *map[string][]string,
	m *goneuvector.WorkloadBrief,
	keys ...string,
) {
	for _, key := range keys {
		// We do not need to handle the error there
		addGroupInfo(infos, m, key)
	}
}

func dataSourceGroupMetadataRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	infos := map[string][]string{}

	APIClient := meta.(*goneuvector.Client)

	name := d.Get("name").(string)
	group, err := APIClient.
		WithContext(ctx).
		GetGroup(name)

	if err != nil {
		return diag.FromErr(err)
	}

	members := group.Group.Members

	for _, m := range members {
		addGroupMetadata(&infos, &m, AllowedMetadata...)
	}

	for k, v := range infos {
		d.Set(k, v)
	}

	id, err := uuid.GenerateUUID()

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return nil
}
