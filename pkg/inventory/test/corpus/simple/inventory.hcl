# Simple inventory with basic hosts and groups
vars {
    environment = "test"
    domain      = "example.com"
}

transport "ssh" {
    user = "admin"
    port = 22
    connection_timeout = "30s"
    use_known_hosts = false
}

group "webservers" {
    vars {
        role = "web"
        port = 8080
    }
}

host "web1" {
    groups = ["webservers"]

    vars {
        ip = "10.0.1.10"
        hostname = "web1.${var.domain}"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

host "web2" {
    groups = ["webservers"]

    vars {
        ip = "10.0.1.11"
        hostname = "web2.${var.domain}"
    }

    transport "ssh" {
        host = "${var.ip}"
    }
}

group "databases" {
    vars {
        role = "database"
        port = 5432
    }
}

host "db1" {
    groups = ["databases"]

    vars {
        ip = "10.0.2.10"
        hostname = "db1.${var.domain}"
    }

    transport "ssh" {
        host = "${var.ip}"
        user = "dbuser"
    }
}
