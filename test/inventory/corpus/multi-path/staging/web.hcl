# Staging environment configuration
vars {
    staging_env = "staging"
    staging_count = 1
}

group "staging_web" {
    parent = "infrastructure"

    vars {
        role = "web"
        environment = "${var.staging_env}"
        instance_count = "${var.staging_count}"
    }
}

host "staging-web-1" {
    groups = ["staging_web"]

    vars {
        ip = "10.2.1.10"
        hostname = "staging-web-1.${var.domain}"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}
