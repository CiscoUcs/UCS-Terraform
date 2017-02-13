provider "ucs" {
  ip_address   = "1.2.3.4"
  username     = "john"
  password     = "supersecret"
  log_level    = 1
  log_filename = "terraform.log"
}

resource "ucs_service_profile" "the-server-name" {
  name                     = "the-server-name"
  target_org               = "some-target-org"
  service_profile_template = "some-service-profile-template"
  vNIC {
    name  = "eth0"
    cidr = "1.2.3.4/24"
  }
}
