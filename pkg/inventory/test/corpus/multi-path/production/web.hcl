# Production environment configuration
vars {
  prod_env = "production"
  prod_count = 3
}

group "prod_web" {
  parent = "infrastructure"
  vars {
    role = "web"
    environment = "${vars.prod_env}"
    instance_count = "${vars.prod_count}"
  }
  host "prod-web-1" {
    vars {
      ip = "10.1.1.10"
      hostname = "prod-web-1.${vars.domain}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "prod-web-2" {
    vars {
      ip = "10.1.1.11"
      hostname = "prod-web-2.${vars.domain}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "prod-web-3" {
    vars {
      ip = "10.1.1.12"
      hostname = "prod-web-3.${vars.domain}"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}
