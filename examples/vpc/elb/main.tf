# Get the subnet where ELB will be created
data "sbercloud_vpc_subnet" "subnet_01" {
  name = "place_here_the_name_of_your_existing_subnet"
}

# Get the latest Ubuntu image
data "sbercloud_images_image" "ubuntu_image" {
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

# Create two ECS
resource "sbercloud_compute_instance" "ecs_01" {
  count = 2

  name              = "ecs-terraform-${count.index}"
  image_id          = data.sbercloud_images_image.ubuntu_image.id
  flavor_id         = "s6.large.2"
  security_groups   = ["default"]
  availability_zone = "ru-moscow-1a"
  key_pair          = "place_here_the_name_of_your_existing_key_pair"
  system_disk_type = "SAS"
  system_disk_size = 16

  network {
    uuid = data.sbercloud_vpc_subnet.subnet_01.id
  }
}

# Create Elastic IP (EIP) for ELB
resource "sbercloud_vpc_eip" "elb_eip" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "elb_bandwidth"
    size        = 3
    share_type  = "PER"
    charge_mode = "bandwidth"
  }
}

# Create ELB
resource "sbercloud_lb_loadbalancer" "elb_01" {
  name          = "elb-frontend"
  vip_subnet_id = data.sbercloud_vpc_subnet.subnet_01.subnet_id
}

# Attach EIP to ELB
resource "sbercloud_networking_eip_associate" "elb_eip_associate" {
  public_ip = sbercloud_vpc_eip.elb_eip.address
  port_id   = sbercloud_lb_loadbalancer.elb_01.vip_port_id
}

# Create ELB listener
resource "sbercloud_lb_listener" "listener_01" {
  name            = "Demo listener"
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = sbercloud_lb_loadbalancer.elb_01.id
}

# Create ECS backend group for ELB
resource "sbercloud_lb_pool" "backend_pool" {
  name        = "Backend servers group for ELB"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = sbercloud_lb_listener.listener_01.id
}

# Create ELB health check policy
resource "sbercloud_lb_monitor" "elb_health_check" {
  name           = "Health check for ECS"
  type           = "HTTP"
  url_path       = "/"
  expected_codes = "200-202"
  delay          = 10
  timeout        = 5
  max_retries    = 3
  pool_id        = sbercloud_lb_pool.backend_pool.id
}

# Add both ECS to the backend server group
resource "sbercloud_lb_member" "backend_server" {
  count = 2

  address = sbercloud_compute_instance.ecs_01[count.index].access_ip_v4

  protocol_port = 80
  pool_id       = sbercloud_lb_pool.backend_pool.id
  subnet_id     = data.sbercloud_vpc_subnet.subnet_01.subnet_id

  depends_on = [
    sbercloud_lb_monitor.elb_health_check
  ]
}
