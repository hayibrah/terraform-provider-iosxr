resource "iosxr_logging_vrf" "example" {
  host_ipv4_addresses = [
    {
      facility            = "local0"
      ipv4_address        = "1.1.1.1"
      ipv4_source_address = "1.1.1.2"
      operator            = "equals"
      severity            = "informational"
    }
  ]
  host_ipv6_addresses = [
    {
      facility            = "local0"
      ipv6_address        = "2001:db8::1"
      ipv6_source_address = "2001:db8::2"
      operator            = "equals-or-higher"
      severity            = "informational"
    }
  ]
  hostnames = [
    {
      facility                = "local0"
      hostname_source_address = "1.1.1.2"
      name                    = "server.cisco.com"
      operator                = "equals"
      severity                = "informational"
    }
  ]
  vrf_name = "default"
}
