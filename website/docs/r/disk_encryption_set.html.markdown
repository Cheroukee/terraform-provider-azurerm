---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_encryption_set"
description: |-
  Manages a Disk Encryption Set.
---

# azurerm_disk_encryption_set

Manages a Disk Encryption Set.


## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                        = "des-example-keyvault"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  sku_name                    = "premium"
  enabled_for_disk_encryption = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_key" "example" {
  name         = "des-example-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  depends_on = [
    azurerm_key_vault_access_policy.example-user
  ]

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
}

resource "azurerm_disk_encryption_set" "example" {
  name                = "des"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  key_vault_key_id    = azurerm_key_vault_key.example.id

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_key_vault_access_policy" "example-disk" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = azurerm_disk_encryption_set.example.identity.0.tenant_id
  object_id = azurerm_disk_encryption_set.example.identity.0.principal_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign"
  ]
}

resource "azurerm_key_vault_access_policy" "example-user" {
  key_vault_id = azurerm_key_vault.example.id

  tenant_id = data.azurerm_client_config.current.tenant_id
  object_id = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign"
  ]
}

resource "azurerm_role_assignment" "example-disk" {
  scope                = azurerm_key_vault.example.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = azurerm_disk_encryption_set.example.identity.0.principal_id
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Disk Encryption Set. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Disk Encryption Set should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Disk Encryption Set exists. Changing this forces a new resource to be created.

* `key_vault_key_id` - (Required) Specifies the URL to a Key Vault Key (either from a Key Vault Key, or the Key URL for the Key Vault Secret).

-> **NOTE** Access to the KeyVault must be granted for this Disk Encryption Set, if you want to further use this Disk Encryption Set in a Managed Disk or Virtual Machine, or Virtual Machine Scale Set. For instructions, please refer to the doc of [Server side encryption of Azure managed disks](https://docs.microsoft.com/azure/virtual-machines/linux/disk-encryption).

-> **NOTE** A KeyVault using [enable_rbac_authorization](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault#enable_rbac_authorization) requires to use `azurerm_role_assignment` to assigne the role `Key Vault Crypto Service Encryption User` to this Disk Encryption Set.
In this case, `azurerm_key_vault_access_policy` is not needed.

* `auto_key_rotation_enabled` - (Optional) Boolean flag to specify whether Azure Disk Encryption Set automatically rotates encryption Key to latest version. Defaults to `false`.

* `encryption_type` - (Optional) The type of key used to encrypt the data of the disk. Possible values are `EncryptionAtRestWithCustomerKey`, `EncryptionAtRestWithPlatformAndCustomerKeys` and `ConfidentialVmEncryptedWithCustomerKey`. Defaults to `EncryptionAtRestWithCustomerKey`.

* `identity` - (Required) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the Disk Encryption Set.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that is configured on this Disk Encryption Set. The only possible value is `SystemAssigned`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Disk Encryption Set.

---

An `identity` block exports the following:

* `principal_id` - The (Client) ID of the Service Principal.

* `tenant_id` - The ID of the Tenant the Service Principal is assigned in.

## Timeouts



The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Disk Encryption Set.
* `update` - (Defaults to 60 minutes) Used when updating the Disk Encryption Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disk Encryption Set.
* `delete` - (Defaults to 60 minutes) Used when deleting the Disk Encryption Set.

## Import

Disk Encryption Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_disk_encryption_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/diskEncryptionSets/encryptionSet1
```
