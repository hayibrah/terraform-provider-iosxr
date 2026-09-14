resource "iosxr_router_bgp" "example" {
  as_league_peers = [
    {
      peer_as_number = "65002"
    }
  ]
  as_lists = [
    {
      as_numbers = [
        {
          as_value = "65010"
        }
      ]
      list_name = "AS-LIST-1"
    }
  ]
  as_number = "65001"
  attribute_filter_groups = [
    {
      attribute_code_ranges = [
        {
          discard = true
          end     = 8
          start   = 4
        }
      ]
      group_name = "GROUP1"
    }
  ]
  bfd_minimum_interval                 = 10
  bfd_multiplier                       = 4
  bgp_as_path_loopcheck                = true
  bgp_auto_policy_soft_reset_disable   = true
  bgp_bestpath_aigp_ignore             = true
  bgp_bestpath_as_path_ignore          = true
  bgp_bestpath_as_path_multipath_relax = true
  bgp_bestpath_compare_routerid        = true
  bgp_bestpath_cost_community_ignore   = true
  bgp_bestpath_igp_metric_sr_policy    = true
  bgp_bestpath_med_missing_as_worst    = true
  bgp_bestpath_origin_as_allow_invalid = true
  bgp_bestpath_origin_as_use_validity  = true
  bgp_bestpath_sr_policy_prefer        = true
  bgp_cluster_id_32bit_format          = 100010
  bgp_confederation_identifier         = "65001"
  bgp_confederation_peers = [
    {
      peer_as_number = "65010"
    }
  ]
  bgp_default_local_preference                   = 200
  bgp_enforce_first_as_disable                   = true
  bgp_fast_external_fallover_disable             = true
  bgp_graceful_restart_enable                    = true
  bgp_graceful_restart_graceful_reset            = true
  bgp_graceful_restart_purge_time                = 120
  bgp_graceful_restart_restart_time              = 90
  bgp_graceful_restart_retain_nbr_routes_disable = true
  bgp_graceful_restart_stalepath_time            = 120
  bgp_install_diversion                          = true
  bgp_log_memory_threshold_critical              = 90
  bgp_log_memory_threshold_warning               = 80
  bgp_log_message_disable                        = true
  bgp_log_neighbor_changes_detail                = true
  bgp_log_total_paths                            = 10000
  bgp_log_total_paths_warn_threshold             = 80
  bgp_lpts_secure_binding                        = true
  bgp_maximum_neighbor                           = 5000
  bgp_multipath_as_path_ignore_onwards           = true
  bgp_multipath_use_cluster_list_length          = true
  bgp_origin_as_validation_signal_ibgp           = true
  bgp_origin_as_validation_time                  = 45
  bgp_redistribute_internal                      = true
  bgp_router_id                                  = "22.22.22.22"
  bgp_scan_time                                  = 30
  bgp_unsafe_ebgp_policy                         = true
  bgp_update_delay                               = 240
  bgp_update_delay_always                        = true
  default_information_originate                  = true
  default_metric                                 = 125
  graceful_maintenance_activate_all_neighbors    = true
  graceful_maintenance_activate_interfaces = [
    {
      interface_name = "TenGigE0/0/0/2"
    }
  ]
  graceful_maintenance_activate_locations = [
    {
      location_value = "0/RP0/CPU0"
    }
  ]
  graceful_maintenance_activate_retain_routes = true
  ibgp_policy_out_enforce_modifications       = true
  mpls_activate_interfaces = [
    {
      interface_name = "TenGigE0/0/0/2"
    }
  ]
  mvpn = true
  neighbors = [
    {
      address                                        = "10.1.1.2"
      advertisement_interval_seconds                 = 10
      allowas_in                                     = 3
      as_override                                    = "enable"
      as_path_loopcheck_out                          = "enable"
      bestpath_origin_as_allow_invalid               = true
      bfd_fast_detect                                = true
      bfd_fast_detect_strict_mode_negotiate_override = true
      bfd_minimum_interval                           = 10
      bfd_multiplier                                 = 4
      bmp_activate_servers = [
        {
          server_number = 1
        }
      ]
      capability_additional_paths_receive                            = true
      capability_additional_paths_send                               = true
      capability_suppress_all                                        = true
      capability_suppress_extended_nexthop_encoding                  = true
      capability_suppress_four_byte_as                               = true
      cluster_id_32bit_format                                        = 100010
      dampening                                                      = "enable"
      default_policy_action_in                                       = "reject"
      default_policy_action_out                                      = "reject"
      description                                                    = "My Neighbor Description"
      dmz_link_bandwidth                                             = true
      dscp                                                           = "ef"
      ebgp_multihop_maximum_hop_count                                = 10
      ebgp_multihop_mpls                                             = true
      ebgp_recv_extcommunity_dmz                                     = true
      ebgp_send_extcommunity_dmz_cumulative                          = true
      egress_engineering                                             = true
      enforce_first_as                                               = "enable"
      fast_fallover                                                  = true
      graceful_maintenance_activate                                  = true
      graceful_maintenance_as_prepends_number                        = 3
      graceful_maintenance_bandwidth_aware_percentage_threshold      = 75
      graceful_maintenance_bandwidth_aware_percentage_threshold_high = 80
      graceful_maintenance_local_preference                          = 200
      graceful_restart                                               = true
      graceful_restart_restart_time                                  = 90
      graceful_restart_stalepath_time                                = 120
      idle_watch_time                                                = 240
      ignore_connected_check                                         = true
      local_address                                                  = "192.168.1.1"
      log_message_in_size                                            = 256
      log_message_out_size                                           = 256
      log_neighbor_changes_detail                                    = true
      maximum_peers                                                  = 5
      origin_as_validation_disable                                   = true
      password                                                       = "12341C2713181F13253920"
      peer_node_sid_index                                            = 50
      receive_buffer_size                                            = 1024
      receive_buffer_size_read                                       = 1024
      remote_as                                                      = "65002"
      send_buffer_size                                               = 4096
      send_buffer_size_write                                         = 4096
      send_extended_community_ebgp                                   = true
      session_open_mode                                              = "active-only"
      shutdown                                                       = false
      tcp_mss_value                                                  = 1460
      tcp_mtu_discovery                                              = true
      timers_holdtime                                                = 30
      timers_holdtime_minimum_acceptable_holdtime                    = 30
      timers_keepalive_interval                                      = 10
      ttl_security                                                   = false
      update_in_error_handling_treat_as_withdraw                     = "enable"
      update_in_filtering_attribute_filter_group                     = "GROUP1"
      update_in_filtering_logging_disable                            = true
      update_in_filtering_message_buffers                            = 10
      update_in_filtering_message_buffers_type                       = "non-circular"
      update_in_labeled_unicast_equivalent                           = true
      update_source                                                  = "GigabitEthernet0/0/0/1"
      use_neighbor_group                                             = "GROUP1"
      use_session_group                                              = "SGROUP1"
    }
  ]
  nexthop_mpls_forwarding_ibgp               = true
  nexthop_resolution_allow_default           = true
  nexthop_validation_color_extcomm_sr_policy = true
  nsr                                        = true
  openconfig_rib_telemetry                   = true
  rpki_routes = [
    {
      max_length    = 24
      origin_as     = 501
      route_address = "172.16.1.0"
      route_prefix  = 24
    }
  ]
  rpki_servers = [
    {
      bind_source_interface = "Loopback0"
      password              = "060506324F41"
      purge_time            = 180
      refresh_time_seconds  = 120
      response_time_seconds = 240
      server                = "192.168.1.200"
      shutdown              = false
      transport_tcp_port    = 3323
      username              = "rpki-user"
    }
  ]
  segment_routing_srv6_locator                             = "locator11"
  segment_routing_srv6_usid_allocation_wide_local_id_block = true
  slow_peer_dynamic                                        = true
  slow_peer_dynamic_threshold                              = 260
  socket_receive_buffer_size                               = 1024
  socket_receive_buffer_size_read                          = 1024
  socket_send_buffer_size                                  = 4096
  socket_send_buffer_size_write                            = 4096
  timers_bgp_holdtime                                      = 30
  timers_bgp_holdtime_minimum_acceptable_holdtime          = 30
  timers_bgp_keepalive_interval                            = 10
  update_in_error_handling_basic_ebgp_disable              = true
  update_in_error_handling_basic_ibgp_disable              = true
  update_in_error_handling_extended_ebgp                   = true
  update_in_error_handling_extended_ibgp                   = true
  update_limit                                             = 20
  update_out_logging                                       = true
}
