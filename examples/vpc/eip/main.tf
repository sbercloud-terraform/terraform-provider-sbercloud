# Create EIP
resource "sbercloud_vpc_eip" "eip_01" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type  = "PER"
    name        = "eip-demo"
    size        = 3
    charge_mode = "bandwidth"
  }
}
