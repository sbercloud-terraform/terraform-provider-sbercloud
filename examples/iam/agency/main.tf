# Create agency
resource "sbercloud_identity_agency" "agency_01" {
  name                   = "tf_test_agency"
  description            = "Allow FG manage ECS"
  delegated_service_name = "op_svc_cff"
  duration               = "ONEDAY"

  project_role {
    project = "ru-moscow-1"
    roles = [
      "ECS User",
    ]
  }
}
