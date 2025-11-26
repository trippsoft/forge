# Production environment configuration
vars {
    prod_env = "production"
    prod_count = 3
}

group "prod_web" {
    parent = "infrastructure"

    vars {
        role = "web"
        environment = "${var.prod_env}"
        instance_count = "${var.prod_count}"
    }
}

host "prod-web-1" {
    groups = ["prod_web"]

    vars {
        ip = "10.1.1.10"
        hostname = "prod-web-1.${var.domain}"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "prod-web-2" {
    groups = ["prod_web"]

    vars {
      ip = "10.1.1.11"
      hostname = "prod-web-2.${var.domain}"
    }

    transport "ssh" {
      host = "${var.ip}"
    }
}

host "prod-web-3" {
    groups = ["prod_web"]

    vars {
        ip = "10.1.1.12"
        hostname = "prod-web-3.${var.domain}"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
