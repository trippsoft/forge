# Production environment
vars {
    prod_environment = "production"
    prod_log_level = "warn"
    prod_replicas = 3
}

group "prod_web" {
    vars {
        role = "web"
        domain = "www.${var.base_domain}"
        app_port = 8080
    }
}

host "prod-web1" {
    groups = ["prod_web"]

    vars {
        ip = "10.1.1.10"
        tier = "primary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "prod-web2" {
    vars {
        ip = "10.1.1.11"
        tier = "secondary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "prod-web3" {
    groups = ["prod_web"]

    vars {
        ip = "10.1.1.12"
        tier = "secondary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

group "prod_api" {
    vars {
        role = "api"
        domain = "api.${var.base_domain}"
        app_port = 9000
    }
}

host "prod-api1" {
    groups = ["prod_api"]

    vars {
        ip = "10.1.2.10"
        tier = "primary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
host "prod-api2" {
    groups = ["prod_api"]

    vars {
        ip = "10.1.2.11"
        tier = "secondary"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
