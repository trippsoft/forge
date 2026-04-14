# Staging environment
vars {
  staging_environment = "staging"
  staging_log_level = "debug"
  staging_replicas = 1
}

group "staging_web" {
    vars {
        role = "web"
        domain = "staging.${var.base_domain}"
        app_port = 8080
    }
}

host "staging-web1" {
    groups = ["staging_web"]

    vars {
        ip = "10.2.1.10"
        tier = "single"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

group "staging_api" {
    vars {
        role = "api"
        domain = "staging-api.${var.base_domain}"
        app_port = 9000
    }
}

host "staging-api1" {
    groups = ["staging_api"]

    vars {
        ip = "10.2.2.10"
        tier = "single"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
