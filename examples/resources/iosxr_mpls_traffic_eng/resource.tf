resource "iosxr_mpls_traffic_eng" "example" {
  traffic_eng = true
  disable = true
  reoptimize_reoptimization_period_in = 3600
  server_ipv4 = "192.0.2.1"
}
