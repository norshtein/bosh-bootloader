data "azurerm_public_ip" "cf-lb" {
  name                = "${var.env_id}-cf-lb-ip"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
}

resource "azurerm_dns_zone" "cf" {
  name                = "${var.env_id}.${var.system_domain}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
}

resource "azurerm_dns_a_record" "cf" {
  name                = "*"
  zone_name           = "${azurerm_dns_zone.cf.name}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  ttl                 = "300"
  records             = ["${data.azurerm_public_ip.cf-lb.ip_address}"]
}

resource "azurerm_dns_a_record" "bosh" {
  name                = "bosh"
  zone_name           = "${azurerm_dns_zone.cf.name}"
  resource_group_name = "${azurerm_resource_group.bosh.name}"
  ttl                 = "300"
  records             = ["${azurerm_public_ip.bosh.ip_address}"]
}
