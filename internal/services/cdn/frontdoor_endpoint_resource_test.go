package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/afdendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorProfileEndpointResource struct{}

func TestAccFrontdoorEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_endpoint", "test")
	r := FrontdoorProfileEndpointResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFrontdoorEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_endpoint", "test")
	r := FrontdoorProfileEndpointResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccFrontdoorEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_endpoint", "test")
	r := FrontdoorProfileEndpointResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccFrontdoorEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_frontdoor_endpoint", "test")
	r := FrontdoorProfileEndpointResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r FrontdoorProfileEndpointResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := afdendpoints.ParseAfdEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cdn.FrontDoorEndpointsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r FrontdoorProfileEndpointResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-afdx-%d"
  location = "%s"
}

resource "azurerm_frontdoor_profile" "test" {
  name                = "acctest-c-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r FrontdoorProfileEndpointResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_frontdoor_endpoint" "test" {
  name                            = "acctest-c-%d"
  frontdoor_profile_id            = azurerm_frontdoor_profile.test.id
  enabled_state                   = ""
  location                        = "%s"
  origin_response_timeout_seconds = 0

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r FrontdoorProfileEndpointResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_endpoint" "import" {
  name                            = azurerm_frontdoor_endpoint.test.name
  frontdoor_profile_id            = azurerm_frontdoor_profile.test.id
  enabled_state                   = ""
  location                        = "%s"
  origin_response_timeout_seconds = 0

  tags = {
    ENV = "Test"
  }
}
`, config, data.Locations.Primary)
}

func (r FrontdoorProfileEndpointResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_endpoint" "test" {
  name                            = "acctest-c-%d"
  frontdoor_profile_id            = azurerm_frontdoor_profile.test.id
  enabled_state                   = ""
  location                        = "%s"
  origin_response_timeout_seconds = 0

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r FrontdoorProfileEndpointResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_frontdoor_endpoint" "test" {
  name                            = "acctest-c-%d"
  frontdoor_profile_id            = azurerm_frontdoor_profile.test.id
  enabled_state                   = ""
  origin_response_timeout_seconds = 0

  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}