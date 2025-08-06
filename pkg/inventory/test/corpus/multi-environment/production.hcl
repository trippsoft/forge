# Production environment
vars {
  prod_environment = "production"
  prod_log_level = "warn"
  prod_replicas = 3
}

group "prod_web" {
  vars {
    role = "web"
    domain = "www.${vars.base_domain}"
    app_port = 8080
  }
  host "prod-web1" {
    vars {
      ip = "10.1.1.10"
      tier = "primary"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "prod-web2" {
    vars {
      ip = "10.1.1.11"
      tier = "secondary"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "prod-web3" {
    vars {
      ip = "10.1.1.12"
      tier = "secondary"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}

group "prod_api" {
  vars {
    role = "api"
    domain = "api.${vars.base_domain}"
    app_port = 9000
  }
  host "prod-api1" {
    vars {
      ip = "10.1.2.10"
      tier = "primary"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
  host "prod-api2" {
    vars {
      ip = "10.1.2.11"
      tier = "secondary"
    }
    transport "ssh" {
      host = "${vars.ip}"
    }
  }
}
