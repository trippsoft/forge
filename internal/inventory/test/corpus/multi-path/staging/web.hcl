# Staging environment configuration
vars {
  staging_env = "staging"
  staging_count = 1
}

group "staging_web" {
  parent = "infrastructure"
  vars {
    role = "web"
    environment = "${vars.staging_env}"
    instance_count = "${vars.staging_count}"
  }
  host "staging-web-1" {
    vars {
      ip = "10.2.1.10"
      hostname = "staging-web-1.${vars.domain}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}
